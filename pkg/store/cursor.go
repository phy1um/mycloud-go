package store

import (
	"context"
	"fmt"
	"log"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type CursorKey uuid.UUID
type Order string

const (
	Descend Order = "DESC"
	Ascend  Order = "ASC"
)

type CursorFunc func(ctx context.Context, tail string, db *sqlx.DB) (*sqlx.Rows, error)

type Cursor struct {
	offset   int
	pageSize int
	order    Order
	orderBy  string
	dirty    bool
	err      error
	exec     CursorFunc
}

func (c Cursor) setError(err error) {
	c.err = err
}

func (c Cursor) setDirty(b bool) {
	c.dirty = b
}

func (c Cursor) nextPage() {
	c.offset += c.pageSize
}

func (c Cursor) query(ctx context.Context, db *sqlx.DB) (*sqlx.Rows, error) {
	return c.exec(ctx, c.queryTail(), db)
}

func (c Cursor) queryTail() string {
	return fmt.Sprintf(
		"LIMIT %d OFFSET %d",
		c.pageSize,
		c.offset,
	)
}

func (c *Client) NewCursor(pageSize int, orderBy string, order Order, fn CursorFunc) CursorKey {
	cursor := Cursor{
		offset:   0,
		pageSize: pageSize,
		order:    order,
		orderBy:  orderBy,
		exec:     fn,
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

func AllFiles() CursorFunc {
	return func(ctx context.Context, tail string, db *sqlx.DB) (*sqlx.Rows, error) {
		query := "SELECT * FROM files " + tail
		log.Printf("running query: %s\n", query)
		return db.QueryxContext(ctx, query)
	}
}

func FileNameSearch(name string) CursorFunc {
	return func(ctx context.Context, tail string, db *sqlx.DB) (*sqlx.Rows, error) {
		query := "SELECT * FROM files WHERE files.name LIKE ? " + tail
		log.Printf("running query: %s\n", query)
		return db.QueryxContext(ctx, query, "%"+name+"%")
	}
}

func FileTagSearch(tag string) CursorFunc {
	return func(ctx context.Context, tail string, db *sqlx.DB) (*sqlx.Rows, error) {
		query := "SELECT (files.id, files.path, files.name, files.created) FROM files JOIN tags WHERE files.id = tags.id AND tags.value LIKE ? " + tail
		log.Printf("running query: %s\n", query)
		return db.QueryxContext(
			ctx,
			query,
			"%"+tag+"%",
		)
	}
}
