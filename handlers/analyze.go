package handlers

import (
	"encoding/json"
	"net/http"

	"cp-advisor/services"
)

func HandleAnalyze(w http.ResponseWriter, r *http.Request) {
	handle := r.URL.Query().Get("handle")
	if handle == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"error": "handle parameter is required"}`))
		return
	}

	result, err := services.AnalyzeUser(handle)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"error": "` + err.Error() + `"}`))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}
