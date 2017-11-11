package database

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
