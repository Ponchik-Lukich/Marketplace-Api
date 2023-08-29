package segment

import (
	"github.com/gin-gonic/gin"
	"market/pkg/dtos"
	"market/pkg/errors"
	"market/pkg/service/segment"
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
		errors.HandleError(ctx, http.StatusBadRequest, errors.EmptySegmentNameErr400, nil)
		return
	}

	if payload.Percent < 0 || payload.Percent > 100 {
		errors.HandleError(ctx, http.StatusBadRequest, errors.InvalidPercentErr, nil)
		return
	}

	if err := h.repo.CreateSegment(payload.Name, payload.Percent); err != nil {
		errors.HandleCustomError(ctx, err.Code(), err)
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
		errors.HandleError(ctx, http.StatusBadRequest, errors.EmptySegmentNameErr400, nil)
		return
	}

	if err := h.repo.DeleteSegment(payload.Name); err != nil {
		errors.HandleCustomError(ctx, err.Code(), err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "success",
	})
}
