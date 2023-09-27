package user

import (
	"github.com/gin-gonic/gin"
)

type handler struct {
	svc Service
}

func NewHandler(s Service) *handler {
	return &handler{
		svc: s,
	}
}

func (h *handler) RegisterUser(c *gin.Context) {
}

func (h *handler) LoginUser(c *gin.Context) {
}
