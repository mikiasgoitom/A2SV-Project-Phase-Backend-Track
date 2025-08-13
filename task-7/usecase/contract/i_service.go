package contract

import "context"

type IPasswordService interface {
	HashPassword(password string) (string, error)
	ComparePassword(inputPassword, hashedPassword string) (bool, error)
}

type ITokenService interface {
	GenerateTokens(ctx context.Context, id string, username string, role string) (accessToken string, refreshToken string, err error)
	ValidateAccessToken(ctx context.Context, tokenString string) (id string, username string, role string, err error)
	ValidateRefreshToken(ctx context.Context, tokenString string) (accessToken string, refreshToken string, err error)
}
