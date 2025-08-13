package controller

import (
	"clean-architecture/delivary/handlerdto"
	"clean-architecture/usecase/contract"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authUseCase contract.IAuthUseCase
}

func NewAuthHandler(authUseCase contract.IAuthUseCase) *AuthHandler {
	return &AuthHandler{
		authUseCase: authUseCase,
	}
}

func (h *AuthHandler) HandleLogin(ctx *gin.Context) {
	var loginRequest handlerdto.LoginRequest
	if err := ctx.ShouldBindJSON(&loginRequest); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	aToken, rToken, err := h.authUseCase.Login(ctx.Request.Context(), loginRequest.Username, loginRequest.Password)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"access_token": aToken, "refresh_token": rToken})
}

func (h *AuthHandler) HandleLogout(ctx *gin.Context) {
	var logoutRequest handlerdto.LogoutRequest
	if err := ctx.ShouldBindJSON(&logoutRequest); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.authUseCase.Logout(ctx.Request.Context(), logoutRequest.RefreshToken)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "logged out successfully"})
}
