package expenses

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/francoggm/go_expenses_api/internal/users"

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

func (h *Handler) ListExpenses(c *gin.Context) {
	sessionToken, _ := c.Cookie("session_token")

	userId := users.GetIdBySession(sessionToken)
	if userId == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "not authenticated"})
		return
	}

	expenses, err := h.srv.ListExpenses(userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var expensesResponse []ExpenseResponse

	for _, expense := range expenses {
		expensesResponse = append(expensesResponse, ExpenseResponse{
			ID:            expense.ID,
			Title:         expense.Title,
			Description:   expense.Description,
			Value:         expense.Value,
			Category:      expense.Category,
			Movimentation: expense.Movimentation,
			CreatedAt:     expense.CreatedAt,
		})
	}

	c.JSON(http.StatusOK, expensesResponse)
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

	c.JSON(http.StatusOK, ExpenseResponse{
		ID:            expense.ID,
		Title:         expense.Title,
		Description:   expense.Description,
		Value:         expense.Value,
		Category:      expense.Category,
		Movimentation: expense.Movimentation,
		CreatedAt:     expense.CreatedAt,
	})
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
		CreatedAt:     req.CreatedAt,
	})
}

func (h *Handler) UpdateExpense(c *gin.Context) {
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

	if err := c.BindJSON(&expense); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	err = h.srv.UpdateExpense(id, userId, expense)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	expense, _ = h.srv.GetExpense(id, userId)

	c.JSON(http.StatusOK, ExpenseResponse{
		ID:            expense.ID,
		Title:         expense.Title,
		Description:   expense.Description,
		Value:         expense.Value,
		Category:      expense.Category,
		Movimentation: expense.Movimentation,
		CreatedAt:     expense.CreatedAt,
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

	_, err = h.srv.GetExpense(id, userId)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
			return
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	err = h.srv.DeleteExpense(id, userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}
