package sqlite

import (
	"fmt"
	"log"

	"taskbot/internal/models"
)

func (s *Storage) GetAllTasks() ([]models.Task, error) {
	var tasks []models.Task

	rows, err := s.db.Query(`SELECT task_id, description, owner_id, assign_name 
	FROM tasks`)
	if err != nil {
		log.Println("failed query of get tasks:", err)

		return nil, fmt.Errorf("failed query of get tasks: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var task models.Task

		err := rows.Scan(&task.TaskID, &task.Description, &task.OwnerID, &task.AssignName)
		if err != nil {
			log.Println("failed to scan task:", err)

			return nil, fmt.Errorf("failed to scan task: %w", err)
		}

		tasks = append(tasks, task)
	}

	return tasks, nil
}

func (s *Storage) NewTask(description, assignName string, ownerID int) (map[int]string, error) {
	query := "INSERT INTO tasks (description, owner_id, assign_name) VALUES (?, ?, ?)"

	res, err := s.db.Exec(query, description, ownerID, "@"+assignName)
	if err != nil {
		log.Println("error adding a new task", err)

		return nil, fmt.Errorf("failed created task: %w", err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		log.Println("couldn't get the ID of the last record", err)

		return nil, fmt.Errorf("failed created user: %w", err)
	}

	if id == 0 {
		log.Println("ID last record is nil", err)

		return nil, fmt.Errorf("failed created user: %w", err)
	}

	answer := make(map[int]string)
	answer[ownerID] = fmt.Sprintf(`Задача "%s" создана, id=%d`, description, id)

	return answer, nil
}

func (s *Storage) Assing(taskID, assignedID int) error {
	row, err := s.db.Exec("UPDATE tasks SET assigned_id = ? WHERE task_id = ?", assignedID, taskID)
	if err != nil {
		log.Println("failed to assigned task:", err)

		return fmt.Errorf("failed to assigned task: %w", err)
	}

	id, err := row.RowsAffected()
	if id == 0 || err != nil {
		return fmt.Errorf("failed to assigned task: %w", err)
	}

	return nil
}

func (s *Storage) Unassing(taskID int) error {
	row, err := s.db.Exec("UPDATE tasks SET assigned_id = NULL WHERE task_id = ?", taskID)
	if err != nil {
		log.Println("failed to assigned task:", err)

		return fmt.Errorf("failed to unassigned task: %w", err)
	}

	id, err := row.RowsAffected()
	if id == 0 || err != nil {
		return fmt.Errorf("failed to unassigned task: %w", err)
	}

	return nil
}

func (s *Storage) Resolve(taskID int) error {
	row, err := s.db.Exec("DELETE FROM tasks WHERE task_id = ?", taskID)
	if err != nil {
		log.Println("failed to resolve task:", err)

		return fmt.Errorf("failed to resolve task: %w", err)
	}

	id, err := row.RowsAffected()
	if id == 0 || err != nil {
		return fmt.Errorf("failed to resolve task: %w", err)
	}

	return nil
}

func (s *Storage) GetMyTasks(myID int) ([]models.Task, error) {
	var tasks []models.Task

	rows, err := s.db.Query(`SELECT task_id, description, owner_id, assigned_id 
	FROM tasks WHERE assigned_id=?`, myID)
	if err != nil {
		log.Println("failed query of get my tasks:", err)

		return nil, fmt.Errorf("failed query of get my tasks: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var task models.Task

		err := rows.Scan(&task.TaskID, &task.Description, &task.OwnerID, &task.AssignName)
		if err != nil {
			log.Println("failed to scan task:", err)

			return nil, fmt.Errorf("failed to scan task: %w", err)
		}

		tasks = append(tasks, task)
	}

	return tasks, nil
}

func (s *Storage) GetOwnTasks(myID int) ([]models.Task, error) {
	var tasks []models.Task

	rows, err := s.db.Query(`SELECT task_id, description, owner_id, assigned_id 
	FROM tasks WHERE owner_id=?`, myID)
	if err != nil {
		log.Println("failed query of get own tasks:", err)

		return nil, fmt.Errorf("failed query of get own tasks: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var task models.Task

		err := rows.Scan(&task.TaskID, &task.Description, &task.OwnerID, &task.AssignName)
		if err != nil {
			log.Println("failed to scan task:", err)

			return nil, fmt.Errorf("failed to scan task: %w", err)
		}

		tasks = append(tasks, task)
	}

	return tasks, nil
}
