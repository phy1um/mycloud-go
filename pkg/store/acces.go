package store

import "sshtest/pkg/data"

func (c *Client) CreateAccessKey(key *data.Access) error {
	_, err := c.db.NamedExec(
		"INSERT INTO access_keys (path,key,user_code,display_name,created,until) VALUES (:path, :key, :user_code, :display_name, :created, :until)",
		key,
	)
	return err
}
