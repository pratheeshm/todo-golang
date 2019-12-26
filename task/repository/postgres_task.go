package repository

import (
	"database/sql"

	"github.com/pratheeshm/todo-golang/models"

	"github.com/pratheeshm/todo-golang/task"
)

type postgresTaskRepository struct {
	*sql.DB
}

// NewPostgresTaskRepository will create an object that represent the task.Repository interface
func NewPostgresTaskRepository(db *sql.DB) task.Repository {
	return &postgresTaskRepository{db}
}
func (p *postgresTaskRepository) Add(task *models.Task) error {
	_, err := p.DB.Query("INSERT INTO task(ida_task, title, status) values(?, ?, ?)",
		task.ID, task.Title, task.Status)
	return err
}
func (p *postgresTaskRepository) List() {

}
func (p *postgresTaskRepository) Delete() {

}
func (p *postgresTaskRepository) Edit() {

}
