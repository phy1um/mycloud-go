package store

import (
	"context"
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

type Cursor interface {
	query(ctx context.Context, db *sqlx.DB) (*sqlx.Rows, error)
	setError(err error)
	setDirty(b bool)
	nextPage()
}

type cursorShared struct {
	offset   int
	pageSize int
	order    Order
	orderBy  string
	dirty    bool
	err      error
}

func (c cursorShared) queryTail() string {
	return fmt.Sprintf(
		"OFFSET %d LIMIT %d ORDER BY %s %s",
		c.offset,
		c.pageSize,
		c.orderBy,
		string(c.order),
	)
}

func (c *Client) NewCursor(cursor Cursor) CursorKey {
	// TODO: handle this error?
	id, _ := uuid.NewUUID()
	key := CursorKey(id)
	c.cursors[key] = cursor
	return key
}

func (c *Client) DestroyCursor(key CursorKey) {
	delete(c.cursors, key)
}

type AllCursor struct {
	cursorShared
}

func (c *AllCursor) query(ctx context.Context, db *sqlx.DB) (*sqlx.Rows, error) {
	return db.QueryxContext(ctx, "SELECT * FROM files "+c.queryTail())
}

func (c *AllCursor) setError(err error) {
	c.err = err
}

func (c *AllCursor) setDirty(b bool) {
	c.dirty = b
}

func (c *AllCursor) nextPage() {
	c.offset += c.pageSize
}

type FileNameCursor struct {
	search string
	cursorShared
}

func (c *FileNameCursor) query(ctx context.Context, db *sqlx.DB) (*sqlx.Rows, error) {
	return db.QueryxContext(ctx, "SELECT * FROM files WHERE name LIKE ? "+c.queryTail(), "%"+c.search+"%")
}

func (c *FileNameCursor) setError(err error) {
	c.err = err
}

func (c *FileNameCursor) setDirty(b bool) {
	c.dirty = b
}

func (c *FileNameCursor) nextPage() {
	c.offset += c.pageSize
}

type TagCursor struct {
	tag string
	cursorShared
}

func (c *TagCursor) query(ctx context.Context, db *sqlx.DB) (*sqlx.Rows, error) {
	return db.QueryxContext(
		ctx,
		"SELECT * FROM files JOIN tags WHERE file.id = tags.id AND tags.value LIKE ?"+c.queryTail(),
		"%"+c.tag+"%",
	)
}

func (c *TagCursor) setError(err error) {
	c.err = err
}

func (c *TagCursor) setDirty(b bool) {
	c.dirty = b
}

func (c *TagCursor) nextPage() {
	c.offset += c.pageSize
}
