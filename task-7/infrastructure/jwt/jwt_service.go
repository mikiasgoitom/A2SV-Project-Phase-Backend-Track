package jwtManager

import (
	"clean-architecture/usecase/contract"
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTManager struct {
	secretKey      string
	userRepository contract.IUserRepository
}

func NewJWTService(secretKey string, userRepo contract.IUserRepository) *JWTManager {
	return &JWTManager{
		secretKey:      secretKey,
		userRepository: userRepo,
	}
}

func (j *JWTManager) GenerateTokens(ctx context.Context, id string, username string, role string) (accessToken string, refreshToken string, err error) {
	// create claims
	accessTokenClaims := jwt.MapClaims{
		"id":       id,
		"username": username,
		"role":     role,
		"type":     "access",
		"exp":      time.Now().Add(time.Hour * 72),
	}
	// create access token
	accessToken, err = jwt.NewWithClaims(jwt.SigningMethodHS256, accessTokenClaims).SignedString(j.secretKey)
	if err != nil {
		return "", "", errors.New("failed to generate access token")
	}
	// create refresh token
	refreshTokenClaims := jwt.MapClaims{
		"id":       id,
		"username": username,
		"role":     role,
		"type":     "refresh",
		"exp":      time.Now().Add(time.Hour * 168).Unix(),
	}
	refreshToken, err = jwt.NewWithClaims(jwt.SigningMethodHS256, refreshTokenClaims).SignedString(j.secretKey)
	if err != nil {
		return "", "", errors.New("failed to generate refresh token")
	}
	if accessToken == "" || refreshToken == "" {
		return "", "", errors.New("failed to generate tokens")
	}
	j.userRepository.UpdateRefreshToken(ctx, id, refreshToken)
	return accessToken, refreshToken, nil
}

func (j *JWTManager) ValidateAccessToken(ctx context.Context, tokenString string) (id string, username string, role string, err error) {
	claims := jwt.MapClaims{}
	_, err = jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(j.secretKey), nil
	})
	if err != nil {
		return "", "", "", errors.New("invalid token")
	}
	return claims["id"].(string), claims["username"].(string), claims["role"].(string), nil
}

func (j *JWTManager) ValidateRefreshToken(ctx context.Context, tokenString string) (accessToken string, refreshToken string, err error) {
	claims := jwt.MapClaims{}
	_, err = jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return "", fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(j.secretKey), nil
	})
	idStr, _ := claims["id"].(string)
	usernameStr, _ := claims["username"].(string)
	roleStr, _ := claims["role"].(string)
	at, rt, err := j.GenerateTokens(ctx, idStr, usernameStr, roleStr)
	return at, rt, err
}
