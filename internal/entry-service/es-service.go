package entry

import (
	"encoding/json"
	"errors"
	"github.com/JohannesF99/bkc-fundus-management/pkg/constants"
	"github.com/JohannesF99/bkc-fundus-management/pkg/models"
	"net/http"
	"strconv"
)

const ItemService = constants.ItemService
const MemberService = constants.MemberService

func createNewEntryOrUpdate(newEntryInfo models.NewEntryInfos) (models.Entry, error) {
	_, err := doesMemberExist(newEntryInfo.MemberId)
	if err != nil {
		return models.Entry{}, err
	}
	item, err := doesItemExist(newEntryInfo.ItemId)
	if err != nil {
		return models.Entry{}, err
	}
	if item.Availability < newEntryInfo.Capacity {
		return models.Entry{}, err
	}
	db, err := connect()
	if err != nil {
		return models.Entry{}, err
	}
	if existingEntry, err := db.getEntryForMemberIdAndItemIdFromDB(newEntryInfo.MemberId, newEntryInfo.ItemId); err == nil {
		updatedEntry, err := updateEntry(existingEntry.Id, newEntryInfo.Capacity)
		if err != nil {
			return models.Entry{}, err
		}
		_, err = changeItemAvailability(updatedEntry.ItemId, newEntryInfo.Capacity)
		if err != nil {
			return models.Entry{}, err
		}
		_, err = changeMemberBorrowCount(updatedEntry.MemberId, newEntryInfo.Capacity)
		if err != nil {
			return models.Entry{}, err
		}
		return updatedEntry, nil
	} else {
		newEntry, err := createNewEntry(newEntryInfo)
		if err != nil {
			return models.Entry{}, err
		}
		_, err = changeItemAvailability(newEntry.ItemId, newEntryInfo.Capacity)
		if err != nil {
			return models.Entry{}, err
		}
		_, err = changeMemberBorrowCount(newEntry.MemberId, newEntryInfo.Capacity)
		if err != nil {
			return models.Entry{}, err
		}
		return newEntry, nil
	}
}

func changeItemAvailability(itemId int, borrowed int) (models.Item, error) {
	req, err := http.NewRequest(http.MethodPut,
		ItemService+strconv.Itoa(itemId)+
			"?borrowed="+strconv.Itoa(borrowed),
		nil)
	if err != nil {
		return models.Item{}, err
	}
	client := &http.Client{}
	putResp, err := client.Do(req)
	if err != nil {
		return models.Item{}, err
	}
	defer putResp.Body.Close()
	var updatedItem models.Item
	err = json.NewDecoder(putResp.Body).Decode(&updatedItem)
	if err != nil {
		return models.Item{}, err
	}
	return updatedItem, nil
}

func changeMemberBorrowCount(memberId int, borrowed int) (models.Member, error) {
	req, err := http.NewRequest(http.MethodPut,
		MemberService+strconv.Itoa(memberId)+
			"?borrowed="+strconv.Itoa(borrowed),
		nil)
	if err != nil {
		return models.Member{}, err
	}
	client := &http.Client{}
	putResp, err := client.Do(req)
	if err != nil {
		return models.Member{}, err
	}
	defer putResp.Body.Close()
	var updatedMember models.Member
	err = json.NewDecoder(putResp.Body).Decode(&updatedMember)
	if err != nil {
		return models.Member{}, err
	}
	return updatedMember, nil
}

func doesMemberExist(memberId int) (models.Member, error) {
	resp, err := http.Get(MemberService + strconv.Itoa(memberId))
	if err != nil {
		return models.Member{}, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return models.Member{}, errors.New("member does not exist in Database. memberId: " + strconv.Itoa(memberId))
	}
	var existingMember models.Member
	err = json.NewDecoder(resp.Body).Decode(&existingMember)
	if err != nil {
		return models.Member{}, err
	}
	return existingMember, nil
}

func doesItemExist(itemId int) (models.Item, error) {
	resp, err := http.Get(ItemService + strconv.Itoa(itemId))
	if err != nil {
		return models.Item{}, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return models.Item{}, errors.New("item does not exist in Database. itemId: " + strconv.Itoa(itemId))
	}
	var itemInfo models.Item
	err = json.NewDecoder(resp.Body).Decode(&itemInfo)
	if err != nil {
		return models.Item{}, err
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
	entry, err := db.getEntriesForEntryIdFromDB(entryId)
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
