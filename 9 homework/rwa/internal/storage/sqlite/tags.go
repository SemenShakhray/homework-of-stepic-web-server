package sqlite

import (
	"database/sql"
	"fmt"
	"log"
)

func (s *Storage) GetOrCreateTagID(tag string) (int64, error) {
	var tagID int64

	err := s.db.QueryRow("SELECT tag_id FROM tags WHERE name = ?", tag).Scan(&tagID)
	if err == sql.ErrNoRows {
		log.Println("tag don't exisits")

		res, err := s.db.Exec("INSERT INTO tags (name) VALUES (?)", tag)
		if err != nil {
			log.Println("failed to add tag:", err)

			return 0, fmt.Errorf("failed to add tag: %w", err)
		}

		id, err := res.LastInsertId()
		if err != nil {

		}

		return id, nil
	} else if err != nil {
		log.Println("error getting tag_id:", err)

		return 0, fmt.Errorf("error getting tag_id: %w", err)
	}

	return tagID, nil
}

func (s *Storage) AddArticleTags(artID, tagID int) error {
	_, err := s.db.Exec("INSERT INTO article_tags (article_id, tag_id) VALUES (?, ?)", artID, tagID)
	if err != nil {
		log.Println("failed to add artileID and tagID in the table:", err)

		return fmt.Errorf("failed to add artileID and tagID in the table: %w", err)
	}

	return nil
}

func (s *Storage) GetTagsbyArticle(artID int) ([]string, error) {
	var tags []string

	rows, err := s.db.Query(
		`SELECT t.name
	FROM tags t 
	JOIN article_tags at ON at.tag_id = t.tag_id
	WHERE at.article_id = ?`, artID)

	if err != nil {
		log.Println("failed to get tags by article:", err)

		return nil, fmt.Errorf("failed to get tags by article:%w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var tag string
		if err := rows.Scan(&tag); err != nil {
			log.Println("failed to scan name of tag:", err)

			return nil, fmt.Errorf("failed to scan name of tag: %w", err)
		}
		tags = append(tags, tag)
	}

	return tags, nil
}
