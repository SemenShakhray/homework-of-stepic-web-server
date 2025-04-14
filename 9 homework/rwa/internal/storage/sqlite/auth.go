package sqlite

import (
	"fmt"
	"log"
	"rwa/internal/models"
)

func (s *Storage) Register(user models.Profile) (models.Profile, error) {
	query := "INSERT INTO users (email, username, password) VALUES (?, ?, ?)"

	_, err := s.db.Exec(query, user.User.Email, user.User.Username, user.User.Password)
	if err != nil {
		log.Println("failed created profile", err)

		return models.Profile{}, fmt.Errorf("failed created profile: %w", err)
	}

	var prof models.Profile

	rows := s.db.QueryRow("SELECT email, username, created_at, updated_at FROM users WHERE email=?", user.User.Email)
	err = rows.Scan(&prof.User.Email, &prof.User.Username, &prof.User.CreatedAt, &prof.User.UpdatedAt)
	if err != nil {
		log.Println("failed get profile after registration", err)

		return models.Profile{}, fmt.Errorf("failed get profile after registration: %w", err)
	}

	return prof, nil
}

func (s *Storage) GetPassword(profile models.Login) (string, error) {
	var pass string

	row := s.db.QueryRow("SELECT password FROM users WHERE email=?", profile.User.Email)
	if err := row.Scan(&pass); err != nil {
		log.Println("failed to get password", err)

		return "", fmt.Errorf("failed to get password")
	}

	return pass, nil
}

func (s *Storage) Login(profile models.Login) (models.Profile, error) {
	var prof models.Profile

	rows := s.db.QueryRow("SELECT email, username, created_at, updated_at FROM users WHERE email=?", profile.User.Email)
	err := rows.Scan(&prof.User.Email, &prof.User.Username, &prof.User.CreatedAt, &prof.User.UpdatedAt)
	if err != nil {
		log.Println("failed get profile after registration", err)

		return models.Profile{}, fmt.Errorf("failed get profile og login: %w", err)
	}

	return prof, nil
}
