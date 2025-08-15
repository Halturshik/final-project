package api

import (
	"encoding/json"
	"net/http"
	"time"
)

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	resp, err := json.Marshal(data)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		resp = []byte(`{"error":"ошибка при формировании JSON"}`)
	} else {
		w.WriteHeader(status)
	}

	w.Write(resp)
}

func afterNow(date1, date2 time.Time) bool {
	date1 = date1.Truncate(24 * time.Hour)
	date2 = date2.Truncate(24 * time.Hour)
	return date1.After(date2)
}
