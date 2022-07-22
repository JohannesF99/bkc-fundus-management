package member

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

func StartMemberService() {
	//Preparation
	if _, err := os.Stat("logs"); errors.Is(err, os.ErrNotExist) {
		err := os.Mkdir("logs", os.ModePerm)
		if err != nil {
			panic(err)
		}
	}
	file, err := os.Create("logs/member-service.log")
	if err != nil {
		panic(err)
	}
	gin.DefaultWriter = io.MultiWriter(file, os.Stdout)
	gin.SetMode(gin.ReleaseMode)
	//Start Server
	r := gin.Default()
	v1 := r.Group("/v1/member")
	{
		v1.GET("/", getAllMember)
		v1.GET("/:id", getMemberWithId)
		v1.POST("/", createNewMember)
		v1.PUT("/:id", updateBorrowCount)
		v1.PUT("/:id/status/:status", changeStatus)
	}
	err = r.Run("localhost:8083")
	if err != nil {
		panic(err)
	}
}

func getAllMember(c *gin.Context) {
	members, err := getAllMembers()
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, members)
}

func getMemberWithId(c *gin.Context) {
	userId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Error{
			Details: err.Error(),
			Path:    "Member Service - getMemberWithId()",
			Object:  "",
			Time:    time.Now(),
		})
		return
	}
	member, err := getMemberWithUserId(userId)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, member)
}

func createNewMember(c *gin.Context) {
	var newAccountInfos models.NewMemberInfos
	err := c.BindJSON(&newAccountInfos)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Error{
			Details: err.Error(),
			Path:    "Member Service - createNewMember()",
			Object:  "",
			Time:    time.Now(),
		})
		return
	}
	id, err := insertNewMember(newAccountInfos)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	member, err := getMemberWithUserId(int(id))
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, member)
}

func updateBorrowCount(c *gin.Context) {
	userId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Error{
			Details: err.Error(),
			Path:    "Member Service - updateBorrowCount()",
			Object:  "",
			Time:    time.Now(),
		})
		return
	}
	borrowed, err := strconv.Atoi(c.DefaultQuery("borrowed", "0"))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Error{
			Details: err.Error(),
			Path:    "Member Service - updateBorrowCount()",
			Object:  "",
			Time:    time.Now(),
		})
		return
	}
	returned, err := strconv.Atoi(c.DefaultQuery("returned", "0"))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Error{
			Details: err.Error(),
			Path:    "Member Service - updateBorrowCount()",
			Object:  "",
			Time:    time.Now(),
		})
		return
	}
	userId, err = updateMemberBorrowCount(userId, borrowed-returned)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	member, err := getMemberWithUserId(userId)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, member)
}

func changeStatus(c *gin.Context) {
	userId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Error{
			Details: err.Error(),
			Path:    "Member Service - changeStatus()",
			Object:  "",
			Time:    time.Now(),
		})
		return
	}
	status, err := strconv.ParseBool(c.Param("status"))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Error{
			Details: err.Error(),
			Path:    "Member Service - changeStatus()",
			Object:  "",
			Time:    time.Now(),
		})
		return
	}
	member, err := changeMemberStatus(userId, status)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, member)
}
