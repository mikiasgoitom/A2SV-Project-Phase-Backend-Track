package usecase

import (
	"clean-architecture/domain"
	"clean-architecture/usecase/contract"
	"context"
	"errors"
	"log"
	"time"

	"github.com/google/uuid"
)

type UserUseCase struct {
	Repository      contract.IUserRepository
	PasswordService contract.IPasswordService
}

func NewUserUseCase(repo contract.IUserRepository, pwdService contract.IPasswordService) *UserUseCase {
	return &UserUseCase{
		Repository:      repo,
		PasswordService: pwdService,
	}
}

func (u *UserUseCase) Register(ctx context.Context, userId, username, role, password string) (*domain.User, error) {
	exists, err := u.Repository.CheckUsernameExist(ctx, username)
	if err != nil {
		return &domain.User{}, errors.New("failed to check if username exists")
	}
	if !exists {
		return &domain.User{}, errors.New("username already exists")
	}

	hashPwd, err := u.PasswordService.HashPassword(password)
	if err != nil {
		return &domain.User{}, errors.New("internal server error")
	}

	newUser := domain.User{
		UserID:    uuid.NewString(),
		Username:  username,
		Password:  hashPwd,
		UserType:  domain.UserRole(role),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}

	if err := u.Repository.CreateUser(ctx, &newUser); err != nil {
		log.Println("failed to save user to repository")
		return &domain.User{}, errors.New("internal server error")
	}
	newUser.Password = ""
	return &newUser, nil

}

func (u *UserUseCase) UpdateUser(ctx context.Context, id, role, password string) (*domain.User, error) {
	foundUser, err := u.Repository.GetUserByID(ctx, id)
	if err != nil {
		return &domain.User{}, errors.New("user not found")
	}

	if role == "" {
		role = string(foundUser.UserType)
	}

	if password == "" {
		password = foundUser.Password
	} else {
		password, err = u.PasswordService.HashPassword(password)
		if err != nil {
			log.Println("error hashing new password")
			return &domain.User{}, errors.New("internal server error")
		}
	}

	updatedUser := domain.User{
		UserID:    id,
		Username:  foundUser.Username,
		Password:  password,
		UserType:  domain.UserRole(role),
		CreatedAt: foundUser.CreatedAt,
		UpdatedAt: time.Now().UTC(),
	}

	if err := u.Repository.UpdateUser(ctx, &updatedUser); err != nil {
		log.Println("error saving the updated user data")
		return &domain.User{}, errors.New("internal server error")
	}
	updatedUser.Password = ""
	return &updatedUser, nil
}

func (u *UserUseCase) DeleteUser(ctx context.Context, id string) error {
	foundUser, err := u.Repository.GetUserByID(ctx, id)
	if err != nil {
		log.Println("user not found by id")
		return errors.New("user not found")
	}

	if err := u.Repository.DeleteUser(ctx, foundUser.UserID); err != nil {
		log.Println("user not deleted due to repository error")
		return errors.New("failed to delete user")
	}

	return nil
}
