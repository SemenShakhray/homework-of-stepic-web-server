package sqlite

import (
	"fmt"
	"log"
	"rwa/internal/models"
)

func (s *Storage) CreateWithResponse(req models.RequestNewArticle, slug, email string) (models.Article, error) {
	return models.Article{}, nil
}

func (s *Storage) Create(req models.RequestNewArticle, slug, email string) (int, error) {
	// user, err := s.GetUser(email)
	// if err != nil {
	// 	return models.Article{}, err
	// }
	res, err := s.db.Exec("INSERT INTO articles (title, body, description) VALUES (?, ?, ?)",
		req.Article.Title, req.Article.Body, req.Article.Description)
	if err != nil {
		log.Println("failed to crate article:", err)

		return 0, fmt.Errorf("failed to crate article: %w", err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		log.Println("failed to get articleID:", err)

		return 0, fmt.Errorf("failed to get articleID: %w", err)
	}

	return int(id), nil
}
