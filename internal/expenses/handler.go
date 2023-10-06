package expenses

import (
	"expenses_api/internal/users"
	"net/http"

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

func (h *Handler) CreateExpense(c *gin.Context) {
	var req Expense

	sessionToken, _ := c.Cookie("session_token")
	if !users.IsAuthenticated(sessionToken) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "please log in to continue!"})
		return
	}

	userId := users.GetIdBySession(sessionToken)
	if userId == -1 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": ""})
		return
	}

	req.UserID = userId

	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := h.srv.CreateExpense(&req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, req)
}
