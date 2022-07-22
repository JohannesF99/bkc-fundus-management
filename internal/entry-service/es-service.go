package entry

import (
	"encoding/json"
	"github.com/JohannesF99/bkc-fundus-management/pkg/constants"
	"github.com/JohannesF99/bkc-fundus-management/pkg/models"
	"net/http"
	"strconv"
	"time"
)

const ItemService = constants.ItemService
const MemberService = constants.MemberService

func deleteEntry(entryId int) (models.Entry, error) {
	db, err := connect()
	if err != nil {
		return models.Entry{}, err
	}
	entry, err := db.getEntryForEntryIdFromDB(entryId)
	if err != nil {
		return models.Entry{}, err
	}
	err = db.deleteEntryFromDB(entryId)
	if err != nil {
		return models.Entry{}, err
	}
	_, err = changeItemAvailability(entry.ItemId, entry.Capacity)
	if err != nil {
		return models.Entry{}, err
	}
	_, err = changeMemberBorrowCount(entry.MemberId, entry.Capacity)
	if err != nil {
		return models.Entry{}, err
	}
	return entry, nil
}

func createNewEntryOrUpdate(newEntryInfo models.NewEntryInfos) (models.Entry, error) {
	db, err := connect()
	if err != nil {
		return models.Entry{}, err
	}
	_, err = doesMemberExist(newEntryInfo.MemberId)
	if err != nil {
		return models.Entry{}, err
	}
	item, err := doesItemExist(newEntryInfo.ItemId)
	if err != nil {
		return models.Entry{}, err
	}
	if item.Availability < newEntryInfo.Capacity {
		return models.Entry{}, models.Error{
			Details: "Availability is greater then Capacity",
			Path:    "Entry Service - createNewEntryOrUpdate()",
			Object:  item.String(),
			Time:    time.Now(),
		}
	}
	if existingEntry, err := db.getEntryForMemberIdAndItemIdFromDB(newEntryInfo.MemberId, newEntryInfo.ItemId); err == nil {
		entryId, err := db.updateEntryInDB(existingEntry.Id, newEntryInfo.Capacity)
		if err != nil {
			return models.Entry{}, err
		}
		updatedEntry, err := db.getEntryForEntryIdFromDB(entryId)
		if err != nil {
			return models.Entry{}, err
		}
		_, err = changeItemAvailability(updatedEntry.ItemId, -newEntryInfo.Capacity)
		if err != nil {
			return models.Entry{}, err
		}
		_, err = changeMemberBorrowCount(updatedEntry.MemberId, -newEntryInfo.Capacity)
		if err != nil {
			return models.Entry{}, err
		}
		return updatedEntry, nil
	} else {
		newEntry, err := createNewEntry(newEntryInfo)
		if err != nil {
			return models.Entry{}, err
		}
		_, err = changeItemAvailability(newEntry.ItemId, -newEntryInfo.Capacity)
		if err != nil {
			return models.Entry{}, err
		}
		_, err = changeMemberBorrowCount(newEntry.MemberId, -newEntryInfo.Capacity)
		if err != nil {
			return models.Entry{}, err
		}
		return newEntry, nil
	}
}

func changeItemAvailability(itemId int, returned int) (models.Item, error) {
	req, err := http.NewRequest(http.MethodPut,
		ItemService+strconv.Itoa(itemId)+
			"?returned="+strconv.Itoa(returned),
		nil)
	if err != nil {
		return models.Item{}, models.Error{
			Details: err.Error(),
			Path:    "Entry Service - changeItemAvailability()",
			Object:  "",
			Time:    time.Now(),
		}
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return models.Item{}, models.Error{
			Details: err.Error(),
			Path:    "Entry Service - changeItemAvailability",
			Object:  "",
			Time:    time.Now(),
		}
	}
	if resp.StatusCode != 200 {
		var apiError models.Error
		err = json.NewDecoder(resp.Body).Decode(&apiError)
		if err != nil {
			return models.Item{}, err
		}
		return models.Item{}, apiError
	}
	defer resp.Body.Close()
	var updatedItem models.Item
	err = json.NewDecoder(resp.Body).Decode(&updatedItem)
	if err != nil {
		return models.Item{}, models.Error{
			Details: err.Error(),
			Path:    "Entry Service - changeItemAvailability",
			Object:  "",
			Time:    time.Now(),
		}
	}
	return updatedItem, nil
}

func changeMemberBorrowCount(memberId int, returned int) (models.Member, error) {
	req, err := http.NewRequest(http.MethodPut,
		MemberService+strconv.Itoa(memberId)+
			"?returned="+strconv.Itoa(returned),
		nil)
	if err != nil {
		return models.Member{}, models.Error{
			Details: err.Error(),
			Path:    "Entry Service - changeMemberBorrowCount()",
			Object:  strconv.Itoa(memberId),
			Time:    time.Now(),
		}
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return models.Member{}, models.Error{
			Details: err.Error(),
			Path:    "Entry Service - changeMemberBorrowCount()",
			Object:  "",
			Time:    time.Now(),
		}
	}
	if resp.StatusCode != 200 {
		var apiError models.Error
		err = json.NewDecoder(resp.Body).Decode(&apiError)
		if err != nil {
			return models.Member{}, err
		}
		return models.Member{}, apiError
	}
	defer resp.Body.Close()
	var updatedMember models.Member
	err = json.NewDecoder(resp.Body).Decode(&updatedMember)
	if err != nil {
		return models.Member{}, models.Error{
			Details: err.Error(),
			Path:    "Entry Service - changeMemberBorrowCount()",
			Object:  strconv.Itoa(memberId),
			Time:    time.Now(),
		}
	}
	return updatedMember, nil
}

func doesMemberExist(memberId int) (models.Member, error) {
	resp, err := http.Get(MemberService + strconv.Itoa(memberId))
	if err != nil {
		return models.Member{}, models.Error{
			Details: err.Error(),
			Path:    "Entry Service - doesMemberExist()",
			Object:  "",
			Time:    time.Now(),
		}
	}
	if resp.StatusCode != 200 {
		var apiError models.Error
		err = json.NewDecoder(resp.Body).Decode(&apiError)
		if err != nil {
			return models.Member{}, err
		}
		return models.Member{}, apiError
	}
	defer resp.Body.Close()
	var existingMember models.Member
	err = json.NewDecoder(resp.Body).Decode(&existingMember)
	if err != nil {
		return models.Member{}, models.Error{
			Details: err.Error(),
			Path:    "Entry Service - doesMemberExist()",
			Object:  "",
			Time:    time.Now(),
		}
	}
	return existingMember, nil
}

func doesItemExist(itemId int) (models.Item, error) {
	resp, err := http.Get(ItemService + strconv.Itoa(itemId))
	if err != nil {
		return models.Item{}, models.Error{
			Details: err.Error(),
			Path:    "Entry Service - doesItemExist()",
			Object:  "",
			Time:    time.Now(),
		}
	}
	if resp.StatusCode != 200 {
		var apiError models.Error
		err = json.NewDecoder(resp.Body).Decode(&apiError)
		if err != nil {
			return models.Item{}, err
		}
		return models.Item{}, apiError
	}
	defer resp.Body.Close()
	var itemInfo models.Item
	err = json.NewDecoder(resp.Body).Decode(&itemInfo)
	if err != nil {
		return models.Item{}, models.Error{
			Details: err.Error(),
			Path:    "Entry Service - doesItemExist()",
			Object:  "",
			Time:    time.Now(),
		}
	}
	return itemInfo, nil
}

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
	entry, err := db.getEntryForEntryIdFromDB(int(entryId))
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
	entry, err := db.getEntryForEntryIdFromDB(entryId)
	if err != nil {
		return models.Entry{}, err
	}
	return entry, nil
}

func updateEntry(entryId int, returned int) (models.Entry, error) {
	db, err := connect()
	if err != nil {
		return models.Entry{}, err
	}
	entry, err := getEntryForEntryId(entryId)
	if err != nil {
		return models.Entry{}, err
	}
	item, err := doesItemExist(entry.ItemId)
	if err != nil {
		return models.Entry{}, err
	}
	if item.Availability-returned > item.Capacity {
		return models.Entry{}, models.Error{
			Details: "You tried to return more, than you borrowed",
			Path:    "Entry Service - updateEntry()",
			Object:  item.String(),
			Time:    time.Now(),
		}
	}
	_, err = doesMemberExist(entry.MemberId)
	if err != nil {
		return models.Entry{}, err
	}
	if entry.Capacity < returned {
		return models.Entry{}, models.Error{
			Details: "You tried to return more, than you had initially borrowed",
			Path:    "Entry Service - updateEntry()",
			Object:  entry.String(),
			Time:    time.Now(),
		}
	}
	if entry.Capacity == returned {
		entry, err := deleteEntry(entryId)
		if err != nil {
			return models.Entry{}, err
		}
		return entry, nil
	}
	_, err = changeItemAvailability(entry.ItemId, returned)
	if err != nil {
		return models.Entry{}, err
	}
	_, err = changeMemberBorrowCount(entry.MemberId, returned)
	if err != nil {
		return models.Entry{}, err
	}
	entryId, err = db.updateEntryInDB(entryId, -returned)
	if err != nil {
		return models.Entry{}, err
	}
	entry, err = db.getEntryForEntryIdFromDB(entryId)
	if err != nil {
		return models.Entry{}, err
	}
	return entry, nil
}

func getAllEntriesByMemberId(memberId int) ([]models.Export, error) {
	db, err := connect()
	if err != nil {
		return nil, err
	}
	_, err = doesMemberExist(memberId)
	if err != nil {
		return nil, err
	}
	allEntries, err := db.getEntriesForMemberIdFromDB(memberId)
	if err != nil {
		return nil, err
	}
	itemInfosForMember := []models.Export{}
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
		return models.Item{}, models.Error{
			Details: err.Error(),
			Path:    "Entry Service - requestItemFromItemService",
			Object:  "",
			Time:    time.Now(),
		}
	}
	if resp.StatusCode != 200 {
		var apiError models.Error
		err = json.NewDecoder(resp.Body).Decode(&apiError)
		if err != nil {
			return models.Item{}, err
		}
		return models.Item{}, apiError
	}
	defer resp.Body.Close()
	var item models.Item
	err = json.NewDecoder(resp.Body).Decode(&item)
	if err != nil {
		return models.Item{}, models.Error{
			Details: err.Error(),
			Path:    "Entry Service - requestItemFromItemService()",
			Object:  "",
			Time:    time.Now(),
		}
	}
	return item, nil
}

func getAllEntriesByItemId(itemId int) ([]models.Export, error) {
	db, err := connect()
	if err != nil {
		return nil, err
	}
	_, err = doesItemExist(itemId)
	if err != nil {
		return nil, err
	}
	allEntries, err := db.getEntriesForItemIdFromDB(itemId)
	if err != nil {
		return nil, err
	}
	memberInfosForItem := []models.Export{}
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
		return models.Member{}, models.Error{
			Details: err.Error(),
			Path:    "Entry Service - requestMemberFromMemberService()",
			Object:  "",
			Time:    time.Now(),
		}
	}
	if resp.StatusCode != 200 {
		var apiError models.Error
		err = json.NewDecoder(resp.Body).Decode(&apiError)
		if err != nil {
			return models.Member{}, err
		}
		return models.Member{}, apiError
	}
	defer resp.Body.Close()
	var member models.Member
	err = json.NewDecoder(resp.Body).Decode(&member)
	if err != nil {
		return models.Member{}, models.Error{
			Details: err.Error(),
			Path:    "Entry Service - requestMemberFromMemberService()",
			Object:  "",
			Time:    time.Now(),
		}
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

func borrowedItemLost(entryId int, diff int) (models.Item, error) {
	db, err := connect()
	if err != nil {
		return models.Item{}, err
	}
	entry, err := db.getEntryForEntryIdFromDB(entryId)
	if err != nil {
		return models.Item{}, err
	}
	item, err := doesItemExist(entry.ItemId)
	if err != nil {
		return models.Item{}, err
	}
	_, err = doesMemberExist(entry.MemberId)
	if err != nil {
		return models.Item{}, err
	}
	if diff <= 0 {
		return models.Item{}, models.Error{
			Details: "The amount of Items Lost has to be 1 or greater",
			Path:    "Entry Service - borrowedItemLost()",
			Object:  strconv.Itoa(diff),
			Time:    time.Now(),
		}
	}
	if entry.Capacity < diff {
		return models.Item{}, models.Error{
			Details: "You tried to return more, than you had initially borrowed",
			Path:    "Entry Service - borrowedItemLost()",
			Object:  entry.String(),
			Time:    time.Now(),
		}
	}
	if entry.Capacity == diff {
		err = db.deleteEntryFromDB(entryId)
		if err != nil {
			return models.Item{}, err
		}
	} else {
		_, err = db.updateEntryInDB(entryId, -diff)
		if err != nil {
			return models.Item{}, err
		}
	}
	_, err = changeMemberBorrowCount(entry.MemberId, diff)
	if err != nil {
		return models.Item{}, err
	}
	item, err = changeItemCapacity(entry.ItemId, diff)
	if err != nil {
		return models.Item{}, err
	}
	return item, nil
}

func changeItemCapacity(itemId int, diff int) (models.Item, error) {
	client := &http.Client{}
	req, err := http.NewRequest(http.MethodPut, ItemService+
		strconv.Itoa(itemId)+
		"/lost/"+strconv.Itoa(diff), nil)
	if err != nil {
		return models.Item{}, models.Error{
			Details: err.Error(),
			Path:    "Entry Service - changeItemCapacity()",
			Object:  "",
			Time:    time.Now(),
		}
	}
	resp, err := client.Do(req)
	if err != nil {
		return models.Item{}, models.Error{
			Details: err.Error(),
			Path:    "Entry Service - changeItemCapacity()",
			Object:  "",
			Time:    time.Now(),
		}
	}
	if resp.StatusCode != 200 {
		var apiError models.Error
		err = json.NewDecoder(resp.Body).Decode(&apiError)
		if err != nil {
			return models.Item{}, err
		}
		return models.Item{}, apiError
	}
	defer resp.Body.Close()
	var item models.Item
	err = json.NewDecoder(resp.Body).Decode(&item)
	if err != nil {
		return models.Item{}, models.Error{
			Details: err.Error(),
			Path:    "Entry Service - changeItemCapacity()",
			Object:  "",
			Time:    time.Now(),
		}
	}
	return item, nil
}
