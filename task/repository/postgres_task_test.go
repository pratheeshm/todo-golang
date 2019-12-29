package repository

import (
	"database/sql"
	"reflect"
	"regexp"
	"testing"

	"github.com/sirupsen/logrus"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/pratheeshm/todo-golang/models"
)

func Test_postgresTaskRepository_Add(t *testing.T) {
	query := "INSERT INTO task(title, status) values(?, ?)"
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
					Status: 0,
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs("Take maths notes", 0).WillReturnRows()
			p := &postgresTaskRepository{
				DB: tt.fields.DB,
			}
			if err := p.Add(tt.args.task); (err != nil) != tt.wantErr {
				t.Errorf("postgresTaskRepository.Add() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_postgresTaskRepository_List(t *testing.T) {
	query := "SELECT * FROM task"
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
				Status: 0,
				Title:  "Make maths note",
			}, &models.Task{
				ID:     1,
				Status: 0,
				Title:  "do physics homework",
			}},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &postgresTaskRepository{
				DB: tt.fields.DB,
			}
			rows := mock.NewRows([]string{"id_task", "status", "title"}).AddRow(tt.want[0].ID, tt.want[0].Status, tt.want[0].Title).
				AddRow(tt.want[1].ID, tt.want[1].Status, tt.want[1].Title)
			mock.ExpectQuery(regexp.QuoteMeta(query)).WillReturnRows(rows)
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
