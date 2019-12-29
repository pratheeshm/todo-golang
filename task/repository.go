package task

import (
	"github.com/pratheeshm/todo-golang/models"
)

//Repository represents task's interface
type Repository interface {
	Add(*models.Task) error
	Delete()
	Edit()
	List() ([]*models.Task, error)
}
