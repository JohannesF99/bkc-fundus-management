package main

import (
	"bytes"
	"encoding/json"
	"github.com/JohannesF99/bkc-fundus-management/pkg/constants"
	"github.com/JohannesF99/bkc-fundus-management/pkg/models"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"strconv"
	"time"
)

const EntryService = constants.EntryService
const ItemService = constants.ItemService
const MemberService = constants.MemberService

func main() {
	r := gin.Default()
	accounts := readValidAccountsFromFile()
	authorized := r.Group("/v1/fundus/", gin.BasicAuth(accounts))
	{
		authorized.GET("member", getAllMember)                         //Check
		authorized.GET("item", getAllItems)                            //Check
		authorized.GET("member/:memberId/items", getAllItemsForMember) //Check
		authorized.GET("item/:itemId/members", getAllMembersForItem)   //Check
		authorized.POST("member", registerNewMember)                   //Check
		authorized.POST("item", registerNewItem)                       //Check
		authorized.POST("entry", borrowItem)
		authorized.PUT("entry/:entryId", changeEntryCapacity)
		authorized.PUT("entry/:entryId/lost/:diff", borrowedItemLost)
		authorized.PUT("member/:memberId/status/:status", activateOrDeactivateMember)
		authorized.DELETE("entry/:entryId", removeEntry)
		authorized.DELETE("item/:itemId", removeItem)
	}
	err := r.Run(":8083")
	if err != nil {
		panic(err)
	}
}

func readValidAccountsFromFile() gin.Accounts {
	file, _ := os.Open("scripts/accounts.json")
	defer file.Close()
	decoder := json.NewDecoder(file)
	accounts := gin.Accounts{}
	err := decoder.Decode(&accounts)
	if err != nil {
		os.Exit(-1)
	}
	return accounts
}

func getAllMember(c *gin.Context) {
	resp, err := http.Get(MemberService)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	var memberList []models.Member
	err = json.NewDecoder(resp.Body).Decode(&memberList)
	if err != nil {
		panic(err)
	}
	c.JSON(http.StatusOK, memberList)
}

func getAllItems(c *gin.Context) {
	resp, err := http.Get(ItemService)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	var itemList []models.Item
	err = json.NewDecoder(resp.Body).Decode(&itemList)
	if err != nil {
		panic(err)
	}
	c.JSON(http.StatusOK, itemList)
}

func getAllItemsForMember(c *gin.Context) {
	memberId, err := strconv.Atoi(c.Param("memberId"))
	if err != nil {
		panic(err)
	}
	resp, err := http.Get(EntryService + "member/" + strconv.Itoa(memberId))
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	var itemInfoList []models.Export
	err = json.NewDecoder(resp.Body).Decode(&itemInfoList)
	if err != nil {
		panic(err)
	}
	c.JSON(http.StatusOK, itemInfoList)
}

func getAllMembersForItem(c *gin.Context) {
	itemId, err := strconv.Atoi(c.Param("itemId"))
	if err != nil {
		panic(err)
	}
	resp, err := http.Get(EntryService + "item/" + strconv.Itoa(itemId))
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	var memberInfoList []models.Export
	err = json.NewDecoder(resp.Body).Decode(&memberInfoList)
	if err != nil {
		panic(err)
	}
	c.JSON(http.StatusOK, memberInfoList)
}

func registerNewMember(c *gin.Context) {
	var newAccountInfos models.NewMemberInfos
	err := c.BindJSON(&newAccountInfos)
	if err != nil {
		panic(err)
	}
	postBody, err := json.Marshal(newAccountInfos)
	if err != nil {
		panic(err)
	}
	responseBody := bytes.NewBuffer(postBody)
	resp, err := http.Post(MemberService, "application/json", responseBody)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	var jsonObj models.Member
	err = json.NewDecoder(resp.Body).Decode(&jsonObj)
	if err != nil {
		panic(err)
	}
	c.JSON(http.StatusOK, jsonObj)
}

func registerNewItem(c *gin.Context) {
	var newItemInfos models.NewItemInfos
	err := c.BindJSON(&newItemInfos)
	if err != nil {
		panic(err)
	}
	postBody, err := json.Marshal(newItemInfos)
	if err != nil {
		panic(err)
	}
	responseBody := bytes.NewBuffer(postBody)
	resp, err := http.Post(ItemService, "application/json", responseBody)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	var jsonObj models.Item
	err = json.NewDecoder(resp.Body).Decode(&jsonObj)
	if err != nil {
		panic(err)
	}
	c.JSON(http.StatusOK, jsonObj)
}

func borrowItem(c *gin.Context) {
	//Parse Request Body To Object
	var newEntryInfos models.NewEntryInfos
	err := c.BindJSON(&newEntryInfos)
	if err != nil {
		c.JSON(0, models.Error{
			Details: err.Error(),
			Path:    c.FullPath(),
			Object:  nil,
			Time:    time.Now(),
		})
		return
	}
	//Forward Request To Entry-Service
	postBody, err := json.Marshal(newEntryInfos)
	if err != nil {
		c.JSON(0, models.Error{
			Details: err.Error(),
			Path:    c.FullPath(),
			Object:  nil,
			Time:    time.Now(),
		})
		return
	}
	responseBody := bytes.NewBuffer(postBody)
	resp, err := http.Post(EntryService, "application/json", responseBody)
	if err != nil {
		c.JSON(0, models.Error{
			Details: err.Error(),
			Path:    c.FullPath(),
			Object:  nil,
			Time:    time.Now(),
		})
		return
	}
	defer resp.Body.Close()
	var newEntry models.Entry
	err = json.NewDecoder(resp.Body).Decode(&newEntry)
	if err != nil {
		c.JSON(0, models.Error{
			Details: err.Error(),
			Path:    c.FullPath(),
			Object:  nil,
			Time:    time.Now(),
		})
		return
	}
	c.JSON(http.StatusOK, newEntry)
}

func changeEntryCapacity(c *gin.Context) {
	returned, err := strconv.Atoi(c.Query("returned"))
	if err != nil {
		panic(err)
	}
	entryId, err := strconv.Atoi(c.Param("entryId"))
	if err != nil {
		panic(err)
	}
	existingEntry, err := doesEntryExistForEntryId(entryId)
	if err != nil {
		c.JSON(http.StatusConflict, gin.H{
			"message": "Entry does not exist.",
			"entry":   entryId,
		})
		return
	}
	if existingEntry.Capacity < returned {
		c.JSON(http.StatusConflict, gin.H{
			"message":  "You want to return more than your current capacity.",
			"capacity": existingEntry.Capacity,
			"return":   returned,
		})
		return
	}
	if existingEntry.Capacity == returned {
		entry, err := removeExistingEntry(entryId)
		if err != nil {
			c.JSON(http.StatusConflict, gin.H{
				"message": "Could not Delete Entry",
				"error":   err,
			})
		}
		c.JSON(http.StatusOK, entry)
		return
	}
	entry, err := updateEntry(entryId, -returned)
	if err != nil {
		c.JSON(http.StatusConflict, gin.H{
			"message": "Failed to update Entry, see Log for more Details",
			"error":   err,
		})
	}
	c.JSON(http.StatusOK, entry)
}

func removeEntry(c *gin.Context) {
	entryId, err := strconv.Atoi(c.Param("entryId"))
	if err != nil {
		panic(err)
	}
	_, err = doesEntryExistForEntryId(entryId)
	if err != nil {
		panic(err)
	}
	entry, err := removeExistingEntry(entryId)
	if err != nil {
		panic(err)
	}
	c.JSON(http.StatusOK, entry)
}

func borrowedItemLost(c *gin.Context) {
	entryId, err := strconv.Atoi(c.Param("entryId"))
	if err != nil {
		panic(err)
	}
	diff, err := strconv.Atoi(c.Param("diff"))
	if err != nil {
		panic(err)
	}
	entry, err := doesEntryExistForEntryId(entryId)
	if err != nil {
		panic(err)
	}
	item, err := doesItemExist(entry.ItemId)
	if err != nil {
		panic(err)
	}
	member, err := doesMemberExist(entry.MemberId)
	if err != nil {
		panic(err)
	}
	if entry.Capacity == diff {
		_, err := removeEntryWithoutSideeffects(entryId)
		if err != nil {
			panic(err)
		}
	}
	member, err = changeMemberBorrowCount(member.Id, -diff)
	if err != nil {
		panic(err)
	}
	item, err = changeItemCapacity(item.Id, diff)
	if err != nil {
		panic(err)
	}
	c.JSON(http.StatusOK, item)
}

func activateOrDeactivateMember(c *gin.Context) {
	memberId, err := strconv.Atoi(c.Param("memberId"))
	if err != nil {
		panic(err)
	}
	status, err := strconv.ParseBool(c.Param("status"))
	if err != nil {
		panic(err)
	}
	member, err := doesMemberExist(memberId)
	if err != nil {
		panic(err)
	}
	if member.BorrowedItemCount > 0 {
		c.JSON(http.StatusConflict, gin.H{
			"message": "The Member has to return all borrowed Items before his Status can be changed",
			"member":  member,
		})
		return
	}
	member, err = changeMemberStatus(memberId, status)
	if err != nil {
		panic(err)
	}
	c.JSON(http.StatusOK, member)
}

func removeItem(c *gin.Context) {
	itemId, err := strconv.Atoi(c.Param("itemId"))
	if err != nil {
		panic(err)
	}
	item, err := doesItemExist(itemId)
	if err != nil {
		panic(err)
	}
	if item.Availability != item.Capacity {
		c.JSON(http.StatusConflict, gin.H{
			"message": "All borrowed Items have to be returned before the Item can be deleted!",
			"item":    item,
		})
		return
	}
	item, err = removeItemForItemId(itemId)
	if err != nil {
		panic(err)
	}
	c.JSON(http.StatusOK, item)
}
