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
