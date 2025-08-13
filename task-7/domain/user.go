package domain

import (
	"time"
)

type UserRole string

const (
	Admin       UserRole = "admin"
	RegularUser UserRole = "user"
)

type User struct {
	UserID       string
	Username     string
	Password     string
	RefreshToken string
	UserType     UserRole
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func IsValidUserRole(role string) bool {
	switch role {
	case string(Admin), string(RegularUser):
		return true
	}
	return false
}
