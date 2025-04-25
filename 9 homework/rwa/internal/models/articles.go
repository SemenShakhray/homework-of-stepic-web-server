package models

import "time"

type Articles struct {
	Article Article `json:"article"`
}

type Article struct {
	Author         Author    `json:"author,omitempty"`
	Body           string    `json:"body"`
	Title          string    `json:"title"`
	Description    string    `json:"description"`
	CreatedAt      time.Time `json:"created_at,omitempty"`
	UpdatedAt      time.Time `json:"updated_at,omitempty"`
	Favorited      bool      `json:"favorited,omitempty"`
	FavoritesCount int       `json:"favoritesCount,omitempty"`
	Slug           string    `json:"slug,omitempty"`
	TagList        []string  `json:"tagList,omitempty"`
}

type RequestNewArticle struct {
	Article struct {
		Body        string   `json:"body" binding:"required"`
		Description string   `json:"description" binding:"required"`
		Title       string   `json:"title" binding:"required"`
		TagList     []string `json:"tagList,omitempty"`
	}
}

type ArticleQueryParams struct {
	Tag       string `form:"tag"`
	Author    string `form:"author"`
	Favorited string `form:"favorited"`
	Limit     int    `form:"limit,default=20"`
	Offset    int    `form:"offset,default=0"`
}

type RequestAllArticleByFiltres struct {
	Articles      []Article `json:"articles"`
	ArticlesCount int       `json:"articlesCount"`
}
