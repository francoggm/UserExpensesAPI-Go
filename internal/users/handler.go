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
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
			return
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
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
		expires: time.Now().Add(30 * time.Minute),
	}

	c.SetCookie("session_token", sessionId, int(2*time.Hour), "/", "localhost", false, true)

	h.srv.SetLastLogin(user.ID, time.Now())

	c.JSON(http.StatusOK, UserResponse{
		Email: user.Email,
		Name:  user.Name,
	})
}

func (h *Handler) Authenticate(c *gin.Context) {
	sessionToken, _ := c.Cookie("session_token")
	if !IsAuthenticated(sessionToken) {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "please login to continue!"})
		return
	}
}
