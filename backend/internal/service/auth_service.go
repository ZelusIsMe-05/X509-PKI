package service

import (
	"errors"
	"x509-pki/internal/model"
	"x509-pki/internal/repository"
)

func Register(user model.User) error {
	if repository.Exists(user.Username) {
		return errors.New("username already exists")
	}

	if err := repository.Save(user); err != nil {
		return errors.New("failed to save user")
	}
	return nil
}

func Login(user model.User) error {
	password, exists := repository.GetPassword(user.Username)

	if !exists || password != user.Password {
		return errors.New("invalid username or password")
	}

	return nil
}

