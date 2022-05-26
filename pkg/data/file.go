package data

import "time"

type File struct {
	Id      string    `db:"id"`
	Path    string    `db:"path"`
	Name    string    `db:"name"`
	Created time.Time `db:"created"`
}

type TagSet []string
