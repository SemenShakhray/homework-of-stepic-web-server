package models

import "time"

type Article struct {
	Author         User      `json:"author,omitempty"`
	Body           string    `json:"body"`
	Title          string    `json:"title"`
	Description    string    `json:"discription"`
	CreatedAt      time.Time `json:"created_at,omitempty"`
	UpdatedAt      time.Time `json:"updated_at,omitempty"`
	Favorited      bool      `json:"favorited,omitempty"`
	FavoritesCount int       `json:"favoritesCount,omitempty"`
	Slug           string    `json:"slug,omitempty"`
	TagList        []string  `json:"tagList,omitempty"`
}

type RequestNewArticle struct {
	Article struct {
		Body        string `json:"body" binding:"required"`
		Description string `json:"description" binding:"required"`
		Title       string `json:"title" binding:"required"`
		TagList     string `json:"tagList,omitempty"`
	}
}
