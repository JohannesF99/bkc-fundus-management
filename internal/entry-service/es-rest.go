package entry

import (
	"github.com/JohannesF99/bkc-fundus-management/pkg/models"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
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
		panic(err)
	}
	entry, err := createNewEntryOrUpdate(newEntry)
	if err != nil {
		panic(err)
	}
	c.JSON(http.StatusOK, entry)
}

func changeEntry(c *gin.Context) {
	entryId, err := strconv.Atoi(c.Param("entryId"))
	if err != nil {
		panic(err)
	}
	borrowed, err := strconv.Atoi(c.DefaultQuery("borrowed", "0"))
	if err != nil {
		panic(err)
	}
	returned, err := strconv.Atoi(c.DefaultQuery("returned", "0"))
	if err != nil {
		panic(err)
	}
	entry, err := updateEntry(entryId, borrowed-returned)
	if err != nil {
		panic(err)
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
