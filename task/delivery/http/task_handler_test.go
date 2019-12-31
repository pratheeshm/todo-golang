package http

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http/httptest"
	"testing"

	"github.com/pratheeshm/todo-golang/task/mocks"

	nethttp "net/http"

	"github.com/pratheeshm/todo-golang/task"
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
				t.Fatalf("expected msg %v, but got %v", msg, tt.message)
			}
		})
	}
}
