package database

import (
	"bytes"
	"database/sql"
	"errors"
	"fmt"
	"time"
)

type Vlan struct {
	ID       sql.NullInt64
	Name     sql.NullString
	Vlan     sql.NullInt64
	Site     sql.NullInt64
	SiteName sql.NullString
	Comment  sql.NullString
	Created  time.Time
	Changed  time.Time
}

type VLANQuery struct {
	db          *DBS
	vlan        Vlan
	idset       bool
	nameset     bool
	vlanset     bool
	siteset     bool
	sitenameset bool
	commentset  bool
}

func (db *DBS) InitiateVLANS() error {
	query := `CREATE TABLE IF NOT EXISTS vlans (
		id INTEGER PRIMARY KEY ASC,
		name VARCHAR(255),
		vlan SMALLINT NOT NULL,
		comment TEXT,
		site INTEGER NOT NULL,
		created TIMESTAMP NOT NULL,
		changed TIMESTAMP,
		FOREIGN KEY (site) REFERENCES sites(id)
	)`
	query = fmt.Sprintf(query)

	_, err := db.db.Exec(query)

	if err != nil {
		return err
	}
	return nil
}

func (db *DBS) GetVLANByName(site string) (int64, error) {
	query := `SELECT id FROM vlans WHERE name = ?`

	stmt, err := db.db.Prepare(query)
	if err != nil {
		return 0, err
	}

	var vlanID int64

	err = stmt.QueryRow(site).Scan(&vlanID)
	if err != nil {
		return 0, err
	}

	stmt.Close()

	return vlanID, nil
}

func (db *DBS) NewVLANQuery() *Vlan {
	q := &Vlan{
		ID:       sql.NullInt64{Valid: false, Int64: 0},
		Name:     sql.NullString{Valid: false, String: ""},
		Vlan:     sql.NullInt64{Valid: false, Int64: 0},
		Site:     sql.NullInt64{Valid: false, Int64: 0},
		SiteName: sql.NullString{Valid: false, String: ""},
		Comment:  sql.NullString{Valid: false, String: ""},
		Created:  time.Time{},
		Changed:  time.Time{},
	}
	return q
}

func (q *Vlan) SetID(id int64) error {
	return q.ID.Scan(id)
}

func (q *Vlan) SetVLAN(vlan int) error {
	if vlan < 0 || vlan > 4095 {
		return errors.New("Vlan.SetVLAN: vlans has to be between 0 and 4095")
	}
	return q.Vlan.Scan(vlan)
}

func (q *Vlan) SetName(name string) error {
	return q.Name.Scan(name)
}

func (q *Vlan) SetSite(site int) error {
	return q.Site.Scan(site)
}

func (q *Vlan) SetSiteName(sitename string) error {
	return q.SiteName.Scan(sitename)
}

func (q *Vlan) Query(db *DBS) ([]Vlan, error) {
	var result []Vlan
	if !q.ID.Valid && !q.Vlan.Valid && !q.Site.Valid && !q.SiteName.Valid {
		return result, errors.New("no query set")
	}
	var args []interface{}
	var query bytes.Buffer
	set := false
	query.WriteString("SELECT id, name, vlan, site, comment, created, changed FROM vlans WHERE ")
	if q.ID.Valid {
		args = append(args, q.ID)
		query.WriteString("id = ?")
		set = true
	}
	if q.Vlan.Valid {
		if set {
			query.WriteString(" AND ")
		}
		args = append(args, q.Vlan)
		query.WriteString("vlan = ?")
	}
	if q.Site.Valid {
		if set {
			query.WriteString(" AND ")
		}
		args = append(args, q.Site)
		query.WriteString("site = ?")
	}

	stmt, err := db.db.Prepare(query.String())
	if err != nil {
		return result, err
	}

	rows, err := stmt.Query(args)
	if err != nil {
		return result, err
	}

	var id sql.NullInt64
	var name sql.NullString
	var vlan sql.NullInt64
	var site sql.NullInt64
	var comment sql.NullString
	var created time.Time
	var changed time.Time

	for rows.Next() {
		err = rows.Scan(&id, &name, &vlan, &site, &comment, &created, &changed)
		if err != nil {
			return result, err
		}
		result = append(result, Vlan{ID: id, Name: name, Site: site, Comment: comment, Created: created, Changed: changed})
	}

	return result, nil
}

func (q *Vlan) Add(db *DBS) error {
	if !q.Vlan.Valid {
		return errors.New("vlan required")
	}

	query := `INSERT INTO vlans (vlan, name, site, comment, created, changed)
	VALUES (?, ?, ?, ?, current_timestamp, current_timestamp)`

	stmt, err := db.db.Prepare(query)
	if err != nil {
		return err
	}

	_, err = stmt.Exec(q.Vlan, q.Name, q.Site, q.Comment)
	if err != nil {
		return err
	}
	return nil
}

func (q *Vlan) PopulateSiteID(db *DBS) error {
	if !q.SiteName.Valid {
		return errors.New("no site name entered")
	}

	query := "SELECT id FROM sites WHERE name = ?"

	stmt, err := db.db.Prepare(query)
	if err != nil {
		return err
	}

	rows, err := stmt.Query(q.SiteName)
	if err != nil {
		return err
	}

	nrows := 0
	site := sql.NullInt64{}
	for rows.Next() {
		if nrows > 0 {
			return errors.New("more than one site returned")
		}
		if err = rows.Scan(&site); err != nil {
			return err
		}
		nrows++
	}
	if nrows == 0 {
		return errors.New(fmt.Sprintf("no site with name %s found", q.SiteName.String))
	}
	if err = q.Site.Scan(site.Int64); err != nil {
		return err
	}
	return nil
}
