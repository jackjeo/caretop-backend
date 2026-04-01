package handlers

import (
	"net/http"

	"caretop-backend/database"
	"caretop-backend/models"
	"caretop-backend/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type RegisterInput struct {
	Username string `json:"username" binding:"required,min=3,max=50"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type LoginInput struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

func Register(c *gin.Context) {
	var input RegisterInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, utils.Error(400, "invalid input: "+err.Error()))
		return
	}

	var existingUser models.User
	if err := database.DB.First(&existingUser, "username = ? OR email = ?", input.Username, input.Email).Error; err == nil {
		c.JSON(http.StatusBadRequest, utils.Error(400, "username or email already exists"))
		return
	}

	hashedPassword, err := utils.HashPassword(input.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.Error(500, "failed to hash password"))
		return
	}

	user := models.User{
		Username:     input.Username,
		Email:        input.Email,
		PasswordHash: hashedPassword,
		Role:         models.RoleUser,
		IsActive:     true,
		IsBanned:     false,
	}

	if err := database.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, utils.Error(500, "failed to create user"))
		return
	}

	token, err := utils.GenerateJWT(user.ID, user.Username, string(user.Role))
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.Error(500, "failed to generate token"))
		return
	}

	utils.Success(c, gin.H{
		"token": token,
		"user":  user,
	})
}

func Login(c *gin.Context) {
	var input LoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, utils.Error(400, "invalid input: "+err.Error()))
		return
	}

	var user models.User
	if err := database.DB.First(&user, "email = ?", input.Email).Error; err != nil {
		c.JSON(http.StatusUnauthorized, utils.Error(401, "invalid email or password"))
		return
	}

	if !utils.CheckPassword(input.Password, user.PasswordHash) {
		c.JSON(http.StatusUnauthorized, utils.Error(401, "invalid email or password"))
		return
	}

	if user.IsBanned {
		c.JSON(http.StatusForbidden, utils.Error(403, "user is banned"))
		return
	}

	token, err := utils.GenerateJWT(user.ID, user.Username, string(user.Role))
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.Error(500, "failed to generate token"))
		return
	}

	utils.Success(c, gin.H{
		"token": token,
		"user":  user,
	})
}

func Logout(c *gin.Context) {
	utils.Success(c, gin.H{"message": "logged out successfully"})
}

func GetMe(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, utils.Error(401, "user not found"))
		return
	}

	u := user.(*models.User)
	utils.Success(c, u)
}

func GetUserByID(c *gin.Context) {
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

	utils.Success(c, user)
}
