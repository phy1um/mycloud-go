package store

import "sshtest/pkg/data"

func (c *Client) CreateAccessKey(key *data.Access) error {
	_, err := c.db.NamedExec(
		"INSERT INTO access_keys (file_id,key,user_code,display_name,created,until) VALUES (:file_id, :key, :user_code, :display_name, :created, :until)",
		key,
	)
	return err
}
