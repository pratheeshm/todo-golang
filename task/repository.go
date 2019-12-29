package task

import (
	"github.com/pratheeshm/todo-golang/models"
)

//Repository represents task's interface
type Repository interface {
	Add(*models.Task) error
	Delete(int) error
	Edit(*models.Task) error
	List() ([]*models.Task, error)
}
