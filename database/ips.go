package database

func (db *dbs) InitiateIPs() error {
	query := `CREATE TABLE IF NOT EXISTS ips (
		id INTEGER PRIMARY KEY ASC,
		name VARCHAR(255),
		ip VARCHAR(15) NOT NULL,
		netmask SMALLINT,
		created TIMESTAMP NOT NULL,
		changed TIMESTAMP
	)`

	_, err := db.db.Exec(query)
	if err != nil {
		return err
	}
	return nil
}
