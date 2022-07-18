package rest

import (
	"github.com/JohannesF99/bkc-fundus-management/internal/member-service/service"
	"github.com/JohannesF99/bkc-fundus-management/pkg/models"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func StartMemberService() {
	r := gin.Default()
	v1 := r.Group("/v1/member")
	{
		v1.GET("/", getAllMember)
		v1.GET("/:id", getMemberWithId)
		v1.POST("/", createNewMember)
		v1.PUT("/:id", updateBorrowCount)
	}
	err := r.Run()
	if err != nil {
		panic(err)
	}
}

func getAllMember(c *gin.Context) {
	members, err := service.GetAllMember()
	if err != nil {
		panic(err)
	}
	c.JSON(http.StatusOK, members)
}

func getMemberWithId(c *gin.Context) {
	userId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		panic(err)
	}
	member, err := service.GetMemberWithId(userId)
	if err != nil {
		panic(err)
	}
	c.JSON(http.StatusOK, member)
}

func createNewMember(c *gin.Context) {
	var apiMember models.ApiMember
	err := c.BindJSON(&apiMember)
	if err != nil {
		panic(err)
	}
	id, err := service.InsertNewMember(*apiMember.ToNormalMember())
	if err != nil {
		panic(err)
	}
	member, err := service.GetMemberWithId(int(id))
	if err != nil {
		panic(err)
	}
	c.JSON(http.StatusOK, member)
}

func updateBorrowCount(c *gin.Context) {
	userId, err := strconv.Atoi(c.Param("id"))
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
	userId, err = service.UpdateBorrowCount(userId, borrowed-returned)
	if err != nil {
		panic(err)
	}
	member, err := service.GetMemberWithId(userId)
	if err != nil {
		panic(err)
	}
	c.JSON(http.StatusOK, member)
}
