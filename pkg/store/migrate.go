package store

func (c Client) Migrate() error {
	_, _ = c.db.Exec(`CREATE TABLE access_keys(file_id text, key text, user_code text, display_name text, until timestamp, created timestamp);`)
	_, _ = c.db.Exec(`CREATE TABLE files(id text, path text, name text, created timestamp);`)
	_, _ = c.db.Exec(`CREATE TABLE tags (id text, value text);`)
	return nil
}
