package middleware

import (
	"net/http"
	"strings"

	"caretop-backend/database"
	"caretop-backend/models"
	"caretop-backend/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, utils.Error(401, "authorization header required"))
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, utils.Error(401, "invalid authorization format"))
			c.Abort()
			return
		}

		claims, err := utils.ValidateJWT(parts[1])
		if err != nil {
			c.JSON(http.StatusUnauthorized, utils.Error(401, "invalid or expired token"))
			c.Abort()
			return
		}

		userID, err := uuid.Parse(claims.UserID)
		if err != nil {
			c.JSON(http.StatusUnauthorized, utils.Error(401, "invalid user id"))
			c.Abort()
			return
		}

		var user models.User
		if err := database.DB.First(&user, "id = ?", userID).Error; err != nil {
			c.JSON(http.StatusUnauthorized, utils.Error(401, "user not found"))
			c.Abort()
			return
		}

		if user.IsBanned {
			c.JSON(http.StatusForbidden, utils.Error(403, "user is banned"))
			c.Abort()
			return
		}

		c.Set("userID", userID)
		c.Set("user", &user)
		c.Next()
	}
}

func AdminRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, utils.Error(401, "authorization header required"))
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, utils.Error(401, "invalid authorization format"))
			c.Abort()
			return
		}

		claims, err := utils.ValidateJWT(parts[1])
		if err != nil {
			c.JSON(http.StatusUnauthorized, utils.Error(401, "invalid or expired token"))
			c.Abort()
			return
		}

		userID, err := uuid.Parse(claims.UserID)
		if err != nil {
			c.JSON(http.StatusUnauthorized, utils.Error(401, "invalid user id"))
			c.Abort()
			return
		}

		var user models.User
		if err := database.DB.First(&user, "id = ?", userID).Error; err != nil {
			c.JSON(http.StatusUnauthorized, utils.Error(401, "user not found"))
			c.Abort()
			return
		}

		if user.Role != models.RoleAdmin {
			c.JSON(http.StatusForbidden, utils.Error(403, "admin access required"))
			c.Abort()
			return
		}

		c.Set("userID", userID)
		c.Set("user", &user)
		c.Next()
	}
}
