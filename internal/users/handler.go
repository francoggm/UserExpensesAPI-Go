package users

import (
	"database/sql"
	"errors"
	"net/http"
	"time"

	"github.com/francoggm/go_expenses_api/configs"
	"github.com/francoggm/go_expenses_api/utils"
	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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

func (h *Handler) Register(c *gin.Context) {
	var req RegisterRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Errorw("internal error",
			zap.Error(err),
			zap.String("email", req.Email),
			zap.String("name", req.Name),
			zap.String("IP", c.RemoteIP()),
			zap.String("handler", "register"),
		)

		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "error in request body",
			"data":    nil,
		})

		return
	}

	_, err := h.srv.GetUserByEmail(req.Email)
	if errors.Is(err, sql.ErrNoRows) {
		h.logger.Warnw("user already exists",
			zap.Error(err),
			zap.String("email", req.Email),
			zap.String("name", req.Name),
			zap.String("IP", c.RemoteIP()),
			zap.String("handler", "register"),
		)

		c.JSON(http.StatusConflict, gin.H{
			"message": "user already exists",
			"data":    nil,
		})

		return
	} else if err != nil {
		h.logger.Errorw("internal error",
			zap.Error(err),
			zap.String("email", req.Email),
			zap.String("name", req.Name),
			zap.String("IP", c.RemoteIP()),
			zap.String("handler", "register"),
		)

		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "internal error, please try again",
			"data":    nil,
		})

		return
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		h.logger.Errorw("internal error",
			zap.Error(err),
			zap.String("email", req.Email),
			zap.String("name", req.Name),
			zap.String("IP", c.RemoteIP()),
			zap.String("handler", "register"),
		)

		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "internal error, please try again",
			"data":    nil,
		})

		return
	}

	req.Password = hashedPassword

	user, err := h.srv.CreateUser(&req)
	if err != nil {
		h.logger.Errorw("internal error",
			zap.Error(err),
			zap.String("email", req.Email),
			zap.String("name", req.Name),
			zap.String("IP", c.RemoteIP()),
			zap.String("handler", "register"),
		)

		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "internal error, please try again",
			"data":    nil,
		})

		return
	}

	h.logger.Infow("success register user",
		zap.Int64("userId", user.ID),
		zap.String("email", user.Email),
		zap.String("name", user.Name),
		zap.String("createdAt", user.CreatedAt.String()),
		zap.String("IP", c.RemoteIP()),
		zap.String("handler", "register"),
	)

	c.JSON(http.StatusOK, gin.H{
		"message": "success",
		"data": UserResponse{
			Email:     user.Email,
			Name:      user.Name,
			LastLogin: user.LastLogin,
		},
	})
}

func (h *Handler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Errorw("internal error",
			zap.Error(err),
			zap.String("email", req.Email),
			zap.String("IP", c.RemoteIP()),
			zap.String("handler", "login"),
		)

		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "error in request body",
			"data":    nil,
		})

		return
	}

	user, err := h.srv.GetUserByEmail(req.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			h.logger.Warnw("user not found",
				zap.String("email", req.Email),
				zap.String("IP", c.RemoteIP()),
				zap.String("handler", "login"),
			)

			c.JSON(http.StatusNotFound, gin.H{
				"message": "invalid account informations",
				"data":    nil,
			})

			return
		} else {
			h.logger.Errorw("internal error",
				zap.Error(err),
				zap.String("email", req.Email),
				zap.String("IP", c.RemoteIP()),
				zap.String("handler", "login"),
			)

			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "internal error, please try again",
				"data":    nil,
			})

			return
		}
	}

	err = utils.CheckHashedPassword(user.Password, req.Password)
	if err != nil {
		h.logger.Warnw("invalid account informations",
			zap.String("email", req.Email),
			zap.String("IP", c.RemoteIP()),
			zap.String("handler", "login"),
		)

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

	c.SetCookie("session_token", sessionId, int(cfg.SessionExpires*time.Second), "/", cfg.CookieDomain, false, true)

	h.srv.SetLastLogin(user.ID, time.Now())

	h.logger.Infow("success login user",
		zap.Int64("userId", user.ID),
		zap.String("email", user.Email),
		zap.String("name", user.Name),
		zap.String("createdAt", user.CreatedAt.String()),
		zap.String("lastLogin", user.LastLogin.String()),
		zap.String("IP", c.RemoteIP()),
		zap.String("handler", "register"),
	)

	c.JSON(http.StatusOK, gin.H{
		"message": "success",
		"data": UserResponse{
			Email:     user.Email,
			Name:      user.Name,
			LastLogin: user.LastLogin,
		},
	})
}

func (h *Handler) Authenticate(c *gin.Context) {
	sessionToken, _ := c.Cookie("session_token")
	if !IsAuthenticated(sessionToken) {
		h.logger.Warnw("not authenticated",
			zap.String("sessionToken", sessionToken),
			zap.String("IP", c.RemoteIP()),
			zap.String("handler", "login"),
		)

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

	c.SetCookie("session_token", sessionId, int(cfg.SessionExpires*time.Second), "/", cfg.CookieDomain, false, true)

	h.logger.Infow("refresh session",
		zap.Int64("userId", sessions[sessionToken].userId),
		zap.String("oldSessionToken", sessionToken),
		zap.String("newSessionToken", sessionId),
		zap.String("IP", c.RemoteIP()),
		zap.String("handler", "register"),
	)

	c.JSON(http.StatusOK, gin.H{
		"message": "success",
		"data":    nil,
	})
}
