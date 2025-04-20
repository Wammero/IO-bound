package handler

import (
	"net/http"

	"github.com/Wammero/IO-bound/api/internal/service"
)

type taskHandler struct {
	service service.TaskService
}

func NewTaskHandler(service service.TaskService) *taskHandler {
	return &taskHandler{service: service}
}

func (h *taskHandler) CreateTask(w http.ResponseWriter, r *http.Request) {

}

func (h *taskHandler) GetTaskByID(w http.ResponseWriter, r *http.Request) {

}
