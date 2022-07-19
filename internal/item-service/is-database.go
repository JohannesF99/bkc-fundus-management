package items

import (
	"database/sql"
	"github.com/JohannesF99/bkc-fundus-management/pkg/models"
	_ "github.com/go-sql-driver/mysql"
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

func (db DB) getAllItemsFromDB() ([]models.Item, error) {
	var items []models.Item
	tx, err := db.db.Begin()
	defer tx.Commit()
	if err != nil {
		return nil, err
	}
	rows, err := tx.Query("SELECT * FROM bkc.items")
	defer rows.Close()
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var newItem models.Item
		err := rows.Scan(
			&newItem.Id,
			&newItem.Name,
			&newItem.Capacity,
			&newItem.Availability,
			&newItem.Description,
			&newItem.Created,
			&newItem.Modified)
		if err != nil {
			return nil, err
		}
		items = append(items, newItem)
	}
	return items, nil
}

func (db DB) getItemWithIdFromDB(itemId int) (models.Item, error) {
	tx, err := db.db.Begin()
	defer tx.Commit()
	if err != nil {
		return models.Item{}, err
	}
	rows, err := tx.Query("SELECT * FROM bkc.items WHERE id=?", itemId)
	defer rows.Close()
	if err != nil {
		return models.Item{}, err
	}
	rows.Next()
	var newItem models.Item
	err = rows.Scan(
		&newItem.Id,
		&newItem.Name,
		&newItem.Capacity,
		&newItem.Availability,
		&newItem.Description,
		&newItem.Created,
		&newItem.Modified)
	if err != nil {
		return models.Item{}, err
	}
	return newItem, nil
}

func (db DB) insertNewItemToDB(item models.Item) (int64, error) {
	tx, err := db.db.Begin()
	defer tx.Commit()
	if err != nil {
		_ = tx.Rollback()
		return -1, err
	}
	rows, err := tx.Exec(
		"INSERT INTO bkc.items(name, capacity, description) VALUES (?,?,?)",
		item.Name, item.Capacity, item.Description)
	if err != nil {
		_ = tx.Rollback()
		return -1, err
	}
	id, err := rows.LastInsertId()
	if err != nil {
		_ = tx.Rollback()
		return -1, err
	}
	return id, nil
}

func (db DB) updateItemAvailabilityInDB(itemId int, diff int) (int, error) {
	tx, err := db.db.Begin()
	defer tx.Commit()
	if err != nil {
		_ = tx.Rollback()
		return -1, err
	}
	rows, err := tx.Exec("UPDATE bkc.items SET available=available+? WHERE id=?",
		diff, itemId)
	if err != nil {
		_ = tx.Rollback()
		return -1, err
	}
	affected, err := rows.RowsAffected()
	if err != nil || affected != 1 {
		_ = tx.Rollback()
		return 0, err
	}
	return itemId, nil
}

func (db DB) deleteItemFromDB(itemId int) error {
	tx, err := db.db.Begin()
	defer tx.Commit()
	if err != nil {
		_ = tx.Rollback()
		return err
	}
	rows, err := tx.Exec("DELETE FROM bkc.items WHERE id=?", itemId)
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
