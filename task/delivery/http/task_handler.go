package http

import (
	"encoding/json"
	nethttp "net/http"

	"github.com/go-chi/chi"
	"github.com/pratheeshm/todo-golang/models"
	"github.com/pratheeshm/todo-golang/task"
)

//TaskHandler represents http handler for task
type TaskHandler struct {
	TaskUsecase task.Usecase
}

//ResponseMessage contains Response message
type ResponseMessage struct {
	Message string
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
		response := ResponseMessage{
			Message: "Could not parse request Body",
		}
		w.WriteHeader(nethttp.StatusBadRequest)
		resData, _ := json.Marshal(response)
		w.Write(resData)
	}
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
