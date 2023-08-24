package segment

import (
	"github.com/gin-gonic/gin"
	"market/pkg/dtos"
	"market/pkg/errors"
	"market/pkg/repository/segment"
	"net/http"
)

type IHandler interface {
	CreateSegment(ctx *gin.Context)
	DeleteSegment(ctx *gin.Context)
}

type Handler struct {
	repo segment.IRepository
}

func NewHandler(segmentRepo segment.IRepository) *Handler {
	return &Handler{repo: segmentRepo}
}

func (h *Handler) CreateSegment(ctx *gin.Context) {
	var payload dtos.SegmentDtoRequest

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		errors.HandleError(ctx, http.StatusInternalServerError, errors.BindingJSONErr, err)
		return
	}

	if payload.Name == "" {
		errors.HandleError(ctx, http.StatusBadRequest, errors.EmptySegmentNameErr, nil)
		return
	}

	if err := h.repo.CreateSegment(payload.Name); err != nil {
		errors.HandleError(ctx, http.StatusBadRequest, errors.CreatingSegmentErr, err)
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"message": "success",
	})
}

func (h *Handler) DeleteSegment(ctx *gin.Context) {
	var payload dtos.SegmentDtoRequest

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		errors.HandleError(ctx, http.StatusInternalServerError, errors.BindingJSONErr, err)
		return
	}

	if payload.Name == "" {
		errors.HandleError(ctx, http.StatusBadRequest, errors.EmptySegmentNameErr, nil)
		return
	}

	if err := h.repo.DeleteSegment(payload.Name); err != nil {
		errors.HandleError(ctx, http.StatusBadRequest, errors.DeletingSegmentErr, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "success",
	})
}
