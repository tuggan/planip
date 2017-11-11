package database

import (
	"database/sql"
	"time"
)

type IP struct {
	ID       int64
	Name     string
	IP       string
	Site     sql.NullInt64
	SiteName sql.NullString
	Vlan     sql.NullInt64
	VlanName sql.NullString
	Netmask  sql.NullInt64
	Created  time.Time
	Changed  time.Time
}

func (db *DBS) InitiateIPs() error {
	query := `CREATE TABLE IF NOT EXISTS ips (
		id INTEGER PRIMARY KEY ASC,
		site INTEGER,
		vlan INTEGER,
		name VARCHAR(255),
		ip VARCHAR(15) NOT NULL,
		netmask SMALLINT,
		created TIMESTAMP NOT NULL,
		changed TIMESTAMP,
		FOREIGN KEY (site) REFERENCES sites(id),
		FOREIGN KEY (vlan) REFERENCES vlans(id)
	)`

	_, err := db.db.Exec(query)
	if err != nil {
		return err
	}
	return nil
}

func (db *DBS) AddIP(ip IP) error {

	var err error

	if ip.SiteName.Valid {
		id, err := db.GetSiteID(ip.SiteName.String)
		if err != nil {
			return err
		}
		err = ip.Site.Scan(id)
		if err != nil {
			return err
		}
	}

	query := `INSERT INTO ips (name, ip, netmask, site, vlan, created, changed) VALUES (
		?, ?, ?, ?, ?, current_timestamp, current_timestamp)`
	stmt, err := db.db.Prepare(query)
	if err != nil {
		return err
	}

	_, err = stmt.Exec(ip.Name, ip.IP, ip.Netmask, ip.Site, ip.Vlan)
	if err != nil {
		return err
	}

	stmt.Close()

	return nil
}
