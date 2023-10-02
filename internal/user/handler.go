package user

import (
	"database/sql"
	"expenses_api/utils"
	"net/http"
	"time"

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

func (h *Handler) Register(c *gin.Context) {
	var req User
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	req.Password = hashedPassword

	err = h.srv.CreateUser(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, UserResponse{
		ID:        req.ID,
		Email:     req.Email,
		Name:      req.Name,
		CreatedAt: req.CreatedAt,
		LastLogin: req.LastLogin,
	})
}

func (h *Handler) Login(c *gin.Context) {
	var req User
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	user, err := h.srv.GetUserByEmail(req.Email)
	if err != nil {
		statusCode := 0

		if err == sql.ErrNoRows {
			statusCode = http.StatusNotFound
		} else {	
			statusCode = http.StatusInternalServerError
		}

		c.JSON(statusCode, gin.H{"error": err.Error()})
		return
	}

	err = utils.CheckHashedPassword(user.Password, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.SetCookie("email", user.Email, int(4*time.Hour), "/", "localhost", false, true)
	c.SetCookie("password", user.Password, int(4*time.Hour), "/", "localhost", false, true)

	c.JSON(http.StatusOK, UserResponse{
		ID:        user.ID,
		Email:     user.Email,
		Name:      user.Name,
		CreatedAt: user.CreatedAt,
		LastLogin: user.LastLogin,
	})
}
