package user

import (
	"expenses_api/utils"
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

func (h *Handler) Register(c *gin.Context) {
	var newUser User
	if err := c.BindJSON(&newUser); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	hashedPassword, err := utils.HashPassword(newUser.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	newUser.Password = hashedPassword

	userId, err := h.srv.CreateUser(&newUser); if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})	
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"email": newUser.Email,
		"id": userId,
	})
}

func (h *Handler) Login(c *gin.Context) {
}
