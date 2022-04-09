package data

import "github.com/jmoiron/sqlx"

func Migrate(db *sqlx.DB) error {
	_, _ = db.Exec(`CREATE TABLE access_keys(path text, key text, user_code text, display_name text, until timestamp, created timestamp);`)
	_, _ = db.Exec(`CREATE TABLE files(id text, path text, created timestamp, tag text);`)
	return nil
}
