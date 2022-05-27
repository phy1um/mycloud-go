package store

import (
	"context"
	"database/sql"

	"github.com/phy1um/mycloud-go/pkg/data"
)

func (c *Client) AddTag(ctx context.Context, fileId string, tag string) error {
	txn, err := c.db.BeginTxx(ctx, &sql.TxOptions{})
	if err != nil {
		return err
	}

	_, err = txn.Exec("INSERT INTO tags (id, value) VALUES (?, ?)", fileId, tag)
	if err != nil {
		return err
	}

	err = txn.Commit()
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) GetTags(ctx context.Context, fileId string) (data.TagSet, error) {
	rows, err := c.db.QueryxContext(ctx, "SELECT value FROM tags WHERE id = ?", fileId)
	if err != nil {
		return nil, err
	}

	var tags data.TagSet
	for rows.Next() {
		var tag string
		err = rows.Scan(&tag)
		if err != nil {
			return nil, err
		}
		tags = append(tags, tag)
	}

	return tags, nil
}
