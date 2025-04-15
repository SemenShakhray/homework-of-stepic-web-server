package sqlite

import (
	"fmt"
	"log"
	"rwa/internal/models"
)

func (s *Storage) Register(email, username, hashPassword string) error {
	query := "INSERT INTO users (email, username, password) VALUES (?, ?, ?)"

	res, err := s.db.Exec(query, email, username, hashPassword)
	if err != nil {
		log.Println("error adding a new user", err)

		return fmt.Errorf("failed created user: %w", err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		log.Println("couldn't get the ID of the last record", err)

		return fmt.Errorf("failed created user: %w", err)
	}

	if id == 0 {
		log.Println("ID last record is nil", err)

		return fmt.Errorf("failed created user: %w", err)
	}

	return nil
}

func (s *Storage) GetPassword(profile models.RequestLogin) (string, error) {
	var pass string

	row := s.db.QueryRow("SELECT password FROM users WHERE email=?", profile.User.Email)
	if err := row.Scan(&pass); err != nil {
		log.Println("failed to get password", err)

		return "", fmt.Errorf("failed to get password")
	}

	return pass, nil
}

func (s *Storage) Login(req models.RequestLogin) (models.Users, error) {
	var user models.Users

	rows := s.db.QueryRow("SELECT email, username, password_hash, created_at, updated_at FROM users WHERE email=?", req.User.Email)
	err := rows.Scan(&user.User.Email, &user.User.Username, &user.User.PasswordHash, &user.User.CreatedAt, &user.User.UpdatedAt)
	if err != nil {
		log.Println("failed get profile after registration", err)

		return models.Users{}, fmt.Errorf("failed get profile og login: %w", err)
	}

	return user, nil
}

func (s *Storage) AddToken(token, email string) error {
	res, err := s.db.Exec("UPDATE users SET token=? WHERE email=?", token, email)
	if err != nil {
		log.Println("failed to added token:", err)

		return fmt.Errorf("failed to added token")
	}

	n, err := res.RowsAffected()
	if err != nil {
		log.Println("failed to added token:", err)

		return fmt.Errorf("failed to retrieve rows affected: %w", err)
	}

	if n == 0 {
		log.Println("failed to added token: user not found")

		return fmt.Errorf("user not found")
	}

	return nil
}

func (s *Storage) GetUser(email string) (models.Users, error) {
	var user models.Users

	row := s.db.QueryRow("SELECT email, username, bio,  image,token, created_at, updated_at FROM users WHERE email=?", email)
	err := row.Scan(&user.User.Email,
		&user.User.Username,
		&user.User.CreatedAt,
		&user.User.UpdatedAt,
		&user.User.Token,
	)
	if err != nil {
		log.Println("failed to scan profile", err)

		return models.Users{}, fmt.Errorf("failed to scan profile")
	}

	return user, nil
}
