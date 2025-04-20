package handler

import (
	"github.com/Wammero/IO-bound/api/internal/service"
	"github.com/go-chi/chi"
)

type Handler struct {
	TaskHandler TaskHandler
}

func New(services *service.Service) *Handler {
	return &Handler{
		TaskHandler: NewTaskHandler(services.TaskService),
	}
}

func (h *Handler) SetupRoutes(r *chi.Mux) {
	r.Route("/tasks", func(r chi.Router) {
		r.Post("/", h.TaskHandler.CreateTask)
		r.Get("/{id}", h.TaskHandler.GetTaskByID)
	})
}
