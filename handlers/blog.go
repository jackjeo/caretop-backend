package handlers

import (
	"net/http"
	"strconv"
	"time"

	"caretop-backend/database"
	"caretop-backend/models"
	"caretop-backend/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type BlogListResponse struct {
	Posts      []models.BlogPost `json:"posts"`
	Total      int64             `json:"total"`
	Page       int               `json:"page"`
	PageSize   int               `json:"page_size"`
	TotalPages int               `json:"total_pages"`
}

func GetBlogPosts(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize := 10
	category := c.Query("category")

	if page < 1 {
		page = 1
	}

	var posts []models.BlogPost
	var total int64

	query := database.DB.Model(&models.BlogPost{}).Where("is_published = ?", true)

	if category != "" {
		query = query.Where("category = ?", category)
	}

	query.Count(&total)
	query.Order("created_at DESC").Offset((page - 1) * pageSize).Limit(pageSize).Preload("Author").Find(&posts)

	totalPages := int(total) / pageSize
	if int(total)%pageSize > 0 {
		totalPages++
	}

	utils.Success(c, BlogListResponse{
		Posts:      posts,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	})
}

func GetBlogPost(c *gin.Context) {
	slug := c.Param("slug")

	var post models.BlogPost
	if err := database.DB.Preload("Author").First(&post, "slug = ? AND is_published = ?", slug, true).Error; err != nil {
		c.JSON(http.StatusNotFound, utils.Error(404, "post not found"))
		return
	}

	database.DB.Model(&post).UpdateColumn("view_count", post.ViewCount+1)
	post.ViewCount++

	utils.Success(c, post)
}

func GetBlogCategories(c *gin.Context) {
	categories := []gin.H{
		{"key": "tech", "label": "技术"},
		{"key": "product", "label": "产品"},
		{"key": "industry", "label": "行业"},
	}
	utils.Success(c, categories)
}

func LikeBlogPost(c *gin.Context) {
	slug := c.Param("slug")

	var post models.BlogPost
	if err := database.DB.First(&post, "slug = ?", slug).Error; err != nil {
		c.JSON(http.StatusNotFound, utils.Error(404, "post not found"))
		return
	}

	database.DB.Model(&post).UpdateColumn("like_count", post.LikeCount+1)

	utils.Success(c, gin.H{"like_count": post.LikeCount + 1})
}

type CreateBlogInput struct {
	Title       string `json:"title" binding:"required"`
	Slug        string `json:"slug" binding:"required"`
	Summary     string `json:"summary"`
	Content     string `json:"content" binding:"required"`
	CoverImageURL string `json:"cover_image_url"`
	Category    string `json:"category" binding:"required"`
	IsPublished bool   `json:"is_published"`
}

func CreateBlogPost(c *gin.Context) {
	var input CreateBlogInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, utils.Error(400, "invalid input: "+err.Error()))
		return
	}

	user, _ := c.Get("user")
	u := user.(*models.User)

	post := models.BlogPost{
		Title:         input.Title,
		Slug:          input.Slug,
		Summary:       input.Summary,
		Content:       input.Content,
		CoverImageURL: input.CoverImageURL,
		Category:      models.BlogCategory(input.Category),
		AuthorID:      u.ID,
		IsPublished:   input.IsPublished,
	}

	if input.IsPublished {
		now := time.Now()
		post.PublishedAt = &now
	}

	if err := database.DB.Create(&post).Error; err != nil {
		c.JSON(http.StatusInternalServerError, utils.Error(500, "failed to create post"))
		return
	}

	database.DB.Preload("Author").First(&post, "id = ?", post.ID)
	utils.Success(c, post)
}

func UpdateBlogPost(c *gin.Context) {
	id := c.Param("id")
	postID, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.Error(400, "invalid post id"))
		return
	}

	var post models.BlogPost
	if err := database.DB.First(&post, "id = ?", postID).Error; err != nil {
		c.JSON(http.StatusNotFound, utils.Error(404, "post not found"))
		return
	}

	var input CreateBlogInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, utils.Error(400, "invalid input: "+err.Error()))
		return
	}

	post.Title = input.Title
	post.Slug = input.Slug
	post.Summary = input.Summary
	post.Content = input.Content
	post.CoverImageURL = input.CoverImageURL
	post.Category = models.BlogCategory(input.Category)
	post.IsPublished = input.IsPublished

	if input.IsPublished && post.PublishedAt == nil {
		now := time.Now()
		post.PublishedAt = &now
	}

	database.DB.Save(&post)
	database.DB.Preload("Author").First(&post, "id = ?", post.ID)
	utils.Success(c, post)
}

func DeleteBlogPost(c *gin.Context) {
	id := c.Param("id")
	postID, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.Error(400, "invalid post id"))
		return
	}

	if err := database.DB.Delete(&models.BlogPost{}, "id = ?", postID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, utils.Error(500, "failed to delete post"))
		return
	}

	utils.Success(c, gin.H{"message": "post deleted successfully"})
}

func GetAllBlogPosts(c *gin.Context) {
	var posts []models.BlogPost
	database.DB.Preload("Author").Order("created_at DESC").Find(&posts)
	utils.Success(c, posts)
}
