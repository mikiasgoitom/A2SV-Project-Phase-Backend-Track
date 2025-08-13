package repository

import (
	"clean-architecture/domain"
	"clean-architecture/usecase/contract"
	"context"
	"errors"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserRepository struct {
	Collection *mongo.Collection
	Logger     contract.ILogger
}

// to implement the IUserRepository interface at compile time
var _ contract.IUserRepository = (*UserRepository)(nil)

func NewUserRepository(colln *mongo.Collection, lg contract.ILogger) *UserRepository {
	if colln == nil {
		panic("mongo collection cannot be nil for UserRepository")
	}
	if lg == nil {
		panic("logger cannot be nil for UserRepository")
	}
	return &UserRepository{
		Collection: colln,
		Logger:     lg,
	}
}

func (r *UserRepository) CreateUser(ctx context.Context, user *domain.User) error {
	_, err := r.Collection.InsertOne(ctx, user)
	if err != nil {
		r.Logger.Error(fmt.Sprintf("failed to store user: %v", err))
		return err
	}
	r.Logger.Info("successfully stored user")
	return nil
}
func (r *UserRepository) GetUserByID(ctx context.Context, id string) (*domain.User, error) {
	filter := bson.M{"user_id": id}
	var user *domain.User
	err := r.Collection.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		r.Logger.Debug(fmt.Sprintf("user not found: %v", err))
		return &domain.User{}, err
	}
	r.Logger.Info("user found")
	return user, nil
}

func (r *UserRepository) GetUserByUsername(ctx context.Context, username string) (*domain.User, error) {
	filter := bson.M{"username": username}
	var user *domain.User
	err := r.Collection.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		r.Logger.Debug(fmt.Sprintf("no user with username: %s", username))
		return &domain.User{}, err
	}
	r.Logger.Info(fmt.Sprintf("found user with username: %s", username))
	return user, nil
}
func (r *UserRepository) UpdateUser(ctx context.Context, user *domain.User) error {
	filter := bson.M{"user_id": user.UserID}
	update := bson.D{
		{Key: "$set", Value: user},
	}
	updateCount, err := r.Collection.UpdateOne(ctx, filter, update)
	if err != nil {
		r.Logger.Debug(fmt.Sprintf("failed to store user update: %v", err))
		return err
	}
	if updateCount.MatchedCount == 0 {
		r.Logger.Debug(fmt.Sprintf("failed to find users with user id: %v", user.UserID))
		return errors.New("user not found")
	}

	r.Logger.Info("user data successfully updated")
	return nil
}
func (r *UserRepository) UpdateRefreshToken(ctx context.Context, id string, token string) error {
	filter := bson.M{"user_id": id}
	var user *domain.User
	err := r.Collection.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		r.Logger.Debug(fmt.Sprintf("no user with id: %s", id))
		return err
	}
	user.RefreshToken = token
	user.UpdatedAt = time.Now().UTC()
	update := bson.D{
		{Key: "$set", Value: user},
	}
	updateCount, err := r.Collection.UpdateOne(ctx, filter, update)
	if err != nil {
		r.Logger.Debug(fmt.Sprintf("failed to store updated refresh token: %v", err))
		return err
	}
	if updateCount.MatchedCount == 0 {
		r.Logger.Debug(fmt.Sprintf("failed to find users with user id: %v", user.UserID))
		return errors.New("user not found")
	}

	r.Logger.Info("successfully updated refresh token")
	return nil
}
func (r *UserRepository) DeleteUser(ctx context.Context, id string) error {
	filter := bson.M{"user_id": id}
	count, err := r.Collection.DeleteOne(ctx, filter)
	if err != nil {
		r.Logger.Error(fmt.Sprintf("failed to delete user: %v", err))
		return err
	}
	if count.DeletedCount == 0 {
		r.Logger.Debug(fmt.Sprintf("no user with id: %s", id))
		return err
	}
	return nil
}

func (r *UserRepository) CheckUsernameExist(ctx context.Context, username string) (bool, error) {
	filter := bson.M{"username": username}
	count, err := r.Collection.CountDocuments(ctx, filter)
	if err != nil {
		r.Logger.Error(fmt.Sprintf("failed to fetch count of username: %s", username))
		return false, err
	}

	if count == 0 {
		r.Logger.Debug(fmt.Sprintf("no user with username: %s", username))
		return false, nil
	}
	r.Logger.Debug(fmt.Sprintf("(%v) user with username (%s) found", count, username))
	return true, nil
}
