package user

import (
	"github.com/gin-gonic/gin"
	"market/pkg/dtos"
	"market/pkg/errors"
	"market/pkg/repository/user"
	"net/http"
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
	var payload dtos.EditUserDtoRequest

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		errors.HandleError(ctx, http.StatusInternalServerError, errors.BindingJSONErr, err)
		return
	}

	if payload.UserID == 0 {
		errors.HandleError(ctx, http.StatusBadRequest, errors.EmptyUserIDErr, nil)
		return
	}

	if err := h.repo.EditUsersSlugs(payload.ToCreate, payload.ToDelete, payload.UserID); err != nil {
		errors.HandleError(ctx, http.StatusBadRequest, errors.EditingUserErr, err)
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"message": "success",
	})
}
