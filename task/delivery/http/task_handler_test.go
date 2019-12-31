package http

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	nethttp "net/http"
	"net/http/httptest"
	"testing"

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
			"status": 0,
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
			"status": 0,
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
			"status": "4",
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
			"status": 4,
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
				t.Fatalf("expected status %v, but got %v", tt.statusCode, res.StatusCode)
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