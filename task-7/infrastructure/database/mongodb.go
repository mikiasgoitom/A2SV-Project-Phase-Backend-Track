package database

import (
	"clean-architecture/usecase/contract"
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func InitMongoDB(ctx context.Context, url, dbName string, logger contract.ILogger) (*mongo.Client, *mongo.Database, error) {
	clientOption := options.Client().ApplyURI(url)
	client, err := mongo.Connect(ctx, clientOption)

	if err != nil {
		logger.Error(fmt.Sprintf("InitMongoDB: failed to connect to MongoDB: (error) %v", err))
		return nil, nil, fmt.Errorf("failed to connect to MongoDB: %w", err)
	}
	if err = client.Ping(ctx, nil); err != nil {
		logger.Error(fmt.Sprintf("InitMongoDB: failed to ping to MongoDB: (error) %v", err))
		return nil, nil, fmt.Errorf("failed to ping to MongoDB: %w", err)
	}
	db := client.Database("task-seven")

	return client, db, nil
}

func CloseMongoDB(ctx context.Context, client *mongo.Client) error {
	if err := client.Disconnect(ctx); err != nil {
		return fmt.Errorf("failed to disconnect MongoDB client: %w", err)
	}
	return nil
}
