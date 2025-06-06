package service

import "fmt"

func (s *Service) GetAllTasks(userID int, userName string) (map[int]string, error) {
	answer := make(map[int]string)

	tasks, err := s.store.GetAllTasks()
	if err != nil {
		return nil, err
	}

	if tasks == nil {
		answer[userID] = "Нет задач"

		return answer, nil
	}

	for _, task := range tasks {
		if _, ok := answer[task.OwnerID]; ok {
			answer[task.OwnerID] += fmt.Sprintf(`
			%s by %s, /assign_%d`, task.Description, task.AssignName, task.TaskID)
		} else {
			answer[task.OwnerID] = fmt.Sprintf("%s by %s, /assign_%d", task.Description, task.AssignName, task.TaskID)
		}
	}

	return answer, nil
}

func (s *Service) NewTask(description, assingName string, ownerID int) (map[int]string, error) {
	return s.store.NewTask(description, assingName, ownerID)
}

func (s *Service) Assing(taskID, assignedID int) error {
	return s.store.Assing(taskID, assignedID)
}

func (s *Service) Unassing(taskID int) error {
	return s.store.Unassing(taskID)
}
