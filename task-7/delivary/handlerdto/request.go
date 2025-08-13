package handlerdto

import (
	"clean-architecture/domain"
	"time"
)

type UserDto struct {
	UserID    string    `json:"user_id"`
	Username  string    `json:"username"`
	Password  string    `json:"password"`
	UserType  string    `json:"user_type"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func UserDtoToEntity(dto UserDto) *domain.User {
	return &domain.User{
		UserID:       dto.UserID,
		Username:     dto.Username,
		Password:     dto.Password,
		RefreshToken: "",
		UserType:     domain.UserRole(dto.UserType),
		CreatedAt:    dto.CreatedAt,
		UpdatedAt:    dto.UpdatedAt,
	}
}
func EntityToUserDto(entity domain.User) *UserDto {
	return &UserDto{
		UserID:    entity.UserID,
		Username:  entity.Username,
		Password:  entity.Password,
		UserType:  string(entity.UserType),
		CreatedAt: entity.CreatedAt,
		UpdatedAt: entity.UpdatedAt,
	}
}

type TaskDto struct {
	TaskID      string    `json:"task_id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Status      string    `json:"status"`
	DueDate     time.Time `json:"due_date"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func TaskDtoToEntity(dto TaskDto) *domain.Task {
	return &domain.Task{
		TaskID:      dto.TaskID,
		Title:       dto.Title,
		Description: dto.Description,
		Status:      domain.TaskStatus(dto.Status),
		DueDate:     dto.DueDate,
		CreatedAt:   dto.CreatedAt,
		UpdatedAt:   dto.UpdatedAt,
	}
}
func EntityToTaskDto(entity domain.Task) *TaskDto {
	return &TaskDto{
		TaskID:      entity.TaskID,
		Title:       entity.Title,
		Description: entity.Description,
		Status:      string(entity.Status),
		DueDate:     entity.DueDate,
		CreatedAt:   entity.CreatedAt,
		UpdatedAt:   entity.UpdatedAt,
	}
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
type LogoutRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}
