package store

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type CursorKey uuid.UUID
type Order string

const (
	Descend Order = "DESC"
	Ascend  Order = "ASC"
)

type Cursor struct {
	offset   int
	pageSize int
	table    string
	order    Order
	orderBy  string
	tag      string
	dirty    bool
	err      error
}

func (c *Client) NewCursor(pageSize int, tag string) CursorKey {
	cursor := Cursor{
		offset:   0,
		pageSize: pageSize,
		tag:      tag,
	}
	// TODO: handle this error?
	id, _ := uuid.NewUUID()
	key := CursorKey(id)
	c.cursors[key] = cursor
	return key
}

func (c *Client) DestroyCursor(key CursorKey) {
	delete(c.cursors, key)
}

func (c *Client) nextPage(ctx context.Context, key CursorKey) (*sqlx.Rows, error) {
	cursor, ok := c.cursors[key]
	if !ok {
		return nil, errors.New("unknown cursor key")
	}
	// I think QueryxContext args will quote strings which I don't think we want here
	var query string
	if cursor.tag != "" {
		query = fmt.Sprintf("SELECT * FROM files JOIN tags WHERE files.id = tags.id AND tags.value = \"%s\"", cursor.tag)
	} else {
		query = fmt.Sprintf("SELECT * FROM files")
	}
	res, err := c.db.QueryxContext(ctx, query)
	if err == nil {
		cursor.offset += cursor.pageSize
		cursor.err = err
	}
	return res, err
}
