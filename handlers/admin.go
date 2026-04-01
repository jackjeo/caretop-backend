package handlers

import (
	"net/http"

	"caretop-backend/database"
	"caretop-backend/models"
	"caretop-backend/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func GetStats(c *gin.Context) {
	var userCount int64
	var postCount int64
	var threadCount int64
	var ticketCount int64

	database.DB.Model(&models.User{}).Count(&userCount)
	database.DB.Model(&models.BlogPost{}).Count(&postCount)
	database.DB.Model(&models.ForumThread{}).Count(&threadCount)
	database.DB.Model(&models.Ticket{}).Count(&ticketCount)

	utils.Success(c, gin.H{
		"users":   userCount,
		"posts":   postCount,
		"threads": threadCount,
		"tickets": ticketCount,
	})
}

func GetUsers(c *gin.Context) {
	var users []models.User
	database.DB.Order("created_at DESC").Find(&users)
	utils.Success(c, users)
}

type UpdateRoleInput struct {
	Role string `json:"role" binding:"required"`
}

func UpdateUserRole(c *gin.Context) {
	id := c.Param("id")
	userID, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.Error(400, "invalid user id"))
		return
	}

	var input UpdateRoleInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, utils.Error(400, "invalid input: "+err.Error()))
		return
	}

	var user models.User
	if err := database.DB.First(&user, "id = ?", userID).Error; err != nil {
		c.JSON(http.StatusNotFound, utils.Error(404, "user not found"))
		return
	}

	user.Role = models.Role(input.Role)
	database.DB.Save(&user)

	utils.Success(c, user)
}

func BanUser(c *gin.Context) {
	id := c.Param("id")
	userID, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.Error(400, "invalid user id"))
		return
	}

	var user models.User
	if err := database.DB.First(&user, "id = ?", userID).Error; err != nil {
		c.JSON(http.StatusNotFound, utils.Error(404, "user not found"))
		return
	}

	user.IsBanned = !user.IsBanned
	database.DB.Save(&user)

	utils.Success(c, gin.H{
		"is_banned": user.IsBanned,
		"message":   map[bool]string{true: "user banned", false: "user unbanned"}[user.IsBanned],
	})
}

func GetForumStats(c *gin.Context) {
	var boardCount int64
	var threadCount int64
	var postCount int64
	var userCount int64

	database.DB.Model(&models.ForumBoard{}).Count(&boardCount)
	database.DB.Model(&models.ForumThread{}).Count(&threadCount)
	database.DB.Model(&models.ForumPost{}).Count(&postCount)
	database.DB.Model(&models.User{}).Count(&userCount)

	utils.Success(c, gin.H{
		"boards":  boardCount,
		"threads": threadCount,
		"posts":   postCount,
		"users":   userCount,
	})
}
