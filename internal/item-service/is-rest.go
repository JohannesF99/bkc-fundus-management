package items

import (
	"errors"
	"github.com/JohannesF99/bkc-fundus-management/pkg/models"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"
)

func StartItemService() {
	//Preparation
	if _, err := os.Stat("logs"); errors.Is(err, os.ErrNotExist) {
		err := os.Mkdir("logs", os.ModePerm)
		if err != nil {
			panic(err)
		}
	}
	file, err := os.Create("logs/item-service.log")
	if err != nil {
		panic(err)
	}
	gin.DefaultWriter = io.MultiWriter(file, os.Stdout)
	gin.SetMode(gin.ReleaseMode)
	//Start Server
	r := gin.Default()
	v1 := r.Group("/v1/item")
	{
		v1.GET("/", allItems)
		v1.GET("/:id", fetchItem)
		v1.POST("/", addItem)
		v1.PUT("/:id", updateItem)
		v1.PUT("/:id/lost/:lost", lostItem)
		v1.DELETE("/:id", removeItem)
	}
	_, err = gin.DefaultWriter.Write([]byte("Starting Item Service on localhost:8082\n"))
	if err != nil {
		panic(err)
	}
	err = r.Run("localhost:8082")
	if err != nil {
		panic(err)
	}
}

func allItems(c *gin.Context) {
	items, err := getAllItems()
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, items)
}

func fetchItem(c *gin.Context) {
	itemId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Error{
			Details: err.Error(),
			Path:    "Item Service - fetchItem()",
			Object:  c.Param("id"),
			Time:    time.Now(),
		})
		return
	}
	item, err := getItemWithId(itemId)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, item)
}

func addItem(c *gin.Context) {
	var newItem models.NewItemInfos
	err := c.BindJSON(&newItem)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Error{
			Details: err.Error(),
			Path:    "Item Service - addItem()",
			Object:  "",
			Time:    time.Now(),
		})
		return
	}
	itemId, err := insertNewItem(newItem)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	item, err := getItemWithId(int(itemId))
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, item)
}

func updateItem(c *gin.Context) {
	itemId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Error{
			Details: err.Error(),
			Path:    "Item Service - updateItem()",
			Object:  c.Param("id"),
			Time:    time.Now(),
		})
		return
	}
	borrowed, err := strconv.Atoi(c.DefaultQuery("borrowed", "0"))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Error{
			Details: err.Error(),
			Path:    "Item Service - updateItem()",
			Object:  c.DefaultQuery("borrowed", "0"),
			Time:    time.Now(),
		})
		return
	}
	returned, err := strconv.Atoi(c.DefaultQuery("returned", "0"))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Error{
			Details: err.Error(),
			Path:    "Item Service - updateItem()",
			Object:  c.DefaultQuery("returned", "0"),
			Time:    time.Now(),
		})
		return
	}
	id, err := updateItemAvailability(itemId, returned-borrowed)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	item, err := getItemWithId(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, item)
}

func lostItem(c *gin.Context) {
	itemId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Error{
			Details: err.Error(),
			Path:    "Item Service - removeItem()",
			Object:  c.Param("id"),
			Time:    time.Now(),
		})
		return
	}
	lost, err := strconv.Atoi(c.Param("lost"))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Error{
			Details: err.Error(),
			Path:    "Item Service - removeItem()",
			Object:  c.Param("id"),
			Time:    time.Now(),
		})
		return
	}
	item, err := updateItemCapacity(itemId, lost)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, item)
}

func removeItem(c *gin.Context) {
	itemId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Error{
			Details: err.Error(),
			Path:    "Item Service - removeItem()",
			Object:  c.Param("id"),
			Time:    time.Now(),
		})
		return
	}
	item, err := deleteItem(itemId)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, item)
}
