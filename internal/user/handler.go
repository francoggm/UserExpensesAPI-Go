package user

import (
	"github.com/gin-gonic/gin"
)

type handler struct {
	*service
}

func NewHandler(s *service) *handler {
	return &handler{
		service: s,
	}
}

func (h *handler) RegisterUser(c *gin.Context) {
	return
}

func (h *handler) LoginUser(c *gin.Context) {
	return
}
