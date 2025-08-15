package api

import (
	"net/http"

	"github.com/Halturshik/final-project/GO/db"
)

type TasksResp struct {
	Tasks []*db.Task `json:"tasks"`
}

func (api *API) tasksHandler(w http.ResponseWriter, r *http.Request) {
	search := r.URL.Query().Get("search")
	tasks, err := api.Store.Tasks(25, search)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]any{"error": err.Error()})
		return
	}
	if tasks == nil {
		tasks = []*db.Task{}
	}
	writeJSON(w, http.StatusOK, TasksResp{Tasks: tasks})
}
