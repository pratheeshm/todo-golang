package http

import (
	"encoding/json"
	nethttp "net/http"

	"github.com/go-chi/chi"
	"github.com/go-playground/validator/v10"
	"github.com/pratheeshm/todo-golang/models"
	"github.com/pratheeshm/todo-golang/task"
)

//TaskHandler represents http handler for task
type TaskHandler struct {
	TaskUsecase task.Usecase
}

// NewTaskHandler will initialize the task/ resources endpoint
func NewTaskHandler(r chi.Router, tu task.Usecase) {
	taskHandler := &TaskHandler{
		TaskUsecase: tu,
	}
	r.Post("/add", taskHandler.Add)
	r.Get("/list", taskHandler.List)
	r.Put("/edit", taskHandler.Edit)
	r.Delete("/task/{id:[0-9]+}", taskHandler.Delete)
}

//Add task handler
func (h *TaskHandler) Add(w nethttp.ResponseWriter, r *nethttp.Request) {
	task := &models.Task{}
	d := json.NewDecoder(r.Body)
	err := d.Decode(task)
	if err != nil {
		w.WriteHeader(nethttp.StatusBadRequest)
		w.Write([]byte("Can not decode body"))
		return
	}
	validate := validator.New()
	err = validate.Struct(task)
	if err != nil {
		w.WriteHeader(nethttp.StatusBadRequest)
		w.Write([]byte("validation error"))
		return
	}
	err = h.TaskUsecase.Add(task)
	if err != nil {
		w.WriteHeader(nethttp.StatusInternalServerError)
		w.Write([]byte("internal server error"))
		return
	}
	w.WriteHeader(nethttp.StatusOK)
	w.Write([]byte("success"))
}

//List handler
func (h *TaskHandler) List(w nethttp.ResponseWriter, r *nethttp.Request) {

}

//Edit handler
func (h *TaskHandler) Edit(w nethttp.ResponseWriter, r *nethttp.Request) {

}

//Delete handler
func (h *TaskHandler) Delete(w nethttp.ResponseWriter, r *nethttp.Request) {

}
