package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID       uint   `json:"id"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"-"`
}

var users = make(map[string]*User)
var countUsers = 0

// secrey key for JWT signing
var jwtSecret = []byte("very secret key")

func main() {
	router := gin.Default()
	router.GET("/", func(ctx *gin.Context) {
		ctx.IndentedJSON(http.StatusOK, gin.H{
			"message": "Welcome to Go Authentication and Authorization",
		})
	})
	router.POST("/register", func(ctx *gin.Context) {
		var user User
		if err := ctx.ShouldBindJSON(&user); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": "Invalid request payload",
			})
			return
		}

		if _, exists := users[user.Email]; exists {
			ctx.JSON(http.StatusConflict, gin.H{
				"message": "User already exists"})
			return
		}
		fmt.Println("Registering user:", user)
		if user.Password == "" {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": "Invalid password. Password cannot be empty.",
			})
			return
		}

		user.ID = uint(countUsers) + 1
		countUsers += 1

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), 5)

		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"message": "Something went wrong. Please try again.",
			})
			log.Println("\nError hashing password:", err)
			return
		}
		user.Password = string(hashedPassword)

		users[user.Email] = &user
		ctx.JSON(http.StatusOK, gin.H{
			"message": "User registered successfully",
		})
	})

	router.POST("/login", func(ctx *gin.Context) {
		var user User

		if err := ctx.ShouldBindJSON(&user); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": "Invalid request payload",
			})
			return
		}

		foundUser, exists := users[user.Email]

		if err := bcrypt.CompareHashAndPassword([]byte(foundUser.Password), []byte(user.Password)); err != nil || !exists {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"message": "invalid email or password",
			})
			return
		}

		token, err := generateJWT(foundUser, jwtSecret)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"message": "something went wrong. Please try again.",
			})
			log.Println("\nError generating JWT token:", err)
			return
		}
		ctx.JSON(http.StatusOK, gin.H{
			"message": "Login successful",
			"token":   token,
		})
	})

	router.GET("/secure", authMiddleware(), func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "This is a secure route"})
	})
	if err := router.Run("localhost:8080"); err != nil {
		panic("failed to start the server at port 8080")
	}
}

func generateJWT(user *User, secret []byte) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"email":   user.Email,
	})

	jwtToken, err := token.SignedString(secret)

	if err != nil {
		log.Println("Error generating JWT token:", err)
		return "", err
	}

	return jwtToken, nil
}

func authMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"message": "Login required",
			})
			log.Println("\nAuthorization header is missing")
			ctx.Abort()
			return
		}

		authParts := strings.Split(authHeader, " ")
		if len(authParts) != 2 || strings.ToLower(authParts[0]) != "bearer" {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"message": "Login required",
			})
			log.Println("\nAuthorization header is malformed")
			ctx.Abort()
			return
		}

		tokenString := authParts[1]

		_, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return jwtSecret, nil
		})

		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"message": "Invalid token",
			})
			log.Println("\nError parsing JWT token:", err)
			ctx.Abort()
			return
		}

		ctx.Next()
	}
}
