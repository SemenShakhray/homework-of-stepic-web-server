package models

type Task struct {
	TaskID      int    `json:"task_id"`
	Description string `json:"description"`
	OwnerID     int    `json:"owner_id"`
	AssignedID  int    `json:"assigned_id"`
}
