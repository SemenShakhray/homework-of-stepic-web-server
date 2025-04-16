package service

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"golang.org/x/crypto/bcrypt"

	"rwa/internal/lib"
	"rwa/internal/models"
)

func (s *Service) Register(req models.RequestNewUser) (models.User, error) {
	hashPass, err := bcrypt.GenerateFromPassword([]byte(req.User.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Println("failed to hash password")

		return models.User{}, fmt.Errorf("failed to hash password")
	}

	err = s.store.Register(req.User.Email, req.User.Username, string(hashPass))
	if err != nil {
		return models.User{}, err
	}

	user, err := s.store.GetUser(req.User.Email)
	if err != nil {
		return models.User{}, err
	}

	return user, nil
}

func (s *Service) Login(req models.RequestLogin, token string) (models.User, error) {
	err := s.store.AddToken(token, req.User.Email)
	if err != nil {
		return models.User{}, err
	}

	user, err := s.store.GetUser(req.User.Email)
	if err != nil {
		return models.User{}, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.User.Password))
	if err != nil {
		log.Println("invalid password", err)

		return models.User{}, fmt.Errorf("infalid password: %w", err)
	}

	return user, nil
}

func (s *Service) GetUser(email string) (models.User, error) {
	return s.store.GetUser(email)
}

func (s *Service) UpdateUser(user models.User, email string, exp time.Duration) (models.User, error) {
	jsonData, err := json.Marshal(user)
	if err != nil {
		log.Println("failed marsalling:", err)

		return models.User{}, fmt.Errorf("failed marsalling")
	}

	userMap := make(map[string]string)

	err = json.Unmarshal(jsonData, &userMap)
	if err != nil {
		log.Println("failed unmarsalling into map:", err)

		return models.User{}, fmt.Errorf("failed marsalling into map")
	}

	err = s.store.UpdateUser(userMap, email)
	if err != nil {
		return models.User{}, err
	}

	if e, ok := userMap["email"]; ok {
		jwtToken, err := lib.CreateJWT(e, exp)
		if err != nil {
			return models.User{}, err
		}

		err = s.store.AddToken(jwtToken, e)
		if err != nil {
			return models.User{}, err
		}

		email = e
	}

	userResp, err := s.store.GetUser(email)
	if err != nil {
		return models.User{}, err
	}

	return userResp, nil
}
