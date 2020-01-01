package http

import (
	"encoding/json"
	nethttp "net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/go-playground/validator/v10"
	"github.com/pratheeshm/todo-golang/core"
	"github.com/pratheeshm/todo-golang/models"
	"github.com/pratheeshm/todo-golang/task"
	"github.com/sirupsen/logrus"
)

//TaskHandler represents http handler for task
type TaskHandler struct {
	TaskUsecase task.Usecase
}

// NewTaskHandler will initialize the task/ resources endpoint
func NewTaskHandler(tu task.Usecase) nethttp.Handler {
	r := chi.NewMux()
	taskHandler := &TaskHandler{
		TaskUsecase: tu,
	}
	r.Post("/add", taskHandler.Add)
	r.Get("/list", taskHandler.List)
	r.Put("/task/{id:[0-9]+}", taskHandler.Edit)
	r.Delete("/task/{id:[0-9]+}", taskHandler.Delete)
	return r
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
	tasks, err := h.TaskUsecase.List()
	if err != nil {
		w.WriteHeader(nethttp.StatusInternalServerError)
		w.Write([]byte("internal server error"))
		return
	}
	w.WriteHeader(nethttp.StatusOK)
	res, _ := json.Marshal(map[string]interface{}{
		"message": "success",
		"tasks":   tasks,
	})
	w.Write(res)
}

//Edit handler
func (h *TaskHandler) Edit(w nethttp.ResponseWriter, r *nethttp.Request) {
	if chi.URLParam(r, "id") == "" {
		w.WriteHeader(nethttp.StatusBadRequest)
		w.Write([]byte("id is empty"))
		return
	}
	id, _ := strconv.Atoi(chi.URLParam(r, "id"))
	task := &models.Task{}
	d := json.NewDecoder(r.Body)
	err := d.Decode(task)
	task.ID = id
	if err != nil {
		logrus.Error(err)
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
	err = h.TaskUsecase.Edit(task)
	if err != nil {
		if err == core.ErrRecordNotFound {
			w.WriteHeader(nethttp.StatusBadRequest)
			w.Write([]byte("failure"))
			return
		}
		logrus.Error(err)
		w.WriteHeader(nethttp.StatusInternalServerError)
		w.Write([]byte("internal server error"))
		return
	}
	w.WriteHeader(nethttp.StatusOK)
	w.Write([]byte("success"))
}

//Delete handler
func (h *TaskHandler) Delete(w nethttp.ResponseWriter, r *nethttp.Request) {
	if chi.URLParam(r, "id") == "" {
		w.WriteHeader(nethttp.StatusBadRequest)
		w.Write([]byte("id is empty"))
		return
	}
	id, _ := strconv.Atoi(chi.URLParam(r, "id"))
	err := h.TaskUsecase.Delete(id)
	if err != nil {
		if err == core.ErrRecordNotFound {
			w.WriteHeader(nethttp.StatusBadRequest)
			w.Write([]byte("failure"))
			return
		}
		logrus.Error(err)
		w.WriteHeader(nethttp.StatusInternalServerError)
		w.Write([]byte("internal server error"))
		return
	}
	w.WriteHeader(nethttp.StatusOK)
	w.Write([]byte("success"))
}
