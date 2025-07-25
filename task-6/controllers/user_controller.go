package controllers

import (
	"log"
	"net/http"
	"task-6/data"
	"task-6/models"

	"github.com/gin-gonic/gin"
)

func GetAllUsers(c *gin.Context) {
	users := data.GetAllUsers()
	if users != nil {
		c.IndentedJSON(http.StatusOK, users)
		return
	}
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "no users found"})
}

func GetUserByID(c *gin.Context) {
	id := c.Param("user_id")

	user := data.GetUserByID(id)
	if user == (models.User{}) {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "user not found"})
		return
	}

	c.IndentedJSON(http.StatusOK, user)
}

func UpdateUser(c *gin.Context) {
	id := c.Param("user_id")
	var updatedUser models.User

	if err := c.ShouldBindJSON(&updatedUser); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Invalid request body"})
		return
	}

	user := data.UpdateUser(id, updatedUser)
	if user == (models.User{}) {
		c.IndentedJSON(http.StatusOK, user)
		return
	}

	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "user not found"})
}

func DeleteUser(c *gin.Context) {
	id := c.Param("user_id")
	err := data.DeleteUser(id)
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}
	c.IndentedJSON(http.StatusOK, gin.H{"message": "user deleted"})
}

func CreateUser(c *gin.Context) {
	var newUser models.User
	if err := c.ShouldBindJSON(&newUser); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	accessToken, err := data.CreateUser(newUser)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "failed to create user: " + err.Error()})
		return
	}
	c.IndentedJSON(http.StatusCreated, []interface{}{newUser, accessToken})
}

func LoginUser(c *gin.Context) {
	var loginData models.Crediential
	if err := c.ShouldBindJSON(&loginData); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
	}
	accessToken, err := data.LoginUser(loginData)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "failed to login user"})
		log.Println("Error logging in user:", err)
		return
	}

	if accessToken == "" {
		c.IndentedJSON(http.StatusUnauthorized, gin.H{"error": "invalid username or password"})
		return
	}
	c.IndentedJSON(http.StatusOK, gin.H{"access_token": accessToken})
}

func LogoutUser(c *gin.Context) {
	id := c.Param("user_id")

	err := data.LogoutUser(id)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "failed to logout user"})
		log.Println("Error logging out user:", err)
		return
	}
	c.IndentedJSON(http.StatusOK, gin.H{"message": "user logged out successfully"})
}
