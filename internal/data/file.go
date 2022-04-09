package data

import "time"

type File struct {
	Id      string
	Path    string
	Tag     string
	Created time.Time
}
