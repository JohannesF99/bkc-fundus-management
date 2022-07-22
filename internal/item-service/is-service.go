package items

import (
	"github.com/JohannesF99/bkc-fundus-management/pkg/models"
	"time"
)

func deleteItem(itemId int) (models.Item, error) {
	db, err := connect()
	if err != nil {
		return models.Item{}, err
	}
	item, err := db.getItemWithIdFromDB(itemId)
	if err != nil {
		return models.Item{}, err
	}
	if item.Availability != item.Capacity {
		return models.Item{}, models.Error{
			Details: "You can only delete this Item, when every borrowed piece has been returned.",
			Path:    "Item Service - deleteItem()",
			Object:  item.String(),
			Time:    time.Now(),
		}
	}
	err = db.deleteItemFromDB(itemId)
	if err != nil {
		return models.Item{}, err
	}
	return item, nil
}

func getAllItems() ([]models.Item, error) {
	db, err := connect()
	if err != nil {
		return nil, err
	}
	items, err := db.getAllItemsFromDB()
	if err != nil {
		return nil, err
	}
	return items, nil
}

func getItemWithId(itemId int) (models.Item, error) {
	db, err := connect()
	if err != nil {
		return models.Item{}, err
	}
	item, err := db.getItemWithIdFromDB(itemId)
	if err != nil {
		return models.Item{}, err
	}
	return item, nil
}

func insertNewItem(newItem models.NewItemInfos) (int64, error) {
	db, err := connect()
	if err != nil {
		return -1, err
	}
	id, err := db.insertNewItemToDB(models.Item{
		Name:        newItem.Name,
		Capacity:    newItem.Capacity,
		Description: newItem.Description,
	})
	if err != nil {
		return -1, err
	}
	return id, nil
}

func updateItemAvailability(itemId int, diff int) (int, error) {
	db, err := connect()
	if err != nil {
		return -1, err
	}
	err = db.updateItemAvailabilityInDB(itemId, diff)
	if err != nil {
		return -1, err
	}
	return itemId, nil
}
