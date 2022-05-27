package store

import (
	"context"
	"database/sql"
	"fmt"
	"sshtest/pkg/data"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type Client struct {
	db      *sqlx.DB
	cursors map[CursorKey]Cursor
}

func NewClient(db *sqlx.DB) Client {
	return Client{
		db:      db,
		cursors: make(map[CursorKey]Cursor),
	}
}

func (c *Client) CreateFile(ctx context.Context, f *data.File, tags data.TagSet) error {
	txn, err := c.db.BeginTxx(ctx, &sql.TxOptions{
		ReadOnly:  false,
		Isolation: sql.LevelDefault,
	})

	_, err = txn.NamedExec(`INSERT INTO files (id,path,created) VALUES(:id, :path, :created)`, &f)
	if err != nil {
		return err
	}

	for _, tag := range tags {
		_, err := txn.Exec("INSERT INTO tags (id,value) VALUES(?, ?)", f.Id, tag)
		if err != nil {
			return err
		}
	}

	err = txn.Commit()
	if err != nil {
		return err
	}

	// all cursors are marked dirty, as the DB has changed
	// this assumes the DB is only accessed by this program!
	for _, cursor := range c.cursors {
		cursor.dirty = true
	}

	return nil
}

func (c *Client) GetTags(ctx context.Context, f *data.File) (data.TagSet, error) {
	rows, err := c.db.QueryxContext(ctx, "SELECT value FROM tags WHERE id = ?", f.Id)
	if err != nil {
		return nil, err
	}

	var tags data.TagSet
	err = rows.Scan(&tags)
	if err != nil {
		return nil, err
	}

	return tags, nil
}

func (c *Client) GetFiles(ctx context.Context, cursor CursorKey) ([]data.File, error) {
	cursorValue, ok := c.cursors[cursor]
	if !ok {
		return nil, fmt.Errorf("no such cursor: %s", uuid.UUID(cursor).String())
	}
	rows, err := cursorValue.query(ctx, c.db)
	cursorValue.nextPage()
	if err != nil {
		return nil, err
	}

	var files []data.File
	for rows.Next() {
		var file data.File
		err = rows.StructScan(&file)
		if err != nil {
			return nil, err
		}
		files = append(files, file)
	}

	return files, err
}
