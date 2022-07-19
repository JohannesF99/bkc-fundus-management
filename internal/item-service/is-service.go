package items

import (
	"github.com/JohannesF99/bkc-fundus-management/pkg/models"
)

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
	id, err := db.updateItemAvailabilityInDB(itemId, diff)
	if err != nil {
		return -1, err
	}
	return id, nil
}

func deleteItem(itemId int) error {
	db, err := connect()
	if err != nil {
		return err
	}
	err = db.deleteItemFromDB(itemId)
	if err != nil {
		return err
	}
	return nil
}
