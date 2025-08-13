package contract

import (
	"clean-architecture/domain"
	"context"
	"time"
)

type IUserUseCase interface {
	Register(ctx context.Context, userId, username, role, password string) (*domain.User, error)
	UpdateUser(ctx context.Context, id, role, password string) (*domain.User, error)
	DeleteUser(ctx context.Context, id string) error
}

type ITaskUseCase interface {
	CreateNewTask(ctx context.Context, title, description string, status domain.TaskStatus, dueDate time.Time) (*domain.Task, error)
	UpdateTask(ctx context.Context, id, title, description string, status domain.TaskStatus, dueDate time.Time) (*domain.Task, error)
	DeleteTask(ctx context.Context, id string) error
	CompleteTask(ctx context.Context, id string) error
	OverDueTasks(ctx context.Context) ([]*domain.Task, error)
}

type IAuthUseCase interface {
	Login(ctx context.Context, username, password string) (accessToken string, refreshToken string, err error)
	Logout(ctx context.Context, userId string) error
}
