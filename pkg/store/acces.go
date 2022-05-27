package store

import (
	"context"
	"sshtest/pkg/data"
)

func (c *Client) CreateAccessKey(key *data.Access) error {
	_, err := c.db.NamedExec(
		"INSERT INTO access_keys (file_id,key,user_code,display_name,created,until) VALUES (:file_id, :key, :user_code, :display_name, :created, :until)",
		key,
	)
	return err
}

func (c *Client) GetAccessKeys(ctx context.Context, file *data.File) ([]*data.Access, error) {
	res, err := c.db.QueryxContext(ctx, "SELECT * FROM access_keys WHERE file_id = ?", file.Id)
	if err != nil {
		return nil, err
	}

	var keys []*data.Access
	for res.Next() {
		var access data.Access
		err = res.StructScan(&access)
		if err != nil {
			return nil, err
		}
		keys = append(keys, &access)
	}

	return keys, nil
}
