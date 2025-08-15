package api

import (
	"github.com/Halturshik/final-project/GO/db"
	"github.com/go-chi/chi/v5"
)

const DateFormat = "20060102"

type API struct {
	Store *db.Store
}

func NewAPI(store *db.Store) *API {
	return &API{Store: store}
}

func (api *API) Init(r *chi.Mux) {
	r.Get("/api/nextdate", nextDateHandler)
	r.Post("/api/task", api.addTaskHandler)
	r.Get("/api/task", api.getTaskHandler)
	r.Put("/api/task", api.updateTaskHandler)
	r.Delete("/api/task", api.deleteTaskHandler)
	r.Post("/api/task/done", api.doneTaskHandler)
	r.Get("/api/tasks", api.tasksHandler)
}
