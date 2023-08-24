package user

import (
	"github.com/gin-gonic/gin"
	"market/pkg/dtos"
	"market/pkg/errors"
	"market/pkg/repository/user"
	"net/http"
	"strconv"
)

type IHandler interface {
	GetUsersSegments(ctx *gin.Context)
	EditUsersSegments(ctx *gin.Context)
}

type Handler struct {
	repo user.IRepository
}

func NewHandler(userRepo user.IRepository) *Handler {
	return &Handler{repo: userRepo}
}

func (h *Handler) GetUsersSegments(ctx *gin.Context) {
	userIdStr, ok := ctx.GetQuery("user_id")
	if !ok {
		errors.HandleError(ctx, http.StatusBadRequest, errors.EmptyUserIDErr, nil)
		return
	}

	userID, err := strconv.ParseUint(userIdStr, 10, 64)
	if err != nil {
		errors.HandleError(ctx, http.StatusBadRequest, errors.ConvertingUserIdErr, err)
		return
	}

	segments, err := h.repo.GetUsersSegments(userID)
	if err != nil {
		errors.HandleError(ctx, http.StatusBadRequest, errors.GettingUserSegmentsErr, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"segments": segments,
	})

}

func (h *Handler) EditUsersSegments(ctx *gin.Context) {
	var payload dtos.EditUserDtoRequest

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		errors.HandleError(ctx, http.StatusInternalServerError, errors.BindingJSONErr, err)
		return
	}

	if payload.UserID == 0 {
		errors.HandleError(ctx, http.StatusBadRequest, errors.EmptyUserIDErr, nil)
		return
	}

	if err := h.repo.EditUsersSegments(payload.ToCreate, payload.ToDelete, payload.UserID); err != nil {
		errors.HandleError(ctx, http.StatusBadRequest, errors.EditingUserErr, err)
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"message": "success",
	})
}
