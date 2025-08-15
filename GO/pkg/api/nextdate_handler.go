package api

import (
	"net/http"
	"time"
)

func nextDateHandler(w http.ResponseWriter, r *http.Request) {
	nowStr := r.URL.Query().Get("now")
	dateStr := r.URL.Query().Get("date")
	repeat := r.URL.Query().Get("repeat")

	var now time.Time
	var err error
	if nowStr == "" {
		now = time.Now()
	} else {
		now, err = time.Parse(DateFormat, nowStr)
		if err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]any{"error": "Некорректный формат now"})
			return
		}
	}

	if dateStr == "" {
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": "Параметр date обязателен"})
		return
	}

	next, err := NextDate(now, dateStr, repeat)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": err.Error()})
		return
	}

	// я бы оставил следующую реализацию, чтобы все выглядело логично и последовательно -
	// writeJSON(w, http.StatusOK, map[string]any{"next": next})
	// но для успешного прохождения тестов она не подходит
	w.Header().Set("Content-Type", "text/plain; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(next))
}
