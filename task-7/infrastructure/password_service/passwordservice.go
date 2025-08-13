package passwordservice

import (
	"clean-architecture/usecase/contract"

	"golang.org/x/crypto/bcrypt"
)

type PasswordService struct {
}

func NewPasswordService() contract.IPasswordService {
	return &PasswordService{}
}

func (h *PasswordService) HashPassword(password string) (string, error) {
	hashedBytePassword, err := bcrypt.GenerateFromPassword([]byte(password), 5)
	if err != nil {
		return "", err
	}
	return string(hashedBytePassword), nil
}

func (h *PasswordService) ComparePassword(inputPassword, hashedPassword string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(inputPassword))
	if err != nil {
		return false, err
	}
	return true, nil
}
