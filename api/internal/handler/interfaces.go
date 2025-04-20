package handler

import "net/http"

type TaskHandler interface {
	CreateTask(w http.ResponseWriter, r *http.Request)
	GetTaskByID(w http.ResponseWriter, r *http.Request)
}
