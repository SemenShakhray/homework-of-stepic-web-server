package service

import (
	"fmt"
	"log"

	"golang.org/x/crypto/bcrypt"

	"rwa/internal/models"
)

func (s *Service) Register(req models.RequestNewUser) (models.Users, error) {
	hashPass, err := bcrypt.GenerateFromPassword([]byte(req.User.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Println("failed to hash password")

		return models.Users{}, fmt.Errorf("failed to hash password")
	}

	err = s.store.Register(req.User.Email, req.User.Username, string(hashPass))
	if err != nil {
		return models.Users{}, err
	}

	user, err := s.store.GetUser(req.User.Email)
	if err != nil {
		return models.Users{}, err
	}

	return user, nil
}

func (s *Service) Login(req models.RequestLogin) (models.Users, error) {
	user, err := s.store.GetUser(req.User.Email)
	if err != nil {
		return models.Users{}, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.User.PasswordHash), []byte(req.User.Password))
	if err != nil {
		log.Println("invalid password", err)

		return models.Users{}, fmt.Errorf("infalid password: %w", err)
	}

	return user, nil
}

func (s *Service) GetUser(email string) (models.Users, error) {
	return s.store.GetUser(email)
}
