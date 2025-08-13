package contract

import (
	"clean-architecture/domain"
	"context"
	"time"
)

type IUserRepository interface {
	CreateUser(ctx context.Context, user *domain.User) error
	GetUserByID(ctx context.Context, id string) (*domain.User, error)
	GetUserByUsername(ctx context.Context, username string) (*domain.User, error)
	UpdateUser(ctx context.Context, user *domain.User) error
	UpdateRefreshToken(ctx context.Context, id string, token string) error
	DeleteUser(ctx context.Context, id string) error
	CheckUsernameExist(ctx context.Context, username string) (bool, error)
}
type ITaskRepository interface {
	CreateTask(ctx context.Context, task *domain.Task) error
	GetTaskByID(ctx context.Context, id string) (*domain.Task, error)
	GetAllOverDueTasks(ctx context.Context, currentTime time.Time) ([]*domain.Task, error)
	// IsTaskOverDue(ctx context.Context, id string) (bool, error)
	UpdateTask(ctx context.Context, task *domain.Task) error
	DeleteTask(ctx context.Context, id string) error
}
