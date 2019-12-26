package models

// Task represents the task model
type Task struct {
	ID     int    `json:"id_task"`
	Title  string `json:"title"`
	Status int    `json:"status"`
}
