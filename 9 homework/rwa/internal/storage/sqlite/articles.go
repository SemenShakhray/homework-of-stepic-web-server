package sqlite

import (
	"bytes"
	"fmt"
	"log"
	"rwa/internal/models"
	"strings"
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

	// log.Println("Article: ", article)

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

func (s *Storage) GetAllArticlesByFiltres(params models.ArticleQueryParams) ([]models.Article, error) {
	var query bytes.Buffer
	var args []interface{}
	var conditions, joins []string

	query.WriteString(`SELECT a.slug, a.title, a.description, a.body, a.created_at, a.updated_at, a.favorites_count,
	u.username, u.bio, u.image
	GROUP_CONCAT(t.name) AS tags
	FROM articles a `)

	if params.Tag != "" {
		joins = append(joins, `
		JOIN article_tags at ON a.article_id=at.article_id
		JOIN tags t ON at.tag_id = t.tag_id
		`)
		conditions = append(conditions, "LOWER(t.name) LIKE LOWER(?)")
		args = append(args, `%`+params.Tag+"%")
	}

	if params.Author != "" {
		joins = append(joins, "JOIN users u ON a.author_id = u.user_id")
		conditions = append(conditions, "LOWER(u.username) LIKE LOWER(?)")
		args = append(args, "%"+params.Author+"%")
	}

	if params.Favorited != "" {
		joins = append(joins, `
		JOIN favorites f ON a.article_id=f.article_id
		JOIN users u_fav ON f.user_id = u_fav.user_id
		`)
		conditions = append(conditions, "LOWER(u.username) LIKE LOWER(?)")
		args = append(args, "%"+params.Favorited+"%")
	}

	for _, j := range joins {
		query.WriteString(j + " ")
	}

	if len(conditions) > 0 {
		query.WriteString("WHERE " + strings.Join(conditions, " AND ") + " ")
	}

	query.WriteString("LIMIT ? ")
	args = append(args, params.Limit)

	query.WriteString("OFFSET ? ")
	args = append(args, params.Offset)

	rows, err := s.db.Query(query.String(), args...)
	if err != nil {
		log.Println("failed query of search articles:", err, "query:", query.String())

		return nil, fmt.Errorf("failed query of search articles: %w", err)
	}
	defer rows.Close()

	var articles []models.Article
	for rows.Next() {
		var a models.Article
		err := rows.Scan(&a.Slug, &a.Title, &a.Description, &a.Body, &a.CreatedAt, &a.UpdatedAt, &a.FavoritesCount)
		if err != nil {
			log.Println("failed to scan article:", err)

			return nil, fmt.Errorf("failed to scan article: %w", err)
		}
		articles = append(articles, a)
	}

	return articles, nil
}
