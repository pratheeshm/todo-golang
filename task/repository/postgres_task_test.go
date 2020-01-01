package repository

import (
	"database/sql"
	"errors"
	"reflect"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/pratheeshm/todo-golang/models"
	"github.com/pratheeshm/todo-golang/task"
	"github.com/sirupsen/logrus"
)

func TestNewPostgresTaskRepository(t *testing.T) {
	type args struct {
		db *sql.DB
	}
	tests := []struct {
		name string
		args args
		want task.Repository
	}{
		{
			name: "Normal Case 1: Return object of type task.Repository",
			args: args{
				db: &sql.DB{},
			},
			want: &postgresTaskRepository{
				DB: &sql.DB{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewPostgresTaskRepository(tt.args.db); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewPostgresTaskRepository() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_postgresTaskRepository_Add(t *testing.T) {
	query := "INSERT INTO task(title, status) values($1, $2)"
	db, mock, err := sqlmock.New()
	if err != nil {
		logrus.Error("expected no error, but got:", err)
		return
	}
	type fields struct {
		DB *sql.DB
	}
	type args struct {
		task *models.Task
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:   "Normal Case 1: Insert task",
			fields: fields{DB: db},
			args: args{
				task: &models.Task{
					ID:     1,
					Title:  "Take maths notes",
					Status: "todo",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock.ExpectQuery(regexp.QuoteMeta(query)).
				WithArgs("Take maths notes", "todo").WillReturnRows()
			p := NewPostgresTaskRepository(tt.fields.DB)
			if err := p.Add(tt.args.task); (err != nil) != tt.wantErr {
				t.Errorf("postgresTaskRepository.Add() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_postgresTaskRepository_List(t *testing.T) {
	query := "SELECT id_task, status, title FROM task"
	db, mock, err := sqlmock.New()
	if err != nil {
		logrus.Error("expected no error, but got:", err)
		return
	}
	type fields struct {
		DB *sql.DB
	}
	tests := []struct {
		name    string
		fields  fields
		want    []*models.Task
		wantErr bool
	}{
		{
			name: "Normal Case 1: List all task",
			fields: fields{
				DB: db,
			},
			want: []*models.Task{&models.Task{
				ID:     0,
				Status: "todo",
				Title:  "Make maths note",
			}, &models.Task{
				ID:     1,
				Status: "todo",
				Title:  "do physics homework",
			}},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewPostgresTaskRepository(tt.fields.DB)
			rows := mock.NewRows([]string{"id_task", "status", "title"}).
				AddRow(tt.want[0].ID, tt.want[0].Status, tt.want[0].Title).
				AddRow(tt.want[1].ID, tt.want[1].Status, tt.want[1].Title)
			mock.ExpectQuery(regexp.QuoteMeta(query)).
				WillReturnRows(rows)
			got, err := p.List()
			if (err != nil) != tt.wantErr {
				t.Errorf("postgresTaskRepository.List() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("postgresTaskRepository.List() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_postgresTaskRepository_Delete(t *testing.T) {
	query := "DELETE FROM task where id_task = $1"
	db, mock, err := sqlmock.New()
	if err != nil {
		logrus.Error(err)
		return
	}
	type fields struct {
		DB *sql.DB
	}
	type args struct {
		id int
	}
	tests := []struct {
		name         string
		fields       fields
		args         args
		wantErr      bool
		rowsAffected int64
		dbError      error
	}{{
		name: "Norml Test 1: Delete a task ",
		fields: fields{
			DB: db,
		},
		args: args{
			id: 1,
		},
		wantErr:      false,
		rowsAffected: 1,
	}, {
		name: "invalid id",
		fields: fields{
			DB: db,
		},
		args: args{
			id: 1,
		},
		wantErr:      true,
		rowsAffected: 0,
		dbError:      errors.New("record not found"),
	}, {
		name: "invalid id",
		fields: fields{
			DB: db,
		},
		args: args{
			id: 1,
		},
		wantErr:      true,
		rowsAffected: 0,
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewPostgresTaskRepository(tt.fields.DB)
			mock.ExpectExec(regexp.QuoteMeta(query)).
				WithArgs(tt.args.id).
				WillReturnResult(sqlmock.NewResult(0, tt.rowsAffected)).
				WillReturnError(tt.dbError)
			if err := p.Delete(tt.args.id); (err != nil) != tt.wantErr {
				t.Errorf("Test %s - got error = %v, wantErr %v", tt.name, err, tt.wantErr)
			}
		})
	}
}

func Test_postgresTaskRepository_Edit(t *testing.T) {
	query := "UPDATE task SET status = $1 , title = $2 where id_task = $3"
	db, mock, err := sqlmock.New()
	if err != nil {
		logrus.Error(err)
		return
	}
	type fields struct {
		DB *sql.DB
	}
	type args struct {
		task *models.Task
	}
	tests := []struct {
		name        string
		fields      fields
		args        args
		wantErr     bool
		rowsUpdated int64
		dbError     error
	}{{
		name: "Normal Case 1: Edit a task Status",
		fields: fields{
			DB: db,
		},
		args: args{
			&models.Task{
				ID:     1,
				Status: "inprogress",
				Title:  "Take math notes",
			},
		},
		wantErr:     false,
		rowsUpdated: 1,
	}, {
		name: "Edit a invalid task",
		fields: fields{
			DB: db,
		},
		args: args{
			&models.Task{
				ID:     1,
				Status: "inprogress",
				Title:  "Take math notes",
			},
		},
		wantErr:     true,
		rowsUpdated: 0,
	}, {
		name: "Edit a task Status but return error",
		fields: fields{
			DB: db,
		},
		args: args{
			&models.Task{
				ID:     1,
				Status: "inprogress",
				Title:  "Take math notes",
			},
		},
		wantErr:     true,
		rowsUpdated: 0,
		dbError:     errors.New("DB error"),
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewPostgresTaskRepository(tt.fields.DB)
			mock.ExpectExec(regexp.QuoteMeta(query)).
				WithArgs(tt.args.task.Status, tt.args.task.Title, tt.args.task.ID).
				WillReturnResult(sqlmock.NewResult(0, tt.rowsUpdated)).
				WillReturnError(tt.dbError)
			err := p.Edit(tt.args.task)
			if (err != nil) != tt.wantErr {
				t.Errorf("Test- %v,error = %v, wantErr %v", tt.name, err, tt.wantErr)
			}
		})
	}
}
