package worker

import (
	"encoding/json"
	"net/http"
)

type executeRequest struct {
	JobRunID string `json:"job_run_id"`
}

func HandleExecute(w http.ResponseWriter, r *http.Request) {
	var req executeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	runID, err := bson.M
}
