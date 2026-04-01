package handlers

import (
	"net/http"

	"caretop-backend/database"
	"caretop-backend/models"
	"caretop-backend/utils"

	"github.com/gin-gonic/gin"
)

func GetProducts(c *gin.Context) {
	var products []models.Product
	if err := database.DB.Where("is_published = ?", true).Order("created_at DESC").Find(&products).Error; err != nil {
		c.JSON(http.StatusInternalServerError, utils.Error(500, "failed to fetch products"))
		return
	}

	utils.Success(c, products)
}

func GetProduct(c *gin.Context) {
	slug := c.Param("slug")

	var product models.Product
	if err := database.DB.First(&product, "slug = ? AND is_published = ?", slug, true).Error; err != nil {
		c.JSON(http.StatusNotFound, utils.Error(404, "product not found"))
		return
	}

	utils.Success(c, product)
}
