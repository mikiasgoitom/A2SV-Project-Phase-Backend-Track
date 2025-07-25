package data

import (
	"context"
	"errors"
	"log"
	"os"
	"task-6/models"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

// customClaim defines the JWT claims structure for user tokens.
type customClaim struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	UserType string `json:"user_type"`
	jwt.RegisteredClaims
}

func GenerateTokens(userId, username, userType string) (string, string, error) {
	atClaim := customClaim{
		UserID:   userId,
		Username: username,
		UserType: userType,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   userId,
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
		},
	}
	accessSecretKey := os.Getenv("ACCESS_SECRET_KEY")
	if accessSecretKey == "" {
		log.Println("Error getting access secret key")
		return "", "", errors.New("failed to get access secret key")
	}
	refreshSecretKey := os.Getenv("REFRESH_SECRET_KEY")
	if refreshSecretKey == "" {
		log.Println("Error getting refresh secret key")
		return "", "", errors.New("failed to get refresh secret key")
	}
	// Generate access token
	accessToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaim).SignedString([]byte(accessSecretKey))
	if err != nil {
		log.Println("Error generating access token:", err)
		return "", "", errors.New("failed to generate access token")
	}
	// Generate refresh token
	rtClaim := customClaim{
		UserID:   userId,
		Username: username,
		UserType: userType,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   userId,
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(168 * time.Hour)),
		},
	}
	refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, rtClaim).SignedString([]byte(refreshSecretKey))
	if err != nil {
		log.Println("Error generating refresh token:", err)
		return "", "", errors.New("failed to generate refresh token")
	}
	// Return the generated tokens
	log.Println("Tokens generated successfully")
	return accessToken, refreshToken, nil
}

func GetAllUsers() []models.User {
	var users []models.User
	cursor, err := userCollection.Find(context.TODO(), bson.D{})
	if err != nil {
		log.Println("Error fetching users:", err)
		return []models.User{}
	}
	defer cursor.Close(context.TODO())

	for cursor.Next(context.TODO()) {
		var user models.User
		if err := cursor.Decode(&user); err != nil {
			log.Println("Error decoding user:", err)
			continue
		}
		users = append(users, user)
	}

	return users
}

func GetUserByID(id string) models.User {
	filter := bson.M{"user_id": id}
	var user models.User
	err := userCollection.FindOne(context.TODO(), filter).Decode(&user)

	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			log.Println("User not found:", id)
			return models.User{}
		}
		log.Println("Error decodeing user:", err)
		return models.User{}
	}

	return user
}

func UpdateUser(id string, updatedUser models.User) models.User {
	updateFields := bson.D{}

	if updatedUser.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(updatedUser.Password), 5)
		if err != nil {
			log.Println("Error hashing password:", err)
			return models.User{}
		}
		updatedUser.Password = string(hashedPassword)
		updateFields = append(updateFields, bson.E{Key: "password", Value: updatedUser.Password})
	}
	if updatedUser.UserType != "" {
		updateFields = append(updateFields, bson.E{Key: "user_type", Value: updatedUser.UserType})
	}

	updateFields = append(updateFields, bson.E{Key: "update"})

	filter := bson.M{"user_id": id}
	update := bson.D{
		{Key: "$set", Value: updateFields},
	}

	_, err := userCollection.UpdateOne(context.TODO(), filter, update)

	if err != nil {
		log.Println("Error updating user:", err)
		return models.User{}
	}
	// Fetch the updated user to return
	var user models.User
	// reused filter to find the updated user
	err = userCollection.FindOne(context.TODO(), filter).Decode(&user)
	if err != nil {
		log.Println("Error fetching updated user:", err)
		return models.User{}
	}
	return user
}
func DeleteUser(id string) error {
	filter := bson.M{"user_id": id}
	result, err := userCollection.DeleteOne(context.TODO(), filter)
	if err != nil {
		log.Println("Error deleting user:", err)
		return err
	}

	if result.DeletedCount == 0 {
		log.Println("No user found with the given ID")
		return errors.New("no user found with the given ID")
	}

	return nil
}

func LoginUser(loginData models.Crediential) (string, error) {
	filter := bson.M{"username": loginData.Username}
	var foundUser models.User
	err := userCollection.FindOne(context.TODO(), filter).Decode(&foundUser)

	if err != nil {
		log.Println("\nfailed to fetch user data")
		return "", errors.New("failed to fetch user data: " + err.Error())
	}

	err = bcrypt.CompareHashAndPassword([]byte(foundUser.Password), []byte(loginData.Password))
	if err != nil {
		log.Println("Invalid username or password:", err)
		return "", errors.New("invalid username or password")
	}
	accessToken, refreshToken, err := GenerateTokens(foundUser.UserID, foundUser.Username, foundUser.UserType)

	if err != nil {
		log.Println("Error generating tokens:", err)
		return "", errors.New("failed to generate tokens: " + err.Error())
	}
	// Update the user's refresh token in the database
	update := bson.D{
		{Key: "$set", Value: bson.D{{Key: "refresh_token", Value: refreshToken}}},
	}

	updatedUser, err := userCollection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		log.Println("Error updating refresh token:", err)
		return "", errors.New("failed to update refresh token: " + err.Error())
	}
	if updatedUser.MatchedCount == 0 {
		log.Println("No user found with the given username")
		return "", errors.New("no user found with the given username")
	}

	return accessToken, nil
}

func LogoutUser(id string) error {
	filter := bson.M{"user_id": id}
	update := bson.D{
		{Key: "$set", Value: bson.D{{Key: "refresh_token", Value: ""}}},
	}

	_, err := userCollection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		log.Println("Error logging out user:", err)
		return err
	}

	log.Println("User logged out successfully")
	return nil
}

func CreateUser(newUser models.User) (string, error) {
	filter := bson.M{"username": newUser.Username}
	var existingUser models.User
	err := userCollection.FindOne(context.TODO(), filter).Decode(&existingUser)

	// check if the user already exists
	if err == nil {
		return "", errors.New("username already exists")
	}
	// if the err is not nil and it doesnt have no document found error, it means there was an error othanr than not finding the document
	if err != mongo.ErrNoDocuments {
		return "", errors.New("Internal server error: " + err.Error())
	}

	accessToken, refreshToken, err := GenerateTokens(newUser.UserID, newUser.Username, newUser.UserType)

	if err != nil {
		return "", errors.New("failed to generate tokens: " + err.Error())
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newUser.Password), 5)
	if err != nil {
		return "", errors.New("failed to hash password: " + err.Error())
	}
	newUser.UserID = uuid.New().String()
	newUser.Password = string(hashedPassword)
	newUser.RefreshToken = refreshToken
	newUser.CreatedAt = time.Now()
	newUser.UpdatedAt = time.Now()

	// Insert the user into the collection
	_, err = userCollection.InsertOne(context.TODO(), newUser)
	if err != nil {
		return "", errors.New("failed to create user: " + err.Error())
	}

	return accessToken, nil
}
