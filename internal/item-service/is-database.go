package items

import (
	"database/sql"
	"github.com/JohannesF99/bkc-fundus-management/pkg/models"
	_ "github.com/go-sql-driver/mysql"
	"strconv"
	"time"
)

type DB struct {
	db *sql.DB
}

func connect() (DB, error) {
	db, err := sql.Open("mysql", "root:2678@/bkc?parseTime=true")
	if err != nil {
		return DB{}, models.Error{
			Details: err.Error(),
			Path:    "Item Service - Database Connection Failed",
			Object:  "",
			Time:    time.Now(),
		}
	}
	return DB{db}, nil
}

func (db DB) getAllItemsFromDB() ([]models.Item, error) {
	var items []models.Item
	tx, err := db.db.Begin()
	defer tx.Commit()
	if err != nil {
		return nil, models.Error{
			Details: err.Error(),
			Path:    "Item Service - getAllItemsFromDB()",
			Object:  "",
			Time:    time.Now(),
		}
	}
	rows, err := tx.Query("SELECT * FROM bkc.items")
	defer rows.Close()
	if err != nil {
		return nil, models.Error{
			Details: err.Error(),
			Path:    "Item Service - getAllItemsFromDB()",
			Object:  "",
			Time:    time.Now(),
		}
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
			return nil, models.Error{
				Details: err.Error(),
				Path:    "Item Service - getAllItemsFromDB()",
				Object:  "",
				Time:    time.Now(),
			}
		}
		items = append(items, newItem)
	}
	return items, nil
}

func (db DB) getItemWithIdFromDB(itemId int) (models.Item, error) {
	tx, err := db.db.Begin()
	defer tx.Commit()
	if err != nil {
		return models.Item{}, models.Error{
			Details: err.Error(),
			Path:    "Item Service - getItemWithIdFromDB()",
			Object:  "",
			Time:    time.Now(),
		}
	}
	rows, err := tx.Query("SELECT * FROM bkc.items WHERE id=?", itemId)
	defer rows.Close()
	if err != nil {
		return models.Item{}, models.Error{
			Details: err.Error(),
			Path:    "Item Service - getItemWithIdFromDB()",
			Object:  "",
			Time:    time.Now(),
		}
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
		return models.Item{}, models.Error{
			Details: err.Error(),
			Path:    "Item Service - getItemWithIdFromDB()",
			Object:  "",
			Time:    time.Now(),
		}
	}
	return newItem, nil
}

func (db DB) insertNewItemToDB(item models.Item) (int64, error) {
	tx, err := db.db.Begin()
	defer tx.Commit()
	if err != nil {
		_ = tx.Rollback()
		return -1, models.Error{
			Details: err.Error(),
			Path:    "Item Service - insertNewItemToDB()",
			Object:  "",
			Time:    time.Now(),
		}
	}
	rows, err := tx.Exec(
		"INSERT INTO bkc.items(name, capacity, description) VALUES (?,?,?)",
		item.Name, item.Capacity, item.Description)
	if err != nil {
		_ = tx.Rollback()
		return -1, models.Error{
			Details: err.Error(),
			Path:    "Item Service - insertNewItemToDB()",
			Object:  "",
			Time:    time.Now(),
		}
	}
	id, err := rows.LastInsertId()
	if err != nil {
		_ = tx.Rollback()
		return -1, models.Error{
			Details: err.Error(),
			Path:    "Item Service - insertNewItemToDB()",
			Object:  "",
			Time:    time.Now(),
		}
	}
	return id, nil
}

func (db DB) updateItemAvailabilityInDB(itemId int, diff int) error {
	tx, err := db.db.Begin()
	defer tx.Commit()
	if err != nil {
		_ = tx.Rollback()
		return models.Error{
			Details: err.Error(),
			Path:    "Item Service - updateItemAvailabilityInDB()",
			Object:  "",
			Time:    time.Now(),
		}
	}
	rows, err := tx.Exec("UPDATE bkc.items SET available=available+? WHERE id=?",
		diff, itemId)
	if err != nil {
		_ = tx.Rollback()
		return models.Error{
			Details: err.Error(),
			Path:    "Item Service - updateItemAvailabilityInDB()",
			Object:  "",
			Time:    time.Now(),
		}
	}
	affected, err := rows.RowsAffected()
	if err != nil {
		_ = tx.Rollback()
		return models.Error{
			Details: err.Error(),
			Path:    "Item Service - updateItemAvailabilityInDB()",
			Object:  "",
			Time:    time.Now(),
		}
	}
	if affected != 1 {
		_ = tx.Rollback()
		return models.Error{
			Details: "Rows affected: " + strconv.Itoa(int(affected)),
			Path:    "Item Service - updateItemAvailabilityInDB()",
			Object:  "",
			Time:    time.Now(),
		}
	}
	return nil
}

func (db DB) deleteItemFromDB(itemId int) error {
	tx, err := db.db.Begin()
	defer tx.Commit()
	if err != nil {
		_ = tx.Rollback()
		return models.Error{
			Details: err.Error(),
			Path:    "Item Service - deleteItemFromDB()",
			Object:  "",
			Time:    time.Now(),
		}
	}
	rows, err := tx.Exec("DELETE FROM bkc.items WHERE id=?", itemId)
	if err != nil {
		_ = tx.Rollback()
		return models.Error{
			Details: err.Error(),
			Path:    "Item Service - deleteItemFromDB()",
			Object:  "",
			Time:    time.Now(),
		}
	}
	affected, err := rows.RowsAffected()
	if err != nil {
		_ = tx.Rollback()
		return models.Error{
			Details: err.Error(),
			Path:    "Item Service - deleteItemFromDB()",
			Object:  "",
			Time:    time.Now(),
		}
	}
	if affected != 1 {
		_ = tx.Rollback()
		return models.Error{
			Details: "Rows affected: " + strconv.Itoa(int(affected)),
			Path:    "Item Service - deleteItemFromDB()",
			Object:  "",
			Time:    time.Now(),
		}
	}
	return nil
}

func (db DB) updateCapacityInDB(itemId int, diff int) error {
	tx, err := db.db.Begin()
	defer tx.Commit()
	if err != nil {
		_ = tx.Rollback()
		return models.Error{
			Details: err.Error(),
			Path:    "Item Service - updateCapacityInDB()",
			Object:  "",
			Time:    time.Now(),
		}
	}
	rows, err := tx.Exec("UPDATE bkc.items SET capacity=capacity+? WHERE id=?",
		diff, itemId)
	if err != nil {
		_ = tx.Rollback()
		return models.Error{
			Details: err.Error(),
			Path:    "Item Service - updateCapacityInDB()",
			Object:  "",
			Time:    time.Now(),
		}
	}
	affected, err := rows.RowsAffected()
	if err != nil {
		_ = tx.Rollback()
		return models.Error{
			Details: err.Error(),
			Path:    "Item Service - updateCapacityInDB()",
			Object:  "",
			Time:    time.Now(),
		}
	}
	if affected != 1 {
		_ = tx.Rollback()
		return models.Error{
			Details: "Rows affected: " + strconv.Itoa(int(affected)),
			Path:    "Item Service - updateCapacityInDB()",
			Object:  "",
			Time:    time.Now(),
		}
	}
	return nil
}
