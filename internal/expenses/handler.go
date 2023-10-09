package expenses

import (
	"database/sql"
	"expenses_api/internal/users"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	srv Service
}

func NewHandler(s Service) *Handler {
	return &Handler{
		srv: s,
	}
}

func (h *Handler) GetExpenses(c *gin.Context) {
	sessionToken, _ := c.Cookie("session_token")

	userId := users.GetIdBySession(sessionToken)
	if userId == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "not authenticated"})
		return
	}

	expenses, err := h.srv.GetExpenses(userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, expenses)
}

func (h *Handler) GetExpense(c *gin.Context) {
	sessionToken, _ := c.Cookie("session_token")

	userId := users.GetIdBySession(sessionToken)
	if userId == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "not authenticated"})
		return
	}

	pathId := c.Param("id")
	id, err := strconv.ParseInt(pathId, 10, 64)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	expense, err := h.srv.GetExpense(id, userId)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
			return
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	c.JSON(http.StatusOK, expense)
}

func (h *Handler) CreateExpense(c *gin.Context) {
	var req Expense

	sessionToken, _ := c.Cookie("session_token")

	userId := users.GetIdBySession(sessionToken)
	if userId == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": ""})
		return
	}

	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	req.UserID = userId

	if err := h.srv.CreateExpense(&req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, ExpenseResponse{
		ID:            req.ID,
		Title:         req.Title,
		Description:   req.Description,
		Value:         req.Value,
		Category:      req.Category,
		Movimentation: req.Movimentation,
	})
}

func (h *Handler) DeleteExpense(c *gin.Context) {
	sessionToken, _ := c.Cookie("session_token")

	userId := users.GetIdBySession(sessionToken)
	if userId == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": ""})
		return
	}

	pathId := c.Param("id")
	id, err := strconv.ParseInt(pathId, 10, 64)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	err = h.srv.DeleteExpense(id, userId)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
			return
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	c.Status(http.StatusOK)
}
