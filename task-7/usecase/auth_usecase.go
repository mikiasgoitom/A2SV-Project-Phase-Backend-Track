package usecase

import (
	"clean-architecture/usecase/contract"
	"context"
	"errors"
	"fmt"
)

const authTag string = "AUTH"
const internalServerError = "internal server error"

type AuthUseCase struct {
	UserRepository  contract.IUserRepository
	PasswordService contract.IPasswordService
	TokenService    contract.ITokenService
	LoggerService   contract.ILogger
}

func NewAuthUseCase(repo contract.IUserRepository, ps contract.IPasswordService, ts contract.ITokenService, lg contract.ILogger) *AuthUseCase {
	return &AuthUseCase{
		UserRepository:  repo,
		PasswordService: ps,
		TokenService:    ts,
		LoggerService:   lg,
	}
}

func (u *AuthUseCase) Login(ctx context.Context, username, password string) (accessToken string, refreshToken string, err error) {
	foundUser, err := u.UserRepository.GetUserByUsername(ctx, username)
	if err != nil {
		u.LoggerService.Secure(authTag, fmt.Sprintf("failed to retrive user by username: %s", err))
		return "", "", errors.New("wrong username or password")
	}

	passwordMatched, err := u.PasswordService.ComparePassword(password, foundUser.Password)
	if err != nil {
		u.LoggerService.Secure(authTag, fmt.Sprintf("failed to match password: %s", err))
		return "", "", errors.New("wrong username or password")
	}
	// check if the password is correct
	if !passwordMatched {
		return "", "", errors.New("wrong username or password")
	}
	accessToken, refreshToken, err = u.TokenService.GenerateTokens(ctx, foundUser.UserID, foundUser.Username, string(foundUser.UserType))

	if err != nil {
		u.LoggerService.Secure(authTag, fmt.Sprintf("failed to generate token: %s", err))
		return "", "", errors.New(internalServerError)
	}

	err = u.UserRepository.UpdateRefreshToken(ctx, foundUser.UserID, refreshToken)

	if err != nil {
		u.LoggerService.Secure(authTag, fmt.Sprintf("Failed to save to db: %s", err))
		return "", "", errors.New(internalServerError)
	}

	return accessToken, refreshToken, nil
}

func (u *AuthUseCase) Logout(ctx context.Context, userId string) error {
	if userId == "" {
		u.LoggerService.Error("userId is empty")
		return errors.New("userId is required")
	}

	err := u.UserRepository.UpdateRefreshToken(ctx, userId, "")
	if err != nil {
		u.LoggerService.Error(fmt.Sprintf("failed to update refresh token: %s", err))
		return errors.New(internalServerError)
	}

	return nil
}
