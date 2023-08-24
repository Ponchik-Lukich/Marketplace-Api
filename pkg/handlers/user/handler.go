package user

import (
	"github.com/gin-gonic/gin"
	"market/pkg/repository/user"
)

type IHandler interface {
	GetUsersSlugs(ctx *gin.Context)
	EditUsersSlugs(ctx *gin.Context)
}

type Handler struct {
	repo user.IRepository
}

func NewHandler(userRepo user.IRepository) *Handler {
	return &Handler{repo: userRepo}
}

func (h *Handler) GetUsersSlugs(ctx *gin.Context) {
	ctx.JSON(200, gin.H{
		"message": "pong",
	})
}

func (h *Handler) EditUsersSlugs(ctx *gin.Context) {
	ctx.JSON(200, gin.H{
		"message": "pong",
	})
}
