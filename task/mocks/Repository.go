package mocks

import (
	"github.com/pratheeshm/todo-golang/models"
)

//MockRepository implements inerface task.Repository
type MockRepository struct {
	Error error
	Tasks []*models.Task
}

//Delete task
func (m *MockRepository) Delete(id int) error {
	return m.Error
}

//Add task
func (m *MockRepository) Add(task *models.Task) error {
	return m.Error
}

//Edit task
func (m *MockRepository) Edit(task *models.Task) error {
	return m.Error
}

//List tasks
func (m *MockRepository) List() ([]*models.Task, error) {
	return m.Tasks, m.Error
}
