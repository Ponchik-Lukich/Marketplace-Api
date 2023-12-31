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
	CreateUserLogs(ctx *gin.Context)
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

	segments, customErr := h.repo.GetUsersSegments(userID)
	if customErr != nil {
		errors.HandleCustomError(ctx, customErr.Code(), customErr)
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
		errors.HandleCustomError(ctx, err.Code(), err)
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"message": "success",
	})
}

func (h *Handler) CreateUserLogs(ctx *gin.Context) {
	var payload dtos.GetLogsDtoRequest

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		errors.HandleError(ctx, http.StatusBadRequest, errors.BindingJSONErr, err)
		return
	}

	if payload.UserID == 0 {
		errors.HandleError(ctx, http.StatusBadRequest, errors.EmptyUserIDErr, nil)
		return
	}

	if payload.Date == "" {
		errors.HandleError(ctx, http.StatusBadRequest, errors.EmptyDateErr, nil)
		return
	}

	filePath, err := h.repo.CreateUserLogs(payload.Date, payload.UserID)
	if err != nil {
		errors.HandleCustomError(ctx, err.Code(), err)
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"file_path": filePath,
	})
}
