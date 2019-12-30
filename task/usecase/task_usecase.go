package usecase

import (
	"github.com/pratheeshm/todo-golang/models"
	"github.com/pratheeshm/todo-golang/task"
)

type taskUsecase struct {
	taskRepo task.Repository
}

// NewTaskUsecase will create new a taskUsecase object representation of task.Usecase interface
func NewTaskUsecase(tr task.Repository) task.Usecase {
	return &taskUsecase{
		taskRepo: tr,
	}
}
func (tu *taskUsecase) Add(task *models.Task) error {
	err := tu.taskRepo.Add(task)
	return err
}
func (tu *taskUsecase) Delete(id int) error {
	err := tu.taskRepo.Delete(id)
	return err
}
func (tu *taskUsecase) Edit(task *models.Task) error {
	err := tu.taskRepo.Edit(task)
	return err
}
func (tu *taskUsecase) List() ([]*models.Task, error) {
	tasks, err := tu.taskRepo.List()
	return tasks, err
}
