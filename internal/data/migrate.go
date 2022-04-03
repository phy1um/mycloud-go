package data

import "github.com/jmoiron/sqlx"

func Migrate(db *sqlx.DB) error {
	_, err := db.Exec(`CREATE TABLE access_keys(path text, key text, user_code text, until timestamp, created timestamp);`)
	return err
}
