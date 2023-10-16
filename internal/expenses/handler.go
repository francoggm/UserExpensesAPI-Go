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
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "not authenticated",
			"data":    nil,
		})

		return
	}

	expenses, err := h.srv.ListExpenses(userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "internal error, please try again",
			"data":    nil,
		})

		return
	}

	var expensesResponse = []ExpenseResponse{}

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

	c.JSON(http.StatusOK, gin.H{
		"message": "success",
		"data":    expensesResponse,
	})
}

func (h *Handler) GetExpense(c *gin.Context) {
	sessionToken, _ := c.Cookie("session_token")

	userId := users.GetIdBySession(sessionToken)
	if userId == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "not authenticated",
			"data":    nil,
		})

		return
	}

	pathId := c.Param("id")
	id, err := strconv.ParseInt(pathId, 10, 64)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "internal error, please try again",
			"data":    nil,
		})

		return
	}

	expense, err := h.srv.GetExpense(id, userId)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{
				"message": "item not found",
				"data":    nil,
			})

			return
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "internal error, please try again",
				"data":    nil,
			})

			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "success",
		"data": ExpenseResponse{
			ID:            expense.ID,
			Title:         expense.Title,
			Description:   expense.Description,
			Value:         expense.Value,
			Category:      expense.Category,
			Movimentation: expense.Movimentation,
			CreatedAt:     expense.CreatedAt,
		},
	})
}

func (h *Handler) CreateExpense(c *gin.Context) {
	var req Expense

	sessionToken, _ := c.Cookie("session_token")

	userId := users.GetIdBySession(sessionToken)
	if userId == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "not authenticated",
			"data":    nil,
		})

		return
	}

	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "internal error, please try again",
			"data":    nil,
		})

		return
	}

	req.UserID = userId

	if err := h.srv.CreateExpense(&req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "internal error, please try again",
			"data":    nil,
		})

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "success",
		"data": ExpenseResponse{
			ID:            req.ID,
			Title:         req.Title,
			Description:   req.Description,
			Value:         req.Value,
			Category:      req.Category,
			Movimentation: req.Movimentation,
			CreatedAt:     req.CreatedAt,
		},
	})
}

func (h *Handler) UpdateExpense(c *gin.Context) {
	sessionToken, _ := c.Cookie("session_token")

	userId := users.GetIdBySession(sessionToken)
	if userId == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "not authenticated",
			"data":    nil,
		})

		return
	}

	pathId := c.Param("id")
	id, err := strconv.ParseInt(pathId, 10, 64)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "internal error, please try again",
			"data":    nil,
		})

		return
	}

	expense, err := h.srv.GetExpense(id, userId)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{
				"message": "item not found",
				"data":    nil,
			})

			return
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "internal error, please try again",
				"data":    nil,
			})

			return
		}
	}

	if err := c.BindJSON(&expense); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "internal error, please try again",
			"data":    nil,
		})

		return
	}

	err = h.srv.UpdateExpense(id, userId, expense)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "internal error, please try again",
			"data":    nil,
		})

		return
	}

	expense, _ = h.srv.GetExpense(id, userId)

	c.JSON(http.StatusOK, gin.H{
		"message": "success",
		"data": ExpenseResponse{
			ID:            expense.ID,
			Title:         expense.Title,
			Description:   expense.Description,
			Value:         expense.Value,
			Category:      expense.Category,
			Movimentation: expense.Movimentation,
			CreatedAt:     expense.CreatedAt,
		},
	})
}

func (h *Handler) DeleteExpense(c *gin.Context) {
	sessionToken, _ := c.Cookie("session_token")

	userId := users.GetIdBySession(sessionToken)
	if userId == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "not authenticated",
			"data":    nil,
		})

		return
	}

	pathId := c.Param("id")
	id, err := strconv.ParseInt(pathId, 10, 64)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "internal error, please try again",
			"data":    nil,
		})

		return
	}

	_, err = h.srv.GetExpense(id, userId)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{
				"message": "item not found",
				"data":    nil,
			})

			return
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "internal error, please try again",
				"data":    nil,
			})

			return
		}
	}

	err = h.srv.DeleteExpense(id, userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "internal error, please try again",
			"data":    nil,
		})

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "success",
		"data":    nil,
	})
}
