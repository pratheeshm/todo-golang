package usecase

import (
	"errors"
	"reflect"
	"testing"

	"github.com/pratheeshm/todo-golang/models"
	"github.com/pratheeshm/todo-golang/task"
	"github.com/pratheeshm/todo-golang/task/mocks"
)

func TestNewTaskUsecase(t *testing.T) {
	type args struct {
		tr task.Repository
	}
	tests := []struct {
		name string
		args args
		want task.Usecase
	}{{
		name: "Normal Test1: Returning value of type task.Usecase",
		args: args{tr: &mocks.MockRepository{}},
		want: &taskUsecase{taskRepo: &mocks.MockRepository{}},
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewTaskUsecase(tt.args.tr); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewTaskUsecase() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_taskUsecase_Add(t *testing.T) {
	type fields struct {
		taskRepo task.Repository
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
			name: "Normal Case:1 Add task",
			fields: fields{
				taskRepo: &mocks.MockRepository{},
			},
			args: args{
				task: &models.Task{
					Status: 0,
					Title:  "Take maths note",
				},
			},
			wantErr: false,
		},
		{
			name: "Case:2 Add task but repository function returns error",
			fields: fields{
				taskRepo: &mocks.MockRepository{
					Error: errors.New("Repository.Error()"),
				},
			},
			args: args{
				task: &models.Task{
					Status: 0,
					Title:  "Take maths note",
				},
			},
			wantErr: true,
		}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tu := NewTaskUsecase(tt.fields.taskRepo)
			if err := tu.Add(tt.args.task); (err != nil) != tt.wantErr {
				t.Errorf("taskUsecase.Add() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_taskUsecase_Delete(t *testing.T) {
	type fields struct {
		taskRepo task.Repository
	}
	type args struct {
		id int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{{
		name: "Normal Case1: Delete Task",
		fields: fields{
			taskRepo: &mocks.MockRepository{},
		},
		args:    args{id: 1},
		wantErr: false,
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tu := NewTaskUsecase(tt.fields.taskRepo)
			if err := tu.Delete(tt.args.id); (err != nil) != tt.wantErr {
				t.Errorf("taskUsecase.Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_taskUsecase_Edit(t *testing.T) {
	type fields struct {
		taskRepo task.Repository
	}
	type args struct {
		task *models.Task
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{{
		name: "Normal Case1: Edit task",
		fields: fields{
			taskRepo: &mocks.MockRepository{},
		},
		args: args{
			task: &models.Task{
				ID:     0,
				Title:  "Take Maths Note",
				Status: 1,
			},
		},
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tu := NewTaskUsecase(tt.fields.taskRepo)
			if err := tu.Edit(tt.args.task); (err != nil) != tt.wantErr {
				t.Errorf("taskUsecase.Edit() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_taskUsecase_List(t *testing.T) {
	type fields struct {
		taskRepo task.Repository
	}
	tests := []struct {
		name    string
		fields  fields
		want    []*models.Task
		wantErr bool
	}{{
		name: "Normal case1: List tasks",
		fields: fields{
			taskRepo: &mocks.MockRepository{
				Tasks: []*models.Task{
					{
						ID:     0,
						Status: 0,
						Title:  "Take Math notes",
					},
				},
			},
		},
		want: []*models.Task{
			{
				ID:     0,
				Status: 0,
				Title:  "Take Math notes",
			},
		},
		wantErr: false,
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tu := NewTaskUsecase(tt.fields.taskRepo)
			got, err := tu.List()
			if (err != nil) != tt.wantErr {
				t.Errorf("taskUsecase.List() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("taskUsecase.List() = %v, want %v", got, tt.want)
			}
		})
	}
}
