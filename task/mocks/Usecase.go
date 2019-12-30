package mocks

import "github.com/pratheeshm/todo-golang/models"

//MockUsecase implements inerface task.Usecase
type MockUsecase struct {
	Error error
	Tasks []*models.Task
}

//Add task
func (m *MockUsecase) Add(*models.Task) error {
	return m.Error
}

//Delete task
func (m *MockUsecase) Delete(int) error {
	return m.Error
}

//Edit task
func (m *MockUsecase) Edit(*models.Task) error {
	return m.Error
}

//List tasks
func (m *MockUsecase) List() ([]*models.Task, error) {
	return m.Tasks, m.Error
}
