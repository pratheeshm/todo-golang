package http

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	nethttp "net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi"
	"github.com/pratheeshm/todo-golang/core"
	"github.com/pratheeshm/todo-golang/task"
	"github.com/pratheeshm/todo-golang/task/mocks"
)

func TestTaskHandler_Add(t *testing.T) {
	url := "localhost:8080/add"
	type fields struct {
		TaskUsecase task.Usecase
	}
	tests := []struct {
		name       string
		fields     fields
		url        string
		statusCode int
		body       map[string]interface{}
		message    string
	}{{
		name: "Normal case1: ",
		fields: fields{
			TaskUsecase: &mocks.MockUsecase{},
		},
		statusCode: 200,
		body: map[string]interface{}{
			"status": "todo",
			"title":  "Test title",
		},
		message: "success",
	}, {
		name: "Usecase returns error",
		fields: fields{
			TaskUsecase: &mocks.MockUsecase{
				Error: errors.New("Usecase.Error()"),
			},
		},
		statusCode: 500,
		body: map[string]interface{}{
			"status": "todo",
			"title":  "Test title",
		},
		message: "internal server error",
	}, {
		name: "request body parse error",
		fields: fields{
			TaskUsecase: &mocks.MockUsecase{},
		},
		statusCode: 400,
		body: map[string]interface{}{
			"status": 3,
			"title":  "Test title",
		},
		message: "Can not decode body",
	}, {
		name: "validation error",
		fields: fields{
			TaskUsecase: &mocks.MockUsecase{},
		},
		statusCode: 400,
		body: map[string]interface{}{
			"status": "completed",
			"title":  "Test title",
		},
		message: "validation error",
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &TaskHandler{
				TaskUsecase: tt.fields.TaskUsecase,
			}
			bodyBytes, err := json.Marshal(tt.body)
			if err != nil {
				t.Fatalf("got error: %v", err)
			}
			req, err := nethttp.NewRequest("POST", url, bytes.NewBuffer(bodyBytes))
			if err != nil {
				t.Fatalf("got error: %v", err)
			}
			rec := httptest.NewRecorder()
			h.Add(rec, req)
			res := rec.Result()
			if res.StatusCode != tt.statusCode {
				t.Fatalf("Test -%s , expected status %v, but got %v", tt.name, tt.statusCode, res.StatusCode)
			}
			defer res.Body.Close()
			respMsg, _ := ioutil.ReadAll(res.Body)
			if msg := string(respMsg); msg != tt.message {
				t.Fatalf("expected msg %v, but got %v", tt.message, msg)
			}
		})
	}
}

func TestNewTaskHandler(t *testing.T) {
	u := &mocks.MockUsecase{}
	urlStatus := map[bool]string{true: "Found", false: "Not found"}
	server := httptest.NewServer(NewTaskHandler(u))
	defer server.Close()
	baseURL := fmt.Sprintf("%s", server.URL)
	tests := []struct {
		name    string
		method  string
		url     string
		isFound bool
	}{{
		name:    "Normal Test1:=",
		method:  "GET",
		url:     "/list",
		isFound: true,
	}, {
		name:    "invalid endpoint",
		method:  "GET",
		url:     "/abc",
		isFound: false,
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var res *nethttp.Response
			var err error
			switch method := tt.method; method {
			case "GET":
				res, err = nethttp.Get(baseURL + tt.url)
			default:
				t.Fatalf("Method not found")
			}
			if err != nil {
				t.Fatalf("got error: %v", err)
			}
			if notFound := (res.StatusCode == nethttp.StatusNotFound); notFound == tt.isFound {
				t.Fatalf("URL - %v  expected it as %v but got %v",
					tt.url, urlStatus[tt.isFound], urlStatus[!notFound])
			}
		})
	}
}

func TestTaskHandler_List(t *testing.T) {
	type fields struct {
		TaskUsecase task.Usecase
	}
	tests := []struct {
		name       string
		fields     fields
		statusCode int
	}{{
		name: "Success case",
		fields: fields{
			TaskUsecase: &mocks.MockUsecase{},
		},
		statusCode: 200,
	}, {
		name: "failure case",
		fields: fields{
			TaskUsecase: &mocks.MockUsecase{
				Error: errors.New("Usecase.error()"),
			},
		},
		statusCode: 500,
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &TaskHandler{
				TaskUsecase: tt.fields.TaskUsecase,
			}
			req := httptest.NewRequest("GET", "localhost:3000/list", nil)
			rec := httptest.NewRecorder()
			h.List(rec, req)
			res := rec.Result()
			if res.StatusCode != tt.statusCode {
				t.Fatalf("expected statusCode %d but got %d", tt.statusCode, res.StatusCode)
			}
		})
	}
}

func TestTaskHandler_Edit(t *testing.T) {
	type fields struct {
		TaskUsecase task.Usecase
	}
	tests := []struct {
		name       string
		fields     fields
		body       map[string]interface{}
		urlParam   map[string]string
		statusCode int
	}{{
		name:   "Normal Test1",
		fields: fields{TaskUsecase: &mocks.MockUsecase{}},
		body: map[string]interface{}{
			"title":  "Take math notes",
			"status": "todo",
		},
		urlParam: map[string]string{
			"id": "3",
		},
		statusCode: 200,
	}, {
		name:   "urlparam is not set",
		fields: fields{TaskUsecase: &mocks.MockUsecase{}},
		body: map[string]interface{}{
			"title":  "Take math notes",
			"status": "todo",
		},
		statusCode: 400,
	}, {
		name: "usecase error",
		fields: fields{TaskUsecase: &mocks.MockUsecase{
			Error: errors.New("Usecase error"),
		}},
		urlParam: map[string]string{
			"id": "3",
		},
		body: map[string]interface{}{
			"title":  "Take math notes",
			"status": "todo",
		},
		statusCode: 500,
	}, {
		name: "record not found error",
		fields: fields{TaskUsecase: &mocks.MockUsecase{
			Error: core.ErrRecordNotFound,
		}},
		urlParam: map[string]string{
			"id": "3",
		},
		body: map[string]interface{}{
			"title":  "Take math notes",
			"status": "todo",
		},
		statusCode: 400,
	}, {
		name: "status as number",
		fields: fields{TaskUsecase: &mocks.MockUsecase{
			Error: core.ErrRecordNotFound,
		}},
		urlParam: map[string]string{
			"id": "3",
		},
		body: map[string]interface{}{
			"title":  "Take math notes",
			"status": 1,
		},
		statusCode: 400,
	}, {
		name: "status is missing",
		fields: fields{TaskUsecase: &mocks.MockUsecase{
			Error: core.ErrRecordNotFound,
		}},
		urlParam: map[string]string{
			"id": "3",
		},
		body: map[string]interface{}{
			"title": "Take math notes",
		},
		statusCode: 400,
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &TaskHandler{
				TaskUsecase: tt.fields.TaskUsecase,
			}
			bodyBytes, err := json.Marshal(tt.body)
			if err != nil {
				t.Fatalf("got error: %v", err)
			}
			req := httptest.NewRequest("PUT", "/task", bytes.NewBuffer(bodyBytes))
			ctx := chi.NewRouteContext()
			for k, v := range tt.urlParam {
				ctx.URLParams.Add(k, v)
			}
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, ctx))
			rec := httptest.NewRecorder()
			h.Edit(rec, req)
			res := rec.Result()
			if res.StatusCode != tt.statusCode {
				t.Fatalf("Test - %s , got statuscode %d but expected %d",
					tt.name, res.StatusCode, tt.statusCode)
			}
		})
	}
}
