package models

// Task represents the task model
type Task struct {
	ID     int    `json:"id_task"`
	Title  string `json:"title" validate:"required"`
	Status int    `json:"status" validate:"gte=0,lte=2"`
}
