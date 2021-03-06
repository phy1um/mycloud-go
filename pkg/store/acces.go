package store

import (
	"context"

	"github.com/phy1um/mycloud-go/pkg/data"

	"github.com/rs/zerolog/log"
)

func (c *Client) CreateAccessKey(key *data.Access) error {
	_, err := c.db.NamedExec(
		"INSERT INTO access_keys (file_id,key,user_code,display_name,created,until) VALUES (:file_id, :key, :user_code, :display_name, :created, :until)",
		key,
	)
	return err
}

func (c *Client) GetAccessKeys(ctx context.Context, file *data.File) ([]*data.Access, error) {
	log.Ctx(ctx).Info().Msgf("get access keys for file %s", file.Id)
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
