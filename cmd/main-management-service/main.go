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
		authorized.GET("member", getAllMember)
		authorized.GET("item", getAllItems)
		authorized.GET("member/:memberId/items", getAllItemsForMember)
		authorized.GET("item/:itemId/members", getAllMembersForItem)
		authorized.POST("member", registerNewMember)
		authorized.POST("item", registerNewItem)
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
	file, err := os.Open("scripts/accounts.json")
	if err != nil {
		os.Exit(-1)
	}
	defer file.Close()
	decoder := json.NewDecoder(file)
	accounts := gin.Accounts{}
	err = decoder.Decode(&accounts)
	if err != nil {
		os.Exit(-1)
	}
	return accounts
}

func getAllMember(c *gin.Context) {
	resp, err := http.Get(MemberService)
	if err != nil {
		c.JSON(0, models.Error{
			Details: err.Error(),
			Path:    c.FullPath(),
			Object:  "",
			Time:    time.Now(),
		})
		return
	}
	defer resp.Body.Close()
	var memberList []models.Member
	err = json.NewDecoder(resp.Body).Decode(&memberList)
	if err != nil {
		c.JSON(0, models.Error{
			Details: err.Error(),
			Path:    c.FullPath(),
			Object:  "",
			Time:    time.Now(),
		})
		return
	}
	c.JSON(http.StatusOK, memberList)
}

func getAllItems(c *gin.Context) {
	resp, err := http.Get(ItemService)
	if err != nil {
		c.JSON(0, models.Error{
			Details: err.Error(),
			Path:    c.FullPath(),
			Object:  "",
			Time:    time.Now(),
		})
		return
	}
	defer resp.Body.Close()
	var itemList []models.Item
	err = json.NewDecoder(resp.Body).Decode(&itemList)
	if err != nil {
		c.JSON(0, models.Error{
			Details: err.Error(),
			Path:    c.FullPath(),
			Object:  "",
			Time:    time.Now(),
		})
		return
	}
	c.JSON(http.StatusOK, itemList)
}

func getAllItemsForMember(c *gin.Context) {
	memberId, err := strconv.Atoi(c.Param("memberId"))
	if err != nil {
		c.JSON(0, models.Error{
			Details: err.Error(),
			Path:    c.FullPath(),
			Object:  "",
			Time:    time.Now(),
		})
		return
	}
	resp, err := http.Get(EntryService + "member/" + strconv.Itoa(memberId))
	if err != nil {
		c.JSON(0, models.Error{
			Details: err.Error(),
			Path:    c.FullPath(),
			Object:  "",
			Time:    time.Now(),
		})
		return
	}
	defer resp.Body.Close()
	var itemInfoList []models.Export
	err = json.NewDecoder(resp.Body).Decode(&itemInfoList)
	if err != nil {
		c.JSON(0, models.Error{
			Details: err.Error(),
			Path:    c.FullPath(),
			Object:  "",
			Time:    time.Now(),
		})
		return
	}
	c.JSON(http.StatusOK, itemInfoList)
}

func getAllMembersForItem(c *gin.Context) {
	itemId, err := strconv.Atoi(c.Param("itemId"))
	if err != nil {
		c.JSON(0, models.Error{
			Details: err.Error(),
			Path:    c.FullPath(),
			Object:  "",
			Time:    time.Now(),
		})
		return
	}
	resp, err := http.Get(EntryService + "item/" + strconv.Itoa(itemId))
	if err != nil {
		c.JSON(0, models.Error{
			Details: err.Error(),
			Path:    c.FullPath(),
			Object:  "",
			Time:    time.Now(),
		})
		return
	}
	defer resp.Body.Close()
	var memberInfoList []models.Export
	err = json.NewDecoder(resp.Body).Decode(&memberInfoList)
	if err != nil {
		c.JSON(0, models.Error{
			Details: err.Error(),
			Path:    c.FullPath(),
			Object:  "",
			Time:    time.Now(),
		})
		return
	}
	c.JSON(http.StatusOK, memberInfoList)
}

func registerNewMember(c *gin.Context) {
	var newAccountInfos models.NewMemberInfos
	err := c.BindJSON(&newAccountInfos)
	if err != nil {
		c.JSON(0, models.Error{
			Details: err.Error(),
			Path:    c.FullPath(),
			Object:  "",
			Time:    time.Now(),
		})
		return
	}
	postBody, err := json.Marshal(newAccountInfos)
	if err != nil {
		c.JSON(0, models.Error{
			Details: err.Error(),
			Path:    c.FullPath(),
			Object:  "",
			Time:    time.Now(),
		})
		return
	}
	responseBody := bytes.NewBuffer(postBody)
	resp, err := http.Post(MemberService, "application/json", responseBody)
	if err != nil {
		c.JSON(0, models.Error{
			Details: err.Error(),
			Path:    c.FullPath(),
			Object:  "",
			Time:    time.Now(),
		})
		return
	}
	defer resp.Body.Close()
	var jsonObj models.Member
	err = json.NewDecoder(resp.Body).Decode(&jsonObj)
	if err != nil {
		c.JSON(0, models.Error{
			Details: err.Error(),
			Path:    c.FullPath(),
			Object:  "",
			Time:    time.Now(),
		})
		return
	}
	c.JSON(http.StatusOK, jsonObj)
}

func registerNewItem(c *gin.Context) {
	var newItemInfos models.NewItemInfos
	err := c.BindJSON(&newItemInfos)
	if err != nil {
		c.JSON(0, models.Error{
			Details: err.Error(),
			Path:    c.FullPath(),
			Object:  "",
			Time:    time.Now(),
		})
		return
	}
	postBody, err := json.Marshal(newItemInfos)
	if err != nil {
		c.JSON(0, models.Error{
			Details: err.Error(),
			Path:    c.FullPath(),
			Object:  "",
			Time:    time.Now(),
		})
		return
	}
	responseBody := bytes.NewBuffer(postBody)
	resp, err := http.Post(ItemService, "application/json", responseBody)
	if err != nil {
		c.JSON(0, models.Error{
			Details: err.Error(),
			Path:    c.FullPath(),
			Object:  "",
			Time:    time.Now(),
		})
		return
	}
	defer resp.Body.Close()
	var jsonObj models.Item
	err = json.NewDecoder(resp.Body).Decode(&jsonObj)
	if err != nil {
		c.JSON(0, models.Error{
			Details: err.Error(),
			Path:    c.FullPath(),
			Object:  "",
			Time:    time.Now(),
		})
		return
	}
	c.JSON(http.StatusOK, jsonObj)
}

func borrowItem(c *gin.Context) {
	var newEntryInfos models.NewEntryInfos
	err := c.BindJSON(&newEntryInfos)
	if err != nil {
		c.JSON(0, models.Error{
			Details: err.Error(),
			Path:    c.FullPath(),
			Object:  "",
			Time:    time.Now(),
		})
		return
	}
	postBody, err := json.Marshal(newEntryInfos)
	if err != nil {
		c.JSON(0, models.Error{
			Details: err.Error(),
			Path:    c.FullPath(),
			Object:  "",
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
			Object:  "",
			Time:    time.Now(),
		})
		return
	}
	if resp.StatusCode != 200 {
		var apiError models.Error
		err = json.NewDecoder(resp.Body).Decode(&apiError)
		if err != nil {
			c.JSON(0, models.Error{
				Details: err.Error(),
				Path:    c.FullPath(),
				Object:  "",
				Time:    time.Now(),
			})
			return
		}
		c.JSON(resp.StatusCode, err)
		return
	}
	defer resp.Body.Close()
	var newEntry models.Entry
	err = json.NewDecoder(resp.Body).Decode(&newEntry)
	if err != nil {
		c.JSON(0, models.Error{
			Details: err.Error(),
			Path:    c.FullPath(),
			Object:  "",
			Time:    time.Now(),
		})
		return
	}
	c.JSON(http.StatusOK, newEntry)
}

func changeEntryCapacity(c *gin.Context) {
	entryId, err := strconv.Atoi(c.Param("entryId"))
	if err != nil {
		c.JSON(0, models.Error{
			Details: err.Error(),
			Path:    c.FullPath(),
			Object:  "",
			Time:    time.Now(),
		})
		return
	}
	returned, err := strconv.Atoi(c.Query("returned"))
	if err != nil {
		c.JSON(0, models.Error{
			Details: err.Error(),
			Path:    c.FullPath(),
			Object:  strconv.Itoa(entryId),
			Time:    time.Now(),
		})
		return
	}
	client := &http.Client{}
	req, err := http.NewRequest(http.MethodPut, EntryService+
		strconv.Itoa(entryId)+
		"?returned="+strconv.Itoa(returned), nil)
	if err != nil {
		c.JSON(0, models.Error{
			Details: err.Error(),
			Path:    c.FullPath(),
			Object:  strconv.Itoa(entryId),
			Time:    time.Now(),
		})
		return
	}
	resp, err := client.Do(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Error{
			Details: err.Error(),
			Path:    "Client-Request",
			Object:  "",
			Time:    time.Now(),
		})
		return
	}
	if resp.StatusCode != 200 {
		var apiError models.Error
		err = json.NewDecoder(resp.Body).Decode(&apiError)
		c.JSON(resp.StatusCode, apiError)
		return
	}
	defer resp.Body.Close()
	var entry models.Entry
	err = json.NewDecoder(resp.Body).Decode(&entry)
	if err != nil {
		c.JSON(0, models.Error{
			Details: err.Error(),
			Path:    c.FullPath(),
			Object:  strconv.Itoa(entryId),
			Time:    time.Now(),
		})
		return
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
	client := &http.Client{}
	req, err := http.NewRequest(http.MethodPut, EntryService+
		strconv.Itoa(entryId)+
		"/lost/"+strconv.Itoa(diff), nil)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Error{
			Details: err.Error(),
			Path:    c.FullPath(),
			Object:  strconv.Itoa(entryId),
			Time:    time.Now(),
		})
		return
	}
	resp, err := client.Do(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	defer resp.Body.Close()
	var item models.Item
	err = json.NewDecoder(resp.Body).Decode(&item)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Error{
			Details: err.Error(),
			Path:    c.FullPath(),
			Object:  strconv.Itoa(entryId),
			Time:    time.Now(),
		})
		return
	}
	c.JSON(http.StatusOK, item)
}

func activateOrDeactivateMember(c *gin.Context) {
	memberId, err := strconv.Atoi(c.Param("memberId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Error{
			Details: err.Error(),
			Path:    c.FullPath(),
			Object:  "",
			Time:    time.Now(),
		})
		return
	}
	status, err := strconv.ParseBool(c.Param("status"))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Error{
			Details: err.Error(),
			Path:    c.FullPath(),
			Object:  "",
			Time:    time.Now(),
		})
		return
	}
	client := &http.Client{}
	req, err := http.NewRequest(http.MethodPut, MemberService+
		strconv.Itoa(memberId)+
		"/status/"+strconv.FormatBool(status), nil)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Error{
			Details: err.Error(),
			Path:    c.FullPath(),
			Object:  "",
			Time:    time.Now(),
		})
		return
	}
	resp, err := client.Do(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Error{
			Details: err.Error(),
			Path:    c.FullPath(),
			Object:  "",
			Time:    time.Now(),
		})
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		var apiError models.Error
		err = json.NewDecoder(resp.Body).Decode(&apiError)
		if err != nil {
			c.JSON(http.StatusBadRequest, models.Error{
				Details: err.Error(),
				Path:    c.FullPath(),
				Object:  "",
				Time:    time.Now(),
			})
			return
		}
		c.JSON(http.StatusBadRequest, apiError)
		return
	}
	var member models.Member
	err = json.NewDecoder(resp.Body).Decode(&member)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Error{
			Details: err.Error(),
			Path:    c.FullPath(),
			Object:  "",
			Time:    time.Now(),
		})
		return
	}
	c.JSON(http.StatusOK, member)
}

func removeEntry(c *gin.Context) {
	entryId, err := strconv.Atoi(c.Param("entryId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Error{
			Details: err.Error(),
			Path:    c.FullPath(),
			Object:  "",
			Time:    time.Now(),
		})
		return
	}
	client := &http.Client{}
	req, err := http.NewRequest(http.MethodDelete, EntryService+strconv.Itoa(entryId), nil)
	if err != nil {
		panic(err)
	}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	var entry models.Entry
	err = json.NewDecoder(resp.Body).Decode(&entry)
	if err != nil {
		panic(err)
	}
	c.JSON(http.StatusOK, entry)
}

func removeItem(c *gin.Context) {
	itemId, err := strconv.Atoi(c.Param("itemId"))
	client := &http.Client{}
	req, err := http.NewRequest(http.MethodDelete, ItemService+strconv.Itoa(itemId), nil)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Error{
			Details: err.Error(),
			Path:    c.FullPath(),
			Object:  "",
			Time:    time.Now(),
		})
		return
	}
	resp, err := client.Do(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Error{
			Details: err.Error(),
			Path:    c.FullPath(),
			Object:  "",
			Time:    time.Now(),
		})
		return
	}
	defer resp.Body.Close()
	var item models.Item
	err = json.NewDecoder(resp.Body).Decode(&item)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Error{
			Details: err.Error(),
			Path:    c.FullPath(),
			Object:  "",
			Time:    time.Now(),
		})
		return
	}
	c.JSON(http.StatusOK, item)
}
