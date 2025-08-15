package api

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/Halturshik/final-project/GO/db"
)

func (api *API) updateTaskHandler(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]any{"error": err.Error()})
		return
	}

	var task db.Task
	if err := json.Unmarshal(body, &task); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": err.Error()})
		return
	}

	if task.ID == "" {
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": "Не указан идентификатор"})
		return
	}
	if task.Title == "" {
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": "Не указан заголовок задачи"})
		return
	}

	if err := checkDate(&task); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": err.Error()})
		return
	}

	if err := api.Store.UpdateTask(&task); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{})
}
