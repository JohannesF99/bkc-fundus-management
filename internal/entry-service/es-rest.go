package entry

import (
	"github.com/JohannesF99/bkc-fundus-management/pkg/models"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
)

func StartEntryService() {
	r := gin.Default()
	v1 := r.Group("/v1/entry")
	{
		v1.GET("", allEntries)
		v1.GET("/:entryId", singleEntry)
		v1.GET("/member/:memberId", allEntriesForMemberId)
		v1.GET("/member/all", allEntriesForAllMember)
		v1.GET("/item/:itemId", allEntriesForItemId)
		v1.GET("/member/:memberId/item/:itemId", entryForMemberIdAndItemId)
		v1.POST("", newEntry)
		v1.PUT("/:entryId", changeEntry)
		v1.PUT("/:entryId/lost/:lost", lostItem)
		v1.DELETE("/:entryId", removeEntry)
	}
	err := r.Run("localhost:8080")
	if err != nil {
		panic(err)
	}
}

func allEntries(c *gin.Context) {
	entries, err := getAllEntries()
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, entries)
}

func singleEntry(c *gin.Context) {
	entryId, err := strconv.Atoi(c.Param("entryId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Error{
			Details: err.Error(),
			Path:    "Entry Service - singleEntry()",
			Object:  "",
			Time:    time.Now(),
		})
		return
	}
	entry, err := getEntryForEntryId(entryId)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, entry)
}

func allEntriesForMemberId(c *gin.Context) {
	memberId, err := strconv.Atoi(c.Param("memberId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Error{
			Details: err.Error(),
			Path:    "Entry Service - allEntriesForMemberId()",
			Object:  "",
			Time:    time.Now(),
		})
		return
	}
	entries, err := getAllEntriesByMemberId(memberId)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, entries)
}

func allEntriesForAllMember(c *gin.Context) {
	entriesForAllMember, err := getEntriesForAllMember()
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, entriesForAllMember)
}

func allEntriesForItemId(c *gin.Context) {
	itemId, err := strconv.Atoi(c.Param("itemId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Error{
			Details: err.Error(),
			Path:    "Entry Service - allEntriesForItemId()",
			Object:  "",
			Time:    time.Now(),
		})
		return
	}
	entries, err := getAllEntriesByItemId(itemId)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, entries)
}

func newEntry(c *gin.Context) {
	var newEntry models.NewEntryInfos
	err := c.BindJSON(&newEntry)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Error{
			Details: err.Error(),
			Path:    "Entry Service - newEntry()",
			Object:  "",
			Time:    time.Now(),
		})
		return
	}
	entry, err := createNewEntryOrUpdate(newEntry)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, entry)
}

func changeEntry(c *gin.Context) {
	entryId, err := strconv.Atoi(c.Param("entryId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Error{
			Details: err.Error(),
			Path:    "Entry Service - changeEntry()",
			Object:  "",
			Time:    time.Now(),
		})
		return
	}
	returned, err := strconv.Atoi(c.DefaultQuery("returned", "0"))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Error{
			Details: err.Error(),
			Path:    "Entry Service - changeEntry()",
			Object:  "",
			Time:    time.Now(),
		})
		return
	}
	entry, err := updateEntry(entryId, returned)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, entry)
}

func removeEntry(c *gin.Context) {
	entryId, err := strconv.Atoi(c.Param("entryId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Error{
			Details: err.Error(),
			Path:    "Entry Service - removeEntry()",
			Object:  "",
			Time:    time.Now(),
		})
		return
	}
	entry, err := deleteEntry(entryId)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, entry)
}

func lostItem(c *gin.Context) {
	entryId, err := strconv.Atoi(c.Param("entryId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Error{
			Details: err.Error(),
			Path:    "Entry Service - lostItem()",
			Object:  "",
			Time:    time.Now(),
		})
		return
	}
	lost, err := strconv.Atoi(c.Param("lost"))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Error{
			Details: err.Error(),
			Path:    "Entry Service - lostItem()",
			Object:  "",
			Time:    time.Now(),
		})
		return
	}
	item, err := borrowedItemLost(entryId, lost)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, item)
}

func entryForMemberIdAndItemId(c *gin.Context) {
	memberId, err := strconv.Atoi(c.Param("memberId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Error{
			Details: err.Error(),
			Path:    "Entry Service - entryForMemberIdAndItemId()",
			Object:  "",
			Time:    time.Now(),
		})
		return
	}
	itemId, err := strconv.Atoi(c.Param("itemId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Error{
			Details: err.Error(),
			Path:    "Entry Service - entryForMemberIdAndItemId()",
			Object:  "",
			Time:    time.Now(),
		})
		return
	}
	entry, err := getEntryForMemberIdAndItemId(memberId, itemId)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, entry)
}
