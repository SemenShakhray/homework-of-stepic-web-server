package service

import (
	"fmt"
	"regexp"
	"rwa/internal/models"
	"strings"
	"time"
)

func (s *Service) CreateArticle(req models.RequestNewArticle, email string) (models.Article, error) {
	slug := generateSlug(req.Article.Title)

	article, err := s.store.CreateWithResponse(req, slug, email)
	if err != nil {
		return models.Article{}, err
	}

	return article, nil
}

func generateSlug(title string) string {
	slug := strings.ToLower(title)
	re := regexp.MustCompile(`[^\w\s-]`)
	slug = re.ReplaceAllString(slug, "")
	slug = strings.ReplaceAll(slug, " ", "-")
	slug = strings.ReplaceAll(slug, "_", "-")
	slug = regexp.MustCompile(`-{2,}`).ReplaceAllString(slug, "-")
	slug = strings.Trim(slug, "-")
	slug += "-" + fmt.Sprintf("%d", time.Now().Unix())

	return slug
}

func (s *Service) GetAllArticlesByFiltres(params models.ArticleQueryParams) (models.RequestAllArticleByFiltres, error) {
	articles, err := s.store.GetAllArticlesByFiltres(params)
	if err != nil {
		return models.RequestAllArticleByFiltres{}, err
	}

	count := len(articles)

	return models.RequestAllArticleByFiltres{
		Articles:      articles,
		ArticlesCount: count,
	}, nil
}
