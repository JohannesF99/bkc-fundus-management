package entry

import "github.com/JohannesF99/bkc-fundus-management/pkg/models"

func getAllEntries() ([]models.EntryInfo, error) {
	db, err := connect()
	if err != nil {
		return nil, err
	}
	allEntries, err := db.getAllEntriesFromDB()
	if err != nil {
		return nil, err
	}
	return allEntries, nil
}

func createNewEntry(newEntry models.NewEntry) (models.EntryInfo, error) {
	db, err := connect()
	if err != nil {
		return models.EntryInfo{}, err
	}
	entryId, err := db.addEntryToDB(newEntry)
	if err != nil {
		return models.EntryInfo{}, err
	}
	entry, err := db.getEntriesForEntryIdFromDB(int(entryId))
	if err != nil {
		return models.EntryInfo{}, err
	}
	return entry, nil
}

func updateEntry(entryId int, diff int) (models.EntryInfo, error) {
	db, err := connect()
	if err != nil {
		return models.EntryInfo{}, err
	}
	entryId, err = db.updateEntryInDB(entryId, diff)
	if err != nil {
		return models.EntryInfo{}, err
	}
	entry, err := db.getEntriesForEntryIdFromDB(int(entryId))
	if err != nil {
		return models.EntryInfo{}, err
	}
	return entry, nil
}

func deleteEntry(entryId int) (models.EntryInfo, error) {
	db, err := connect()
	if err != nil {
		return models.EntryInfo{}, err
	}
	entryInfo, err := db.getEntriesForEntryIdFromDB(entryId)
	if err != nil {
		return models.EntryInfo{}, err
	}
	err = db.deleteEntryFromDB(entryId)
	if err != nil {
		return models.EntryInfo{}, err
	}
	return entryInfo, nil
}

func getAllEntriesByMemberId(memberId int) ([]models.EntryInfo, error) {
	db, err := connect()
	if err != nil {
		return nil, err
	}
	allEntries, err := db.getEntriesForMemberIdFromDB(memberId)
	if err != nil {
		return nil, err
	}
	return allEntries, nil
}

func getAllEntriesByItemId(itemId int) ([]models.EntryInfo, error) {
	db, err := connect()
	if err != nil {
		return nil, err
	}
	allEntries, err := db.getEntriesForItemIdFromDB(itemId)
	if err != nil {
		return nil, err
	}
	return allEntries, nil
}
