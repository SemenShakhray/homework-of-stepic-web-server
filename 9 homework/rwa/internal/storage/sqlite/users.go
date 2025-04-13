package sqlite

import (
	"fmt"
	"log"
	"rwa/internal/models"
)

func (s *Storage) Register(profile models.Profile) (models.Profile, error) {
	query := "INSERT INTO users (email, username, password) VALUES (?, ?, ?)"

	_, err := s.db.Exec(query, profile.Email, profile.Username, profile.Password)
	if err != nil {
		log.Println("failed created profile", err)

		return models.Profile{}, fmt.Errorf("failed created profile: %w", err)
	}

	var prof models.Profile

	rows := s.db.QueryRow("SELECT (email, username, create_at, update_at) FROM users WHERE email=?", profile.Email)
	err = rows.Scan(&prof.Email, &prof.Username, &prof.CreatedAt, &prof.UpdatedAt)
	if err != nil {
		log.Println("failed get profile after registration", err)

		return models.Profile{}, fmt.Errorf("failed get profile after registration: %w", err)
	}

	return prof, nil
}
