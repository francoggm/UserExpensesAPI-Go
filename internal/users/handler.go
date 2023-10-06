package users

import (
	"database/sql"
	"expenses_api/utils"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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
		ID:    req.ID,
		Email: req.Email,
		Name:  req.Name,
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

	sessionId, _ := c.Cookie("session_token")
	delete(sessions, sessionId)

	sessionId = uuid.New().String()
	sessions[sessionId] = session{
		userId:  user.ID,
		expires: time.Now().Add(2 * time.Hour),
	}

	c.SetCookie("session_token", sessionId, int(2*time.Hour), "/", "localhost", false, true)

	c.JSON(http.StatusOK, UserResponse{
		ID:    user.ID,
		Email: user.Email,
		Name:  user.Name,
	})
}
