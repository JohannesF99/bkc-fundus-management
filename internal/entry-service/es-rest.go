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
		panic(err)
	}
	c.JSON(http.StatusOK, entries)
}

func singleEntry(c *gin.Context) {
	entryId, err := strconv.Atoi(c.Param("entryId"))
	if err != nil {
		panic(err)
	}
	entry, err := getEntryForEntryId(entryId)
	if err != nil {
		panic(err)
	}
	c.JSON(http.StatusOK, entry)
}

func allEntriesForMemberId(c *gin.Context) {
	memberId, err := strconv.Atoi(c.Param("memberId"))
	if err != nil {
		panic(err)
	}
	entries, err := getAllEntriesByMemberId(memberId)
	if err != nil {
		panic(err)
	}
	c.JSON(http.StatusOK, entries)
}

func allEntriesForItemId(c *gin.Context) {
	itemId, err := strconv.Atoi(c.Param("itemId"))
	if err != nil {
		panic(err)
	}
	entries, err := getAllEntriesByItemId(itemId)
	if err != nil {
		panic(err)
	}
	c.JSON(http.StatusOK, entries)
}

func newEntry(c *gin.Context) {
	var newEntry models.NewEntryInfos
	err := c.BindJSON(&newEntry)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Error{
			Details: "Problem Binding Request Body",
			Path:    "/v1/entry",
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
			Details: "Problem parsing Parameter Entry-ID",
			Path:    "/v1/entry/:entryId",
			Object:  c.Param("entryId"),
			Time:    time.Now(),
		})
		return
	}
	returned, err := strconv.Atoi(c.DefaultQuery("returned", "0"))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Error{
			Details: "Problem parsing Parameter Returned",
			Path:    "/v1/entry/:entryId",
			Object:  c.DefaultQuery("returned", "0"),
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
		panic(err)
	}
	entry, err := deleteEntry(entryId)
	if err != nil {
		panic(err)
	}
	c.JSON(http.StatusOK, entry)
}

func lostItem(c *gin.Context) {
	entryId, err := strconv.Atoi(c.Param("entryId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Error{
			Details: "Problems parsing Parameter Entry-ID",
			Path:    "/v1/entry/:id/lost/:diff",
			Object:  "",
			Time:    time.Now(),
		})
		return
	}
	lost, err := strconv.Atoi(c.Param("lost"))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Error{
			Details: "Problems parsing Parameter Diff",
			Path:    "/v1/entry/:id/lost/:diff",
			Object:  strconv.Itoa(entryId),
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
		panic(err)
	}
	itemId, err := strconv.Atoi(c.Param("itemId"))
	if err != nil {
		panic(err)
	}
	entry, err := getEntryForMemberIdAndItemId(memberId, itemId)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error":    "No matching entry could be found",
			"memberId": memberId,
			"itemId":   itemId,
		})
		return
	}
	c.JSON(http.StatusOK, entry)
}
