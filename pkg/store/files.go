package store

import (
	"context"
	"database/sql"
	"fmt"
	"sshtest/pkg/data"
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

func (c *Client) CreateFile(ctx context.Context, id string, path string, tags data.TagSet) error {
	txn, err := c.db.BeginTxx(ctx, &sql.TxOptions{
		ReadOnly:  false,
		Isolation: sql.LevelDefault,
	})

	_, err = txn.Exec(`INSERT INTO files (id,path,created) VALUES(?, ?, ?)`, id, path, time.Now())
	if err != nil {
		return err
	}

	for _, tag := range tags {
		_, err := txn.Exec("INSERT INTO tags (id,value) VALUES(?, ?)", id, tag)
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

func (c *Client) GetFiles(ctx context.Context, cursor CursorKey) ([]*data.File, error) {
	cursorValue, ok := c.cursors[cursor]
	if !ok {
		return nil, fmt.Errorf("no such cursor: %s", uuid.UUID(cursor).String())
	}
	rows, err := cursorValue.query(ctx, c.db)
	cursorValue.nextPage()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get files")
	}

	var files []*data.File
	for rows.Next() {
		var file data.File
		err = rows.StructScan(&file)
		if err != nil {
			return nil, err
		}
		files = append(files, &file)
	}

	return files, err
}
