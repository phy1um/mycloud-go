package store

import (
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
