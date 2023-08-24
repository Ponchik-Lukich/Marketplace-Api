package segment

import (
	"github.com/gin-gonic/gin"
	"market/pkg/repository/segment"
)

type IHandler interface {
	CreateSegments(ctx *gin.Context)
	DeleteSegments(ctx *gin.Context)
}

type Handler struct {
	repo segment.IRepository
}

func NewHandler(segmentRepo segment.IRepository) *Handler {
	return &Handler{repo: segmentRepo}
}

func (h *Handler) CreateSegments(ctx *gin.Context) {
	ctx.JSON(200, gin.H{
		"message": "pong",
	})
}

func (h *Handler) DeleteSegments(ctx *gin.Context) {
	ctx.JSON(200, gin.H{
		"message": "pong",
	})
}
