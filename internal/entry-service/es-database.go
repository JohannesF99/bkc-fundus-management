package entry

import (
	"database/sql"
	"github.com/JohannesF99/bkc-fundus-management/pkg/models"
	_ "github.com/go-sql-driver/mysql"
	"time"
)

type DB struct {
	db *sql.DB
}

func connect() (DB, error) {
	db, err := sql.Open("mysql", "root:2678@/bkc?parseTime=true")
	if err != nil {
		return DB{}, err
	}
	return DB{db}, nil
}

func (db DB) getAllEntriesFromDB() ([]models.Entry, error) {
	var allEntries []models.Entry
	tx, err := db.db.Begin()
	defer tx.Commit()
	if err != nil {
		return nil, err
	}
	rows, err := tx.Query("SELECT * FROM bkc.entries")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var newEntry models.Entry
		err := rows.Scan(
			&newEntry.Id,
			&newEntry.MemberId,
			&newEntry.ItemId,
			&newEntry.Capacity,
			&newEntry.Created,
			&newEntry.Modified)
		if err != nil {
			return nil, err
		}
		allEntries = append(allEntries, newEntry)
	}
	return allEntries, nil
}

func (db DB) addEntryToDB(newEntry models.NewEntryInfos) (int64, error) {
	tx, err := db.db.Begin()
	defer tx.Commit()
	if err != nil {
		_ = tx.Rollback()
		return -1, err
	}
	rows, err := tx.Exec(
		"INSERT INTO bkc.entries(member_id, item_id, capacity) VALUES (?,?,?)",
		newEntry.MemberId, newEntry.ItemId, newEntry.Capacity)
	if err != nil {
		_ = tx.Rollback()
		return -1, err
	}
	entryId, err := rows.LastInsertId()
	if err != nil {
		_ = tx.Rollback()
		return -1, err
	}
	return entryId, nil
}

func (db DB) updateEntryInDB(entryId int, diff int) (int, error) {
	tx, err := db.db.Begin()
	defer tx.Commit()
	if err != nil {
		_ = tx.Rollback()
		return -1, err
	}
	rows, err := tx.Exec(
		"UPDATE bkc.entries SET capacity=capacity+? WHERE id=?",
		diff, entryId)
	if err != nil {
		_ = tx.Rollback()
		return -1, err
	}
	affected, err := rows.RowsAffected()
	if err != nil || affected != 1 {
		_ = tx.Rollback()
		return -1, err
	}
	return entryId, nil
}

func (db DB) deleteEntryFromDB(entryId int) error {
	tx, err := db.db.Begin()
	defer tx.Commit()
	if err != nil {
		_ = tx.Rollback()
		return err
	}
	rows, err := tx.Exec(
		"DELETE FROM bkc.entries WHERE id=?",
		entryId)
	if err != nil {
		_ = tx.Rollback()
		return err
	}
	affected, err := rows.RowsAffected()
	if err != nil || affected != 1 {
		_ = tx.Rollback()
		return err
	}
	return nil
}

func (db DB) getEntriesForMemberIdFromDB(memberId int) ([]models.Entry, error) {
	var allEntries []models.Entry
	tx, err := db.db.Begin()
	defer tx.Commit()
	if err != nil {
		return nil, err
	}
	rows, err := tx.Query("SELECT * FROM bkc.entries WHERE member_id=?", memberId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var newEntry models.Entry
		err := rows.Scan(
			&newEntry.Id,
			&newEntry.MemberId,
			&newEntry.ItemId,
			&newEntry.Capacity,
			&newEntry.Created,
			&newEntry.Modified)
		if err != nil {
			return nil, err
		}
		allEntries = append(allEntries, newEntry)
	}
	return allEntries, nil
}

func (db DB) getEntriesForItemIdFromDB(itemId int) ([]models.Entry, error) {
	var allEntries []models.Entry
	tx, err := db.db.Begin()
	defer tx.Commit()
	if err != nil {
		return nil, err
	}
	rows, err := tx.Query("SELECT * FROM bkc.entries WHERE item_id=?", itemId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var newEntry models.Entry
		err := rows.Scan(
			&newEntry.Id,
			&newEntry.MemberId,
			&newEntry.ItemId,
			&newEntry.Capacity,
			&newEntry.Created,
			&newEntry.Modified)
		if err != nil {
			return nil, err
		}
		allEntries = append(allEntries, newEntry)
	}
	return allEntries, nil
}

func (db DB) getEntryForEntryIdFromDB(entryId int) (models.Entry, error) {
	tx, err := db.db.Begin()
	defer tx.Commit()
	if err != nil {
		return models.Entry{}, models.Error{
			Details: err.Error(),
			Path:    "ES-Database",
			Object:  "",
			Time:    time.Now(),
		}
	}
	rows, err := tx.Query("SELECT * FROM bkc.entries WHERE id=?", entryId)
	if err != nil {
		return models.Entry{}, models.Error{
			Details: err.Error(),
			Path:    "ES-Database",
			Object:  "",
			Time:    time.Now(),
		}
	}
	defer rows.Close()
	rows.Next()
	var newEntry models.Entry
	err = rows.Scan(
		&newEntry.Id,
		&newEntry.MemberId,
		&newEntry.ItemId,
		&newEntry.Capacity,
		&newEntry.Created,
		&newEntry.Modified)
	if err != nil {
		return models.Entry{}, models.Error{
			Details: err.Error(),
			Path:    "ES-Database",
			Object:  "",
			Time:    time.Now(),
		}
	}
	return newEntry, nil
}

func (db DB) getEntryForMemberIdAndItemIdFromDB(memberId int, itemId int) (models.Entry, error) {
	tx, err := db.db.Begin()
	defer tx.Commit()
	if err != nil {
		return models.Entry{}, err
	}
	rows, err := tx.Query("SELECT * FROM bkc.entries WHERE member_id=? AND item_id=?", memberId, itemId)
	if err != nil {
		return models.Entry{}, err
	}
	defer rows.Close()
	rows.Next()
	var newEntry models.Entry
	err = rows.Scan(
		&newEntry.Id,
		&newEntry.MemberId,
		&newEntry.ItemId,
		&newEntry.Capacity,
		&newEntry.Created,
		&newEntry.Modified)
	if err != nil {
		return models.Entry{}, err
	}
	return newEntry, nil
}
