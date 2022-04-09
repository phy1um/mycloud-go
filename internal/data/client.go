package data

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/jmoiron/sqlx"

	_ "github.com/mattn/go-sqlite3"
)

type client struct {
	db       *sqlx.DB
	rootPath string
}

func NewClient(db *sqlx.DB, serveFrom string) (*client, error) {

	return &client{
		db:       db,
		rootPath: serveFrom,
	}, nil
}

func (c client) Fetch(key string, code string) ([]byte, error) {
	if key == "" {
		return nil, errors.New("cannot fetch null key, bad request")
	}

	log.Printf("fetching key=%s, code=%s\n", key, code)
	a := Access{}
	err := c.db.Get(&a, "SELECT * from access_keys WHERE key=$1", key)
	if err != nil {
		return nil, err
	}

	ok, err := a.Can(code)
	if !ok {
		return nil, err
	}

	return c.getFileContent(a.Path)
}

func (c client) getFileContent(path string) ([]byte, error) {
	fullPath := fmt.Sprintf("%s%s", c.rootPath, path)
	log.Printf("fetching file: %s\n", fullPath)
	return ioutil.ReadFile(fullPath)
}
