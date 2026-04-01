package handlers

import (
	"net/http"

	"caretop-backend/database"
	"caretop-backend/models"
	"caretop-backend/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type CreateTicketInput struct {
	Type    string `json:"type" binding:"required"`
	Title   string `json:"title" binding:"required,min=5,max=300"`
	Content string `json:"content" binding:"required,min=10"`
}

func GetTickets(c *gin.Context) {
	user, _ := c.Get("user")
	u := user.(*models.User)

	var tickets []models.Ticket
	database.DB.Where("user_id = ?", u.ID).Order("created_at DESC").Find(&tickets)

	utils.Success(c, tickets)
}

func GetTicket(c *gin.Context) {
	id := c.Param("id")
	ticketID, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.Error(400, "invalid ticket id"))
		return
	}

	user, _ := c.Get("user")
	u := user.(*models.User)

	var ticket models.Ticket
	if err := database.DB.First(&ticket, "id = ?", ticketID).Error; err != nil {
		c.JSON(http.StatusNotFound, utils.Error(404, "ticket not found"))
		return
	}

	if ticket.UserID != u.ID && u.Role != models.RoleAdmin {
		c.JSON(http.StatusForbidden, utils.Error(403, "not authorized to view this ticket"))
		return
	}

	var replies []models.TicketReply
	database.DB.Where("ticket_id = ?", ticketID).Preload("User").Order("created_at ASC").Find(&replies)

	utils.Success(c, gin.H{
		"ticket":  ticket,
		"replies": replies,
	})
}

func CreateTicket(c *gin.Context) {
	var input CreateTicketInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, utils.Error(400, "invalid input: "+err.Error()))
		return
	}

	user, _ := c.Get("user")
	u := user.(*models.User)

	ticket := models.Ticket{
		UserID:  u.ID,
		Type:    models.TicketType(input.Type),
		Title:   input.Title,
		Content: input.Content,
		Status:  models.TicketStatusPending,
	}

	if err := database.DB.Create(&ticket).Error; err != nil {
		c.JSON(http.StatusInternalServerError, utils.Error(500, "failed to create ticket"))
		return
	}

	database.DB.Preload("User").First(&ticket, "id = ?", ticket.ID)
	utils.Success(c, ticket)
}

type ReplyTicketInput struct {
	Content string `json:"content" binding:"required,min=1"`
}

func ReplyTicket(c *gin.Context) {
	id := c.Param("id")
	ticketID, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.Error(400, "invalid ticket id"))
		return
	}

	user, _ := c.Get("user")
	u := user.(*models.User)

	var ticket models.Ticket
	if err := database.DB.First(&ticket, "id = ?", ticketID).Error; err != nil {
		c.JSON(http.StatusNotFound, utils.Error(404, "ticket not found"))
		return
	}

	if ticket.UserID != u.ID && u.Role != models.RoleAdmin {
		c.JSON(http.StatusForbidden, utils.Error(403, "not authorized to reply to this ticket"))
		return
	}

	var input ReplyTicketInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, utils.Error(400, "invalid input: "+err.Error()))
		return
	}

	reply := models.TicketReply{
		TicketID: ticketID,
		UserID:   u.ID,
		Content:  input.Content,
	}

	if err := database.DB.Create(&reply).Error; err != nil {
		c.JSON(http.StatusInternalServerError, utils.Error(500, "failed to create reply"))
		return
	}

	if ticket.Status == models.TicketStatusPending {
		database.DB.Model(&ticket).Update("status", models.TicketStatusProcessing)
	}

	database.DB.Preload("User").First(&reply, "id = ?", reply.ID)
	utils.Success(c, reply)
}

func GetAllTickets(c *gin.Context) {
	var tickets []models.Ticket
	database.DB.Preload("User").Order("created_at DESC").Find(&tickets)
	utils.Success(c, tickets)
}

type UpdateTicketStatusInput struct {
	Status string `json:"status" binding:"required"`
}

func UpdateTicketStatus(c *gin.Context) {
	id := c.Param("id")
	ticketID, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.Error(400, "invalid ticket id"))
		return
	}

	var input UpdateTicketStatusInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, utils.Error(400, "invalid input: "+err.Error()))
		return
	}

	var ticket models.Ticket
	if err := database.DB.First(&ticket, "id = ?", ticketID).Error; err != nil {
		c.JSON(http.StatusNotFound, utils.Error(404, "ticket not found"))
		return
	}

	ticket.Status = models.TicketStatus(input.Status)
	database.DB.Save(&ticket)

	utils.Success(c, ticket)
}
