package data

import (
	"context"
	"fmt"
	"log"
	"os"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var mongoClient *mongo.Client
var taskCollection *mongo.Collection
var userCollection *mongo.Collection

// ConnectToMongoDB initializes the MongoDB connection and prepares the task collection
func ConnectToMongoDB() error {
	mongoURL := os.Getenv("DATABASE_URL")
	fmt.Println(len(mongoURL))
	clientOptions := options.Client().ApplyURI(mongoURL)
	client, err := mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		return err
	}

	err = client.Ping(context.TODO(), nil)

	if err != nil {
		return err
	}

	mongoClient = client
	taskCollection = mongoClient.Database("taskdb").Collection("tasks")
	userCollection = mongoClient.Database("taskdb").Collection("users")
	log.Println("Connected to mongodb")

	return nil
}

func CloseMongoDBConnection() {
	if mongoClient != nil {
		if err := mongoClient.Disconnect(context.TODO()); err != nil {
			log.Println("Error disconnecting from MongoDB:", err)
		} else {
			log.Println("Disconnected from MongoDB")
		}
	}
}

func EnsureUsernameIndex(userCollection *mongo.Collection) error {
	indexModel := mongo.IndexModel{
		Keys:    bson.D{{Key: "username", Value: 1}},
		Options: options.Index().SetUnique(true),
	}

	indexName, err := userCollection.Indexes().CreateOne(context.TODO(), indexModel)

	if err != nil {
		log.Println("Error creating index on username:", err)
		return err
	}

	log.Printf("\nIndex on username created successfully: %s\n", indexName)
	return nil
}

func CreateUserIndex() {
	EnsureUsernameIndex(userCollection)
}
