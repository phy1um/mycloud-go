package data

import "time"

type File struct {
	Id      string    `db:"id"`
	Path    string    `db:"path"`
	Name    string    `db:"name"`
	Tag     string    `db:"tag"`
	Created time.Time `db:"created"`
}
