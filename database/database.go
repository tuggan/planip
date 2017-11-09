package database

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type dbs struct {
	db *sql.DB
}

type Site struct {
	ID      int64
	Name    string
	Comment sql.NullString
	Created time.Time
	Changed time.Time
}

func Open() (*dbs, error) {
	var db dbs
	var err error
	db.db, err = sql.Open("sqlite3", "./foo.db")
	if err != nil {
		return nil, err
	}
	return &db, nil
}

func (db *dbs) Close() error {
	db.db.Close()
	return nil
}

func (db *dbs) InitiateDatabase() error {
	name := "planip"
	query := fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s", name)
	_, err := db.db.Exec(query)

	if err != nil {
		return err
	}
	return nil
}

func (db *dbs) InitiateSites() error {
	tableName := "sites"
	id := "id INTEGER PRIMARY KEY ASC"
	name := "name VARCHAR(255) NOT NULL"
	comment := "comment VARCHAR(255)"
	created := "created TIMESTAMP NOT NULL"
	changed := "changed TIMESTAMP"
	query := `CREATE TABLE IF NOT EXISTS %s (%s, %s, %s, %s, %s)`
	query = fmt.Sprintf(query, tableName, id, name, comment, created, changed)

	_, err := db.db.Exec(query)

	if err != nil {
		return err
	}
	return nil
}

func (db *dbs) InitiateVLANS() error {
	tableName := "VLANS"
	id := "id INTEGER PRIMARY KEY ASC"
	name := "name VARCHAR(255) NOT NULL"
	vlan := "vlan SMALLINT NOT NULL"
	created := "created TIMESTAMP NOT NULL"
	changed := "changed TIMSTAMP"
	site := "site INT NOT NULL, FOREIGN KEY (site) REFERENCES sites(id)"
	query := `CREATE TABLE IF NOT EXISTS %s (%s, %s, %s, %s, %s, %s)`
	query = fmt.Sprintf(query, tableName, id, name, vlan, created, changed, site)

	_, err := db.db.Exec(query)

	if err != nil {
		return err
	}
	return nil
}

func (db *dbs) Init() error {
	err := db.InitiateSites()
	if err != nil {
		return err
	}

	err = db.InitiateVLANS()
	if err != nil {
		return err
	}
	return nil
}

func (db *dbs) AddSite(site string) error {
	tableName := "sites"
	//name := fmt.Sprintf("name=%s", site)
	//created := "CURRENT_TIMESTAMP"
	//changed := "CURRENT_TIMESTAMP"

	query := `INSERT INTO %s (name, created, changed) values (?, current_timestamp, current_timestamp)`
	query = fmt.Sprintf(query, tableName)

	stmt, err := db.db.Prepare(query)
	if err != nil {
		return err
	}

	_, err = stmt.Exec(site)
	if err != nil {
		return err
	}

	stmt.Close()

	return nil
}

func (db *dbs) GetSites() ([]Site, error) {
	var sites []Site
	tableName := "sites"

	query := fmt.Sprintf("SELECT * FROM %s", tableName)

	rows, err := db.db.Query(query)
	if err != nil {
		return nil, err
	}
	var id int64
	var name string
	var comment sql.NullString
	var created time.Time
	var changed time.Time

	for rows.Next() {
		err = rows.Scan(&id, &name, &comment, &created, &changed)
		if err != nil {
			return nil, err
		}
		sites = append(sites, Site{id, name, comment, created, changed})
	}

	return sites, nil
}
