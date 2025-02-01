package service

import (
	"fmt"

	adapter "github.com/rismapa/go-banking/adapter/repository"
	config "github.com/rismapa/go-banking/config"

	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	LoginAccount(username, password string) (string, error)
}

type AuthAdapterDB struct {
	repo adapter.AccountRepository
}

func NewAuthService(repo adapter.AccountRepository) *AuthAdapterDB {
	return &AuthAdapterDB{repo: repo}
}

func (u *AuthAdapterDB) LoginAccount(username, password string) (string, error) {
	user, err := u.repo.GetAccountByUsername(username)
	if err != nil {
		return "", err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", fmt.Errorf("invalid password: %v", err)
	}

	token, err := config.GenerateJWT(user.ID, user.Username)
	if err != nil {
		return "", fmt.Errorf("could not generate token: %v", err)
	}

	return token, nil
}
