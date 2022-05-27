package fetch

import (
	"context"
	"fmt"
	"io/ioutil"

	"github.com/phy1um/mycloud-go/pkg/data"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"

	_ "github.com/mattn/go-sqlite3"
)

type client struct {
	db       *sqlx.DB
	rootPath string
}

type keyCrossFile struct {
	data.Access
	data.File
}

func NewClient(db *sqlx.DB, serveFrom string) (*client, error) {

	return &client{
		db:       db,
		rootPath: serveFrom,
	}, nil
}

func (c client) Fetch(ctx context.Context, key string, code string) ([]byte, error) {
	if key == "" {
		return nil, fmt.Errorf("no resource key provided")
	}

	log := log.Ctx(ctx).With().Str("key", key).Logger()
	log.Info().Msg("fetching resource")
	ctx = log.WithContext(ctx)

	a := keyCrossFile{}
	err := c.db.Get(&a, "SELECT * from access_keys JOIN files WHERE key=$1", key)
	if err != nil {
		return nil, errors.Wrap(err, "failed to query database for key")
	}

	ok, err := a.Can(code)
	if !ok {
		return nil, errors.Wrap(err, "access not granted to resource")
	}

	return c.getFileContent(ctx, a.Path)
}

func (c client) getFileContent(ctx context.Context, path string) ([]byte, error) {
	fullPath := fmt.Sprintf("%s%s", c.rootPath, path)
	logger := log.Ctx(ctx).With().Str("file-path", fullPath).Logger()
	logger.Info().Msg("read file from filesystem")
	return ioutil.ReadFile(fullPath)
}
