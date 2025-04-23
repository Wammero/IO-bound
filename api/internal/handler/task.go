package handler

import (
	"net/http"

	"github.com/Wammero/IO-bound/api/internal/service"
	"github.com/Wammero/IO-bound/api/pkg/responsemaker"
	"github.com/google/uuid"
)

type taskHandler struct {
	service service.TaskService
}

func NewTaskHandler(service service.TaskService) *taskHandler {
	return &taskHandler{service: service}
}

func (h *taskHandler) CreateTask(w http.ResponseWriter, r *http.Request) {
	id := uuid.New().String()
	id, err := h.service.CreateTask(r.Context(), id, "test", `"test"`)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	responsemaker.WriteJSONResponse(w, id, http.StatusCreated)
}

func (h *taskHandler) GetTaskByID(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	task, err := h.service.GetTaskByID(r.Context(), id)
	if err != nil {
		responsemaker.WriteJSONError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	responsemaker.WriteJSONResponse(w, task, http.StatusOK)
}
