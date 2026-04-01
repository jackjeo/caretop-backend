package handlers

import (
	"net/http"
	"strconv"

	"caretop-backend/database"
	"caretop-backend/models"
	"caretop-backend/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ThreadListResponse struct {
	Threads    []models.ForumThread `json:"threads"`
	Total      int64                 `json:"total"`
	Page       int                   `json:"page"`
	PageSize   int                   `json:"page_size"`
	TotalPages int                   `json:"total_pages"`
}

func GetBoards(c *gin.Context) {
	var boards []models.ForumBoard
	database.DB.Order("sort_order ASC").Find(&boards)
	utils.Success(c, boards)
}

func GetBoardThreads(c *gin.Context) {
	slug := c.Param("slug")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize := 20
	filter := c.DefaultQuery("filter", "latest")

	if page < 1 {
		page = 1
	}

	var board models.ForumBoard
	if err := database.DB.First(&board, "slug = ?", slug).Error; err != nil {
		c.JSON(http.StatusNotFound, utils.Error(404, "board not found"))
		return
	}

	var threads []models.ForumThread
	var total int64

	query := database.DB.Model(&models.ForumThread{}).Where("board_id = ?", board.ID)

	switch filter {
	case "hot":
		query = query.Order("view_count DESC, created_at DESC")
	case "essential":
		query = query.Where("is_essential = ?", true).Order("created_at DESC")
	default:
		query = query.Order("is_pinned DESC, created_at DESC")
	}

	query.Count(&total)
	query.Offset((page - 1) * pageSize).Limit(pageSize).Preload("User").Find(&threads)

	totalPages := int(total) / pageSize
	if int(total)%pageSize > 0 {
		totalPages++
	}

	utils.Success(c, ThreadListResponse{
		Threads:    threads,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	})
}

func GetThread(c *gin.Context) {
	id := c.Param("id")
	threadID, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.Error(400, "invalid thread id"))
		return
	}

	var thread models.ForumThread
	if err := database.DB.Preload("User").Preload("Board").First(&thread, "id = ?", threadID).Error; err != nil {
		c.JSON(http.StatusNotFound, utils.Error(404, "thread not found"))
		return
	}

	database.DB.Model(&thread).UpdateColumn("view_count", thread.ViewCount+1)
	thread.ViewCount++

	var posts []models.ForumPost
	database.DB.Where("thread_id = ?", threadID).Preload("User").Order("created_at ASC").Find(&posts)

	utils.Success(c, gin.H{
		"thread": thread,
		"posts":  posts,
	})
}

type CreateThreadInput struct {
	BoardID string `json:"board_id" binding:"required"`
	Title   string `json:"title" binding:"required,min=5,max=300"`
	Content string `json:"content" binding:"required,min=10"`
}

func CreateThread(c *gin.Context) {
	var input CreateThreadInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, utils.Error(400, "invalid input: "+err.Error()))
		return
	}

	user, _ := c.Get("user")
	u := user.(*models.User)

	boardID, err := uuid.Parse(input.BoardID)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.Error(400, "invalid board id"))
		return
	}

	thread := models.ForumThread{
		BoardID: boardID,
		UserID:  u.ID,
		Title:   input.Title,
		Content: input.Content,
	}

	if err := database.DB.Create(&thread).Error; err != nil {
		c.JSON(http.StatusInternalServerError, utils.Error(500, "failed to create thread"))
		return
	}

	database.DB.Preload("User").Preload("Board").First(&thread, "id = ?", thread.ID)
	utils.Success(c, thread)
}

type UpdateThreadInput struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

func UpdateThread(c *gin.Context) {
	id := c.Param("id")
	threadID, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.Error(400, "invalid thread id"))
		return
	}

	user, _ := c.Get("user")
	u := user.(*models.User)

	var thread models.ForumThread
	if err := database.DB.First(&thread, "id = ?", threadID).Error; err != nil {
		c.JSON(http.StatusNotFound, utils.Error(404, "thread not found"))
		return
	}

	if thread.UserID != u.ID && u.Role != models.RoleAdmin {
		c.JSON(http.StatusForbidden, utils.Error(403, "not authorized to edit this thread"))
		return
	}

	var input UpdateThreadInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, utils.Error(400, "invalid input: "+err.Error()))
		return
	}

	if input.Title != "" {
		thread.Title = input.Title
	}
	if input.Content != "" {
		thread.Content = input.Content
	}

	database.DB.Save(&thread)
	database.DB.Preload("User").Preload("Board").First(&thread, "id = ?", thread.ID)
	utils.Success(c, thread)
}

func DeleteThread(c *gin.Context) {
	id := c.Param("id")
	threadID, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.Error(400, "invalid thread id"))
		return
	}

	user, _ := c.Get("user")
	u := user.(*models.User)

	var thread models.ForumThread
	if err := database.DB.First(&thread, "id = ?", threadID).Error; err != nil {
		c.JSON(http.StatusNotFound, utils.Error(404, "thread not found"))
		return
	}

	if thread.UserID != u.ID && u.Role != models.RoleAdmin {
		c.JSON(http.StatusForbidden, utils.Error(403, "not authorized to delete this thread"))
		return
	}

	database.DB.Delete(&models.ForumThread{}, "id = ?", threadID)
	utils.Success(c, gin.H{"message": "thread deleted successfully"})
}

type ReplyInput struct {
	Content  string     `json:"content" binding:"required,min=1"`
	ParentID *uuid.UUID `json:"parent_id"`
}

func ReplyThread(c *gin.Context) {
	id := c.Param("id")
	threadID, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.Error(400, "invalid thread id"))
		return
	}

	var thread models.ForumThread
	if err := database.DB.First(&thread, "id = ?", threadID).Error; err != nil {
		c.JSON(http.StatusNotFound, utils.Error(404, "thread not found"))
		return
	}

	if thread.IsLocked {
		c.JSON(http.StatusForbidden, utils.Error(403, "thread is locked"))
		return
	}

	user, _ := c.Get("user")
	u := user.(*models.User)

	var input ReplyInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, utils.Error(400, "invalid input: "+err.Error()))
		return
	}

	post := models.ForumPost{
		ThreadID: threadID,
		UserID:   u.ID,
		ParentID: input.ParentID,
		Content:  input.Content,
	}

	if err := database.DB.Create(&post).Error; err != nil {
		c.JSON(http.StatusInternalServerError, utils.Error(500, "failed to create reply"))
		return
	}

	database.DB.Preload("User").First(&post, "id = ?", post.ID)
	utils.Success(c, post)
}

func LikeThread(c *gin.Context) {
	id := c.Param("id")
	threadID, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.Error(400, "invalid thread id"))
		return
	}

	user, _ := c.Get("user")
	u := user.(*models.User)

	var existingLike models.ForumLike
	if err := database.DB.First(&existingLike, "user_id = ? AND thread_id = ?", u.ID, threadID).Error; err == nil {
		c.JSON(http.StatusBadRequest, utils.Error(400, "already liked this thread"))
		return
	}

	like := models.ForumLike{
		UserID:   u.ID,
		ThreadID: threadID,
	}

	if err := database.DB.Create(&like).Error; err != nil {
		c.JSON(http.StatusInternalServerError, utils.Error(500, "failed to like thread"))
		return
	}

	database.DB.Model(&models.ForumThread{}).Where("id = ?", threadID).UpdateColumn("like_count", database.DB.Raw("like_count + 1"))

	utils.Success(c, gin.H{"message": "liked successfully"})
}

func CollectThread(c *gin.Context) {
	id := c.Param("id")
	threadID, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.Error(400, "invalid thread id"))
		return
	}

	user, _ := c.Get("user")
	u := user.(*models.User)

	var existingCollection models.ForumCollection
	if err := database.DB.First(&existingCollection, "user_id = ? AND thread_id = ?", u.ID, threadID).Error; err == nil {
		c.JSON(http.StatusBadRequest, utils.Error(400, "already collected this thread"))
		return
	}

	collection := models.ForumCollection{
		UserID:   u.ID,
		ThreadID: threadID,
	}

	if err := database.DB.Create(&collection).Error; err != nil {
		c.JSON(http.StatusInternalServerError, utils.Error(500, "failed to collect thread"))
		return
	}

	utils.Success(c, gin.H{"message": "collected successfully"})
}

func GetMyPosts(c *gin.Context) {
	user, _ := c.Get("user")
	u := user.(*models.User)

	var threads []models.ForumThread
	database.DB.Where("user_id = ?", u.ID).Preload("Board").Order("created_at DESC").Find(&threads)

	var collections []models.ForumThread
	database.DB.Model(&models.ForumCollection{}).Where("user_id = ?", u.ID).Preload("Thread.Board").Find(&collections)

	var collThreads []models.ForumThread
	for _, coll := range collections {
		var t models.ForumThread
		database.DB.Preload("Board").First(&t, "id = ?", coll.ID)
		collThreads = append(collThreads, t)
	}

	utils.Success(c, gin.H{
		"threads":      threads,
		"collections":  collThreads,
	})
}
