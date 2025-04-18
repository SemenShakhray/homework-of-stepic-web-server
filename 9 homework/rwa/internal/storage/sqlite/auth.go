package sqlite

import (
	"database/sql"
	"fmt"
	"log"
	"rwa/internal/models"
	"strings"
)

func (s *Storage) Register(email, username, hashPassword string) error {
	query := "INSERT INTO users (email, username, password_hash) VALUES (?, ?, ?)"

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

func (s *Storage) GetUser(email string) (models.User, error) {
	var user models.User

	row := s.db.QueryRow("SELECT user_id, email, username, bio,  image, token, created_at, updated_at, password_hash FROM users WHERE email=?", email)
	err := row.Scan(&user.UserID,
		&user.Email,
		&user.Username,
		&user.Bio,
		&user.Image,
		&user.Token,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.PasswordHash,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Println("recors don't exists:", err)

			return models.User{}, fmt.Errorf("recors don't exists")
		}
		log.Println("failed to scan profile", err)

		return models.User{}, fmt.Errorf("failed to scan profile")
	}

	return user, nil
}

func (s *Storage) UpdateUser(user map[string]string, email string) error {
	var args []interface{}

	query := "UPDATE users SET "

	for field, value := range user {
		query += field + " = ?, "
		args = append(args, value)
	}

	query = strings.TrimSuffix(query, ", ")
	query += " WHERE email = ?"
	args = append(args, email)

	log.Println("query of update user", query)

	row, err := s.db.Exec(query, args...)
	if err != nil {
		log.Println("failed to update user:", err)

		return fmt.Errorf("failed to update user: %w", err)
	}

	n, err := row.RowsAffected()
	if err != nil {
		log.Println("failed to update user:", err)

		return fmt.Errorf("failed to retrieve rows affected: %w", err)
	}

	if n == 0 {
		log.Println("failed to update user: user not found")

		return fmt.Errorf("user not found")
	}

	return nil
}

// func (s *Storage) GetUserID(email string) (int, error) {
// 	var id int

// 	err := s.db.QueryRow("SELECT user_id FROM users WHERE email=?").Scan(&id)
// 	if err != nil {
// 		log.Println("failed to get userID:", err)

// 		return 0, fmt.Errorf("failed to get userID - %w", err)
// 	}

// 	return id, nil
// }
