package store

func (c Client) Migrate() error {
	_, _ = c.db.Exec(`CREATE TABLE access_keys(path text, key text, user_code text, display_name text, until timestamp, created timestamp);`)
	_, _ = c.db.Exec(`CREATE TABLE files(id text, path text, created timestamp, tag text);`)
	return nil
}
