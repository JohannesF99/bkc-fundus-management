package entry

import (
	"encoding/json"
	"github.com/JohannesF99/bkc-fundus-management/pkg/constants"
	"github.com/JohannesF99/bkc-fundus-management/pkg/models"
	"net/http"
	"strconv"
)

const ItemService = constants.ItemService
const MemberService = constants.MemberService

func getAllEntries() ([]models.Entry, error) {
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

func createNewEntry(newEntry models.NewEntryInfos) (models.Entry, error) {
	db, err := connect()
	if err != nil {
		return models.Entry{}, err
	}
	entryId, err := db.addEntryToDB(newEntry)
	if err != nil {
		return models.Entry{}, err
	}
	entry, err := db.getEntriesForEntryIdFromDB(int(entryId))
	if err != nil {
		return models.Entry{}, err
	}
	return entry, nil
}

func getEntryForEntryId(entryId int) (models.Entry, error) {
	db, err := connect()
	if err != nil {
		return models.Entry{}, err
	}
	entry, err := db.getEntriesForEntryIdFromDB(entryId)
	if err != nil {
		return models.Entry{}, err
	}
	return entry, nil
}

func updateEntry(entryId int, diff int) (models.Entry, error) {
	db, err := connect()
	if err != nil {
		return models.Entry{}, err
	}
	entryId, err = db.updateEntryInDB(entryId, diff)
	if err != nil {
		return models.Entry{}, err
	}
	entry, err := db.getEntriesForEntryIdFromDB(int(entryId))
	if err != nil {
		return models.Entry{}, err
	}
	return entry, nil
}

func deleteEntry(entryId int) (models.Entry, error) {
	db, err := connect()
	if err != nil {
		return models.Entry{}, err
	}
	entryInfo, err := db.getEntriesForEntryIdFromDB(entryId)
	if err != nil {
		return models.Entry{}, err
	}
	err = db.deleteEntryFromDB(entryId)
	if err != nil {
		return models.Entry{}, err
	}
	return entryInfo, nil
}

func getAllEntriesByMemberId(memberId int) ([]models.Export, error) {
	db, err := connect()
	if err != nil {
		return nil, err
	}
	allEntries, err := db.getEntriesForMemberIdFromDB(memberId)
	if err != nil {
		return nil, err
	}
	var itemInfosForMember []models.Export
	for _, entry := range allEntries {
		item, err := requestItemFromItemService(entry.ItemId)
		if err != nil {
			return nil, err
		}
		itemInfosForMember = append(itemInfosForMember, models.Export{
			Id:       item.Id,
			Name:     item.Name,
			Capacity: entry.Capacity,
			Date:     entry.Modified,
		})
	}
	return itemInfosForMember, nil
}

func requestItemFromItemService(itemId int) (models.Item, error) {
	resp, err := http.Get(ItemService + strconv.Itoa(itemId))
	if err != nil {
		return models.Item{}, err
	}
	defer resp.Body.Close()
	var item models.Item
	err = json.NewDecoder(resp.Body).Decode(&item)
	if err != nil {
		return models.Item{}, err
	}
	return item, nil
}

func getAllEntriesByItemId(itemId int) ([]models.Export, error) {
	db, err := connect()
	if err != nil {
		return nil, err
	}
	allEntries, err := db.getEntriesForItemIdFromDB(itemId)
	if err != nil {
		return nil, err
	}
	var memberInfosForItem []models.Export
	for _, entry := range allEntries {
		member, err := requestMemberFromMemberService(entry.MemberId)
		if err != nil {
			return nil, err
		}
		memberInfosForItem = append(memberInfosForItem, models.Export{
			Id:       member.Id,
			Name:     member.Name,
			Capacity: entry.Capacity,
			Date:     entry.Modified,
		})
	}
	return memberInfosForItem, nil
}

func requestMemberFromMemberService(memberId int) (models.Member, error) {
	resp, err := http.Get(MemberService + strconv.Itoa(memberId))
	if err != nil {
		return models.Member{}, err
	}
	defer resp.Body.Close()
	var member models.Member
	err = json.NewDecoder(resp.Body).Decode(&member)
	if err != nil {
		return models.Member{}, err
	}
	return member, nil
}

func getEntryForMemberIdAndItemId(memberId int, itemId int) (models.Entry, error) {
	db, err := connect()
	if err != nil {
		return models.Entry{}, err
	}
	entry, err := db.getEntryForMemberIdAndItemIdFromDB(memberId, itemId)
	if err != nil {
		return models.Entry{}, err
	}
	return entry, nil
}
