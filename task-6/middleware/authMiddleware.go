package middleware

import (
	"fmt"
	"net/http"
	"os"
	"slices"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware(requiredRole ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.IndentedJSON(http.StatusUnauthorized, gin.H{"message": "Authorization header is required"})
			c.Abort()
			return
		}

		authParts := strings.Split(authHeader, " ")

		if len(authParts) != 2 || strings.ToLower(authParts[0]) != "bearer" {
			c.IndentedJSON(http.StatusUnauthorized, gin.H{"message": "faulty authentication header"})
			c.Abort()
			return
		}

		// Extract the token from the header
		accessSecretKey := os.Getenv("ACCESS_SECRET_KEY")
		if accessSecretKey == "" {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "access secret key not set"})
			c.Abort()
			return
		}

		tokenString := authParts[1]
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(accessSecretKey), nil
		})

		if err != nil {
			c.IndentedJSON(http.StatusUnauthorized, gin.H{"message": "invalid token"})
			c.Abort()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok || !token.Valid {
			c.IndentedJSON(http.StatusUnauthorized, gin.H{"message": "invalid token claims"})
			c.Abort()
			return
		}
		userRole, _ := claims["user_type"].(string)
		if slices.Contains(requiredRole, userRole) {
			c.Set("user_id", claims["user_id"])
			c.Set("username", claims["username"])
			c.Set("user_type", userRole)
			c.Next()
			return
		}
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"message": "access denied"})
	}
}
