package controller

import (
	"clean-architecture/delivary/handlerdto"
	"clean-architecture/usecase/contract"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	UserUseCase contract.IUserUseCase
}

func NewUserHandler(userUC contract.IUserUseCase) *UserHandler {
	return &UserHandler{
		UserUseCase: userUC,
	}
}
func (h *UserHandler) HandleRegister(ctx *gin.Context) {
	// instantiate UserDto
	var userDto handlerdto.UserDto
	// bind the request body to userdto
	if err := ctx.ShouldBindJSON(&userDto); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	user, err := h.UserUseCase.Register(ctx.Request.Context(), userDto.UserID, userDto.Username, userDto.UserType, userDto.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, gin.H{"user": user})
}

func (h *UserHandler) UpdateUser(ctx *gin.Context) {
	var userDto handlerdto.UserDto
	if err := ctx.ShouldBindJSON(&userDto); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updatedUser, err := h.UserUseCase.UpdateUser(ctx.Request.Context(), userDto.UserID, userDto.UserType, userDto.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"user": updatedUser})
}

func (h *UserHandler) DeleteUser(ctx *gin.Context) {
	id := ctx.Param("id")
	if err := h.UserUseCase.DeleteUser(ctx.Request.Context(), id); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusNoContent, nil)
}
