package repository

import (
	"database/sql"

	"github.com/pratheeshm/todo-golang/core"
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
	_, err := p.DB.Query("INSERT INTO task(title, status) values($1, $2)",
		task.Title, task.Status)
	return err
}
func (p *postgresTaskRepository) List() ([]*models.Task, error) {
	tasks := make([]*models.Task, 0)
	rows, err := p.DB.Query("SELECT id_task, status, title FROM task")
	if err != nil {
		return tasks, err
	}
	defer rows.Close()
	for rows.Next() {
		task := &models.Task{}
		rows.Scan(&task.ID, &task.Status, &task.Title)
		tasks = append(tasks, task)
	}
	if err = rows.Err(); err != nil {
		return []*models.Task{}, err
	}
	return tasks, err
}
func (p *postgresTaskRepository) Delete(id int) error {
	result, err := p.DB.Exec("DELETE FROM task where id_task = $1", id)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if rows == 0 {
		return core.ErrRecordNotFound
	}
	return err
}
func (p *postgresTaskRepository) Edit(task *models.Task) error {
	result, err := p.DB.Exec("UPDATE task SET status = $1 , title = $2 where id_task = $3",
		task.Status, task.Title, task.ID)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if rows == 0 {
		return core.ErrRecordNotFound
	}
	return err
}
