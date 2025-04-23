package sqlite

import (
	"fmt"
	"log"
	"rwa/internal/models"
)

func (s *Storage) CreateWithResponse(req models.RequestNewArticle, slug, email string) (models.Article, error) {
	user, err := s.GetUser(email)
	if err != nil {
		return models.Article{}, err
	}

	article, articleID, err := s.CreateArticle(req, slug, user.UserID)
	if err != nil {
		return models.Article{}, err
	}

	for _, tag := range req.Article.TagList {
		tagID, err := s.GetOrCreateTagID(tag)
		if err != nil {
			return models.Article{}, err
		}

		err = s.AddArticleTags(articleID, int(tagID))
		if err != nil {
			return models.Article{}, err
		}
	}

	tags, err := s.GetTagsbyArticle(articleID)
	if err != nil {
		return models.Article{}, err
	}

	article.Author.Bio = user.Bio
	article.Author.Username = user.Username
	article.Author.Image = user.Image
	article.TagList = tags

	log.Println("Article: ", article)

	return article, nil
}

func (s *Storage) CreateArticle(req models.RequestNewArticle, slug string, authorID int) (models.Article, int, error) {
	res, err := s.db.Exec("INSERT INTO articles (title, body, description, slug, author_id) VALUES (?, ?, ?, ?, ?)",
		req.Article.Title, req.Article.Body, req.Article.Description, slug, authorID)
	if err != nil {
		log.Println("failed to create article:", err)

		return models.Article{}, 0, fmt.Errorf("failed to crate article: %w", err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		log.Println("failed to get articleID:", err)

		return models.Article{}, 0, fmt.Errorf("failed to get articleID: %w", err)
	}

	article, err := s.GetArticle(int(id))
	if err != nil {
		return models.Article{}, 0, err
	}

	return article, int(id), nil
}

func (s *Storage) GetArticle(artID int) (models.Article, error) {
	var article models.Article
	err := s.db.QueryRow(`SELECT slug, 
 	title, 
 	description, 
 	body,
  	created_at,
   	updated_at,
    favorites_count
	FROM articles WHERE article_id = ?`, artID).
		Scan(&article.Slug,
			&article.Title,
			&article.Description,
			&article.Body,
			&article.CreatedAt,
			&article.UpdatedAt,
			&article.FavoritesCount,
		)
	if err != nil {
		log.Println("failed to scan parametrs article:", err)

		return models.Article{}, fmt.Errorf("failed to scan parametrs article: %w", err)
	}

	return article, nil
}
