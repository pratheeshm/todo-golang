package task

import "github.com/pratheeshm/todo-golang/models"

//Usecase represents task's interface
type Usecase interface {
	Add(*models.Task) error
	Delete(int) error
	Edit(*models.Task) error
	List() ([]*models.Task, error)
}
