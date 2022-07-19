package items

import (
	"github.com/JohannesF99/bkc-fundus-management/pkg/models"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func StartItemService() {
	r := gin.Default()
	v1 := r.Group("/v1/item")
	{
		v1.GET("/", allItems)
		v1.GET("/:id", fetchItem)
		v1.POST("/", addItem)
		v1.PUT("/:id", updateItem)
		v1.DELETE("/:id", removeItem)
	}
	err := r.Run("localhost:8081")
	if err != nil {
		panic(err)
	}
}

func allItems(c *gin.Context) {
	items, err := getAllItems()
	if err != nil {
		panic(err)
	}
	c.JSON(http.StatusOK, items)
}

func fetchItem(c *gin.Context) {
	itemId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		panic(err)
	}
	item, err := getItemWithId(itemId)
	if err != nil {
		panic(err)
	}
	c.JSON(http.StatusOK, item)
}

func addItem(c *gin.Context) {
	var newItem models.NewItemInfos
	err := c.BindJSON(&newItem)
	if err != nil {
		panic(err)
	}
	itemId, err := insertNewItem(newItem)
	if err != nil {
		panic(err)
	}
	item, err := getItemWithId(int(itemId))
	if err != nil {
		panic(err)
	}
	c.JSON(http.StatusOK, item)
}

func updateItem(c *gin.Context) {
	itemId, err := strconv.Atoi(c.Param("id"))
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
	id, err := updateItemAvailability(itemId, returned-borrowed)
	if err != nil {
		panic(err)
	}
	item, err := getItemWithId(id)
	if err != nil {
		panic(err)
	}
	c.JSON(http.StatusOK, item)
}

func removeItem(c *gin.Context) {
	itemId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		panic(err)
	}
	item, err := getItemWithId(itemId)
	if err != nil {
		panic(err)
	}
	err = deleteItem(itemId)
	if err != nil {
		panic(err)
	}
	c.JSON(http.StatusOK, item)
}
