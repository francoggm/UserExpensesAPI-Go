package expenses

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/francoggm/go_expenses_api/internal/users"
	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	srv    Service
	logger *zap.SugaredLogger
}

func NewHandler(s Service, logger *zap.SugaredLogger) *Handler {
	return &Handler{
		srv:    s,
		logger: logger,
	}
}

func (h *Handler) ListExpenses(c *gin.Context) {
	sessionToken, _ := c.Cookie("session_token")

	userId := users.GetIdBySession(sessionToken)
	if userId == 0 {
		h.logger.Warnw("not authenticated",
			zap.String("sessionToken", sessionToken),
			zap.Int64("userId", userId),
			zap.String("IP", c.RemoteIP()),
			zap.String("handler", "listExpenses"),
		)

		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "not authenticated",
			"data":    nil,
		})

		return
	}

	expenses, err := h.srv.ListExpenses(userId)
	if err != nil {
		h.logger.Errorw("internal error",
			zap.Error(err),
			zap.String("sessionToken", sessionToken),
			zap.Int64("userId", userId),
			zap.String("IP", c.RemoteIP()),
			zap.String("handler", "listExpenses"),
		)

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

	h.logger.Infow("success list expenses",
		zap.String("sessionToken", sessionToken),
		zap.Int64("userId", userId),
		zap.Int("expensesCount", len(expensesResponse)),
		zap.String("IP", c.RemoteIP()),
		zap.String("handler", "listExpenses"),
	)

	c.JSON(http.StatusOK, gin.H{
		"message": "success",
		"data":    expensesResponse,
	})
}

func (h *Handler) GetExpense(c *gin.Context) {
	sessionToken, _ := c.Cookie("session_token")

	userId := users.GetIdBySession(sessionToken)
	if userId == 0 {
		h.logger.Warnw("not authenticated",
			zap.String("sessionToken", sessionToken),
			zap.Int64("userId", userId),
			zap.String("IP", c.RemoteIP()),
			zap.String("handler", "getExpense"),
		)

		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "not authenticated",
			"data":    nil,
		})

		return
	}

	pathId := c.Param("id")

	id, err := strconv.ParseInt(pathId, 10, 64)
	if err != nil {
		h.logger.Errorw("internal error",
			zap.Error(err),
			zap.String("sessionToken", sessionToken),
			zap.Int64("userId", userId),
			zap.Int64("expenseId", id),
			zap.String("IP", c.RemoteIP()),
			zap.String("handler", "getExpense"),
		)

		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "internal error, please try again",
			"data":    nil,
		})

		return
	}

	expense, err := h.srv.GetExpense(id, userId)
	if err != nil {
		if err == sql.ErrNoRows {
			h.logger.Warnw("item not found",
				zap.String("sessionToken", sessionToken),
				zap.Int64("userId", userId),
				zap.Int64("expenseId", id),
				zap.String("IP", c.RemoteIP()),
				zap.String("handler", "getExpense"),
			)

			c.JSON(http.StatusNotFound, gin.H{
				"message": "item not found",
				"data":    nil,
			})

			return
		} else {
			h.logger.Errorw("internal error",
				zap.Error(err),
				zap.String("sessionToken", sessionToken),
				zap.Int64("userId", userId),
				zap.Int64("expenseId", id),
				zap.String("IP", c.RemoteIP()),
				zap.String("handler", "getExpense"),
			)

			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "internal error, please try again",
				"data":    nil,
			})

			return
		}
	}

	h.logger.Infow("success get expense",
		zap.String("sessionToken", sessionToken),
		zap.Int64("userId", userId),
		zap.Int64("expenseId", id),
		zap.String("IP", c.RemoteIP()),
		zap.String("handler", "getExpense"),
	)

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
		h.logger.Warnw("not authenticated",
			zap.String("sessionToken", sessionToken),
			zap.Int64("userId", userId),
			zap.String("IP", c.RemoteIP()),
			zap.String("handler", "createExpense"),
		)

		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "not authenticated",
			"data":    nil,
		})

		return
	}

	if err := c.BindJSON(&req); err != nil {
		h.logger.Errorw("internal error",
			zap.Error(err),
			zap.String("sessionToken", sessionToken),
			zap.Int64("userId", userId),
			zap.String("IP", c.RemoteIP()),
			zap.String("handler", "createExpense"),
		)

		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "internal error, please try again",
			"data":    nil,
		})

		return
	}

	req.UserID = userId

	if err := h.srv.CreateExpense(&req); err != nil {
		h.logger.Errorw("internal error",
			zap.Error(err),
			zap.String("sessionToken", sessionToken),
			zap.Int64("userId", userId),
			zap.Int64("expenseId", req.ID),
			zap.String("IP", c.RemoteIP()),
			zap.String("handler", "createExpense"),
		)

		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "internal error, please try again",
			"data":    nil,
		})

		return
	}

	h.logger.Infow("success create expense",
		zap.String("sessionToken", sessionToken),
		zap.Int64("userId", userId),
		zap.Int64("expenseId", req.ID),
		zap.String("IP", c.RemoteIP()),
		zap.String("handler", "getExpense"),
	)

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
		h.logger.Warnw("not authenticated",
			zap.String("sessionToken", sessionToken),
			zap.Int64("userId", userId),
			zap.String("IP", c.RemoteIP()),
			zap.String("handler", "updateExpense"),
		)

		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "not authenticated",
			"data":    nil,
		})

		return
	}

	pathId := c.Param("id")

	id, err := strconv.ParseInt(pathId, 10, 64)
	if err != nil {
		h.logger.Errorw("internal error",
			zap.Error(err),
			zap.String("sessionToken", sessionToken),
			zap.Int64("userId", userId),
			zap.Int64("expenseId", id),
			zap.String("IP", c.RemoteIP()),
			zap.String("handler", "updateExpense"),
		)

		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "internal error, please try again",
			"data":    nil,
		})

		return
	}

	expense, err := h.srv.GetExpense(id, userId)
	if err != nil {
		if err == sql.ErrNoRows {
			h.logger.Warnw("item not found",
				zap.String("sessionToken", sessionToken),
				zap.Int64("userId", userId),
				zap.Int64("expenseId", id),
				zap.String("IP", c.RemoteIP()),
				zap.String("handler", "updateExpense"),
			)

			c.JSON(http.StatusNotFound, gin.H{
				"message": "item not found",
				"data":    nil,
			})

			return
		} else {
			h.logger.Errorw("internal error",
				zap.Error(err),
				zap.String("sessionToken", sessionToken),
				zap.Int64("userId", userId),
				zap.Int64("expenseId", id),
				zap.String("IP", c.RemoteIP()),
				zap.String("handler", "updateExpense"),
			)

			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "internal error, please try again",
				"data":    nil,
			})

			return
		}
	}

	if err := c.BindJSON(&expense); err != nil {
		h.logger.Errorw("internal error",
			zap.Error(err),
			zap.String("sessionToken", sessionToken),
			zap.Int64("userId", userId),
			zap.Int64("expenseId", id),
			zap.String("IP", c.RemoteIP()),
			zap.String("handler", "updateExpense"),
		)

		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "internal error, please try again",
			"data":    nil,
		})

		return
	}

	err = h.srv.UpdateExpense(id, userId, expense)
	if err != nil {
		h.logger.Errorw("internal error",
			zap.Error(err),
			zap.String("sessionToken", sessionToken),
			zap.Int64("userId", userId),
			zap.Int64("expenseId", id),
			zap.String("IP", c.RemoteIP()),
			zap.String("handler", "updateExpense"),
		)

		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "internal error, please try again",
			"data":    nil,
		})

		return
	}

	expense, _ = h.srv.GetExpense(id, userId)

	h.logger.Infow("success update expense",
		zap.String("sessionToken", sessionToken),
		zap.Int64("userId", userId),
		zap.Int64("expenseId", id),
		zap.String("IP", c.RemoteIP()),
		zap.String("handler", "updateExpense"),
	)

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
		h.logger.Warnw("not authenticated",
			zap.String("sessionToken", sessionToken),
			zap.Int64("userId", userId),
			zap.String("IP", c.RemoteIP()),
			zap.String("handler", "deleteExpense"),
		)

		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "not authenticated",
			"data":    nil,
		})

		return
	}

	pathId := c.Param("id")

	id, err := strconv.ParseInt(pathId, 10, 64)
	if err != nil {
		h.logger.Errorw("internal error",
			zap.Error(err),
			zap.String("sessionToken", sessionToken),
			zap.Int64("userId", userId),
			zap.Int64("expenseId", id),
			zap.String("IP", c.RemoteIP()),
			zap.String("handler", "deleteExpense"),
		)

		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "internal error, please try again",
			"data":    nil,
		})

		return
	}

	_, err = h.srv.GetExpense(id, userId)
	if err != nil {
		if err == sql.ErrNoRows {
			h.logger.Warnw("item not found",
				zap.String("sessionToken", sessionToken),
				zap.Int64("userId", userId),
				zap.Int64("expenseId", id),
				zap.String("IP", c.RemoteIP()),
				zap.String("handler", "deleteExpense"),
			)

			c.JSON(http.StatusNotFound, gin.H{
				"message": "item not found",
				"data":    nil,
			})

			return
		} else {
			h.logger.Errorw("internal error",
				zap.Error(err),
				zap.String("sessionToken", sessionToken),
				zap.Int64("userId", userId),
				zap.Int64("expenseId", id),
				zap.String("IP", c.RemoteIP()),
				zap.String("handler", "deleteExpense"),
			)

			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "internal error, please try again",
				"data":    nil,
			})

			return
		}
	}

	err = h.srv.DeleteExpense(id, userId)
	if err != nil {
		h.logger.Errorw("internal error",
			zap.Error(err),
			zap.String("sessionToken", sessionToken),
			zap.Int64("userId", userId),
			zap.Int64("expenseId", id),
			zap.String("IP", c.RemoteIP()),
			zap.String("handler", "deleteExpense"),
		)

		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "internal error, please try again",
			"data":    nil,
		})

		return
	}

	h.logger.Infow("success delete expense",
		zap.String("sessionToken", sessionToken),
		zap.Int64("userId", userId),
		zap.Int64("expenseId", id),
		zap.String("IP", c.RemoteIP()),
		zap.String("handler", "deleteExpense"),
	)

	c.JSON(http.StatusOK, gin.H{
		"message": "success",
		"data":    nil,
	})
}
