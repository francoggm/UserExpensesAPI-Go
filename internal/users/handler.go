package users

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/francoggm/go_expenses_api/configs"
	"github.com/francoggm/go_expenses_api/utils"

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
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "internal error, please try again",
			"data":    nil,
		})

		return
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "internal error, please try again",
			"data":    nil,
		})

		return
	}

	req.Password = hashedPassword

	err = h.srv.CreateUser(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "internal error, please try again",
			"data":    nil,
		})

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "success",
		"data": UserResponse{
			Email: req.Email,
			Name:  req.Name,
		},
	})
}

func (h *Handler) Login(c *gin.Context) {
	var req User
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "internal error, please try again",
			"data":    nil,
		})

		return
	}

	user, err := h.srv.GetUserByEmail(req.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "invalid account informations",
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

	err = utils.CheckHashedPassword(user.Password, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "invalid account informations",
			"data":    nil,
		})

		return
	}

	sessionId, _ := c.Cookie("session_token")
	delete(sessions, sessionId)

	cfg := configs.GetConfigs()

	sessionId = uuid.New().String()
	sessions[sessionId] = session{
		userId:  user.ID,
		expires: time.Now().Add(cfg.SessionExpires * time.Second),
	}

	c.SetCookie("session_token", sessionId, int(cfg.SessionExpires*time.Second), "/", "localhost", false, true)

	h.srv.SetLastLogin(user.ID, time.Now())

	c.JSON(http.StatusOK, gin.H{
		"message": "success",
		"data": UserResponse{
			Email: user.Email,
			Name:  user.Name,
		},
	})
}

func (h *Handler) Authenticate(c *gin.Context) {
	sessionToken, _ := c.Cookie("session_token")
	if !IsAuthenticated(sessionToken) {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"message": "please login to continue!",
			"data":    nil,
		})
		return
	}
}

func (h *Handler) RefreshSession(c *gin.Context) {
	sessionToken, _ := c.Cookie("session_token")

	cfg := configs.GetConfigs()

	sessionId := uuid.New().String()
	sessions[sessionId] = session{
		userId:  sessions[sessionToken].userId,
		expires: time.Now().Add(cfg.SessionExpires * time.Second),
	}

	delete(sessions, sessionToken)

	c.SetCookie("session_token", sessionId, int(cfg.SessionExpires*time.Second), "/", "localhost", false, true)

	c.JSON(http.StatusOK, gin.H{
		"message": "success",
		"data":    nil,
	})
}
