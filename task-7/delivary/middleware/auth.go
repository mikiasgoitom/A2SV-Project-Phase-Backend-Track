package middleware

import (
	"clean-architecture/usecase/contract"
	"strings"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware(jwtService contract.ITokenService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			ctx.JSON(401, gin.H{"error": "Authorization header is required"})
			ctx.Abort()
			return
		}
		authParts := strings.Split(authHeader, " ")
		if len(authParts) != 2 || strings.ToLower(authParts[0]) != "bearer" {
			ctx.JSON(401, gin.H{"error": "Invalid Authorization header format"})
			ctx.Abort()
			return
		}

		token := authParts[1]
		if token == "" {
			ctx.JSON(401, gin.H{"error": "Invalid token"})
			ctx.Abort()
			return
		}

		id, username, role, err := jwtService.ValidateAccessToken(ctx.Request.Context(), token)
		if err != nil {
			ctx.JSON(401, gin.H{"error": "Invalid token"})
			ctx.Abort()
			return
		}

		// Set user information in the context
		ctx.Set("userID", id)
		ctx.Set("username", username)
		ctx.Set("role", role)

		ctx.Next()
	}
}
