package main

import (
	"encoding/json"
	"errors"
	"github.com/JohannesF99/bkc-fundus-management/pkg/models"
	"net/http"
	"strconv"
)

func removeItemForItemId(itemId int) (models.Item, error) {
	client := &http.Client{}
	req, err := http.NewRequest(http.MethodDelete, ItemService+strconv.Itoa(itemId), nil)
	if err != nil {
		return models.Item{}, err
	}
	resp, err := client.Do(req)
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

func changeMemberStatus(memberId int, status bool) (models.Member, error) {
	client := &http.Client{}
	req, err := http.NewRequest(http.MethodPut, MemberService+
		strconv.Itoa(memberId)+
		"/status/"+strconv.FormatBool(status), nil)
	if err != nil {
		return models.Member{}, err
	}
	resp, err := client.Do(req)
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

func changeItemCapacity(itemId int, diff int) (models.Item, error) {
	client := &http.Client{}
	req, err := http.NewRequest(http.MethodPut, ItemService+
		strconv.Itoa(itemId)+
		"/lost/"+strconv.Itoa(diff), nil)
	if err != nil {
		return models.Item{}, err
	}
	resp, err := client.Do(req)
	if err != nil {
		return models.Item{}, nil
	}
	defer resp.Body.Close()
	var item models.Item
	err = json.NewDecoder(resp.Body).Decode(&item)
	if err != nil {
		return models.Item{}, err
	}
	return item, nil
}

func removeEntryWithoutSideeffects(entryId int) (models.Entry, error) {
	client := &http.Client{}
	req, err := http.NewRequest(http.MethodDelete, EntryService+strconv.Itoa(entryId), nil)
	if err != nil {
		return models.Entry{}, err
	}
	resp, err := client.Do(req)
	if err != nil {
		return models.Entry{}, nil
	}
	defer resp.Body.Close()
	var entry models.Entry
	err = json.NewDecoder(resp.Body).Decode(&entry)
	if err != nil {
		return models.Entry{}, err
	}
	return entry, nil
}

func updateEntry(entryId int, capacity int) (models.Entry, error) {
	client := &http.Client{}
	req, err := http.NewRequest(http.MethodPut, EntryService+
		strconv.Itoa(entryId)+
		"?borrowed="+strconv.Itoa(capacity), nil)
	if err != nil {
		return models.Entry{}, err
	}
	resp, err := client.Do(req)
	if err != nil {
		return models.Entry{}, nil
	}
	defer resp.Body.Close()
	var entry models.Entry
	err = json.NewDecoder(resp.Body).Decode(&entry)
	if err != nil {
		return models.Entry{}, err
	}
	_, err = changeItemAvailability(entry.ItemId, capacity)
	if err != nil {
		panic(err)
	}
	_, err = changeMemberBorrowCount(entry.MemberId, capacity)
	if err != nil {
		panic(err)
	}
	return entry, nil
}

func removeExistingEntry(entryId int) (models.Entry, error) {
	entry, err := removeEntryWithoutSideeffects(entryId)
	if err != nil {
		return models.Entry{}, err
	}
	_, err = changeItemAvailability(entry.ItemId, -entry.Capacity)
	if err != nil {
		panic(err)
	}
	_, err = changeMemberBorrowCount(entry.MemberId, -entry.Capacity)
	if err != nil {
		panic(err)
	}
	return entry, nil
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

func doesEntryExistForEntryId(entryId int) (models.Entry, error) {
	resp, err := http.Get(EntryService + strconv.Itoa(entryId))
	if err != nil {
		return models.Entry{}, err
	}
	defer resp.Body.Close()
	var existingEntry models.Entry
	err = json.NewDecoder(resp.Body).Decode(&existingEntry)
	if err != nil {
		return models.Entry{}, err
	}
	if resp.StatusCode != 200 {
		return models.Entry{}, errors.New("no existing Entry found")
	}
	return existingEntry, nil
}
