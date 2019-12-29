package repository

import (
	"database/sql"

	"github.com/sirupsen/logrus"

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
	_, err := p.DB.Query("INSERT INTO task(title, status) values(?, ?)",
		task.Title, task.Status)
	return err
}
func (p *postgresTaskRepository) List() ([]*models.Task, error) {
	tasks := make([]*models.Task, 0)
	rows, err := p.DB.Query("SELECT * FROM task")
	if err != nil {
		return tasks, err
	}
	defer func() {
		err = rows.Close()
		if err != nil {
			logrus.Warn(err)
		}
	}()
	for rows.Next() {
		task := &models.Task{}
		err := rows.Scan(&task.ID, &task.Status, &task.Title)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}
	return tasks, err
}
func (p *postgresTaskRepository) Delete(id int) error {
	_, err := p.DB.Query("DELETE FROM task where id_task = ?", id)
	return err
}
func (p *postgresTaskRepository) Edit(task *models.Task) error {
	_, err := p.DB.Query("UPDATE task SET status = ?, title = ? where id_task = ?",
		task.Status, task.Title, task.ID)
	return err
}
