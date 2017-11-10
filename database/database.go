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

type Vlan struct {
	ID       int64
	Name     string
	Vlan     int16
	Site     int64
	SiteName string
	Comment  sql.NullString
	Created  time.Time
	Changed  time.Time
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
	/*tableName := "sites"
	id := "id INTEGER PRIMARY KEY ASC"
	name := "name VARCHAR(255) NOT NULL UNIQUE"
	comment := "comment VARCHAR(255)"
	created := "created TIMESTAMP NOT NULL"
	changed := "changed TIMESTAMP"
	query := `CREATE TABLE IF NOT EXISTS %s (%s, %s, %s, %s, %s)`
	*/
	query := `CREATE TABLE IF NOT EXISTS sites (
		id INTEGER PRIMARY KEY ASC,
		name VARCHAR(255) NOT NULL,
		comment TEXT,
		created TIMESTAMP NOT NULL,
		changed TIMESTAMP,
		UNIQUE (name)
	)`
	//query = fmt.Sprintf(query, tableName, id, name, comment, created, changed)

	_, err := db.db.Exec(query)

	if err != nil {
		return err
	}
	return nil
}

func (db *dbs) InitiateVLANS() error {
	tableName := "vlans"
	id := "id INTEGER PRIMARY KEY ASC"
	name := "name VARCHAR(255)"
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

func (db *dbs) GetSiteID(site string) (int64, error) {
	query := `SELECT id FROM sites WHERE name = ?`

	stmt, err := db.db.Prepare(query)
	if err != nil {
		return 0, err
	}

	var siteID int64

	err = stmt.QueryRow(site).Scan(&siteID)
	if err != nil {
		return 0, err
	}

	stmt.Close()

	return siteID, nil
}

func (db *dbs) GetVLANS() ([]Vlan, error) {
	query := `SELECT vlans.vlan, vlans.name, sites.name FROM vlans INNER JOIN sites ON vlans.site = sites.id`

	rows, err := db.db.Query(query)
	if err != nil {
		return nil, err
	}
	var vlan int16
	var name string
	var sitename string

	var vlans []Vlan

	for rows.Next() {
		err = rows.Scan(&vlan, &name, &sitename)
		if err != nil {
			return nil, err
		}
		vlans = append(vlans, Vlan{Vlan: vlan, Name: name, SiteName: sitename})
	}

	return vlans, nil
}

func (db *dbs) AddVLAN(vlan Vlan) error {

	var err error

	vlan.Site, err = db.GetSiteID(vlan.SiteName)
	if err != nil {
		return err
	}

	query := `INSERT INTO vlans (name, vlan, created, changed, site) VALUES (?, ?, current_timestamp, current_timestamp, ?)`

	stmt, err := db.db.Prepare(query)
	if err != nil {
		return err
	}

	_, err = stmt.Exec(vlan.Name, vlan.Vlan, vlan.Site)
	if err != nil {
		return err
	}

	stmt.Close()

	return nil

}
