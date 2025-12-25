package worker

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/theweird-kid/blaze/internal/config"
	"github.com/theweird-kid/blaze/internal/db"
	"github.com/theweird-kid/blaze/internal/repository"
	"go.mongodb.org/mongo-driver/v2/bson"
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

	runID, err := bson.ObjectIDFromHex(req.JobRunID)
	if err != nil {
		http.Error(w, "inavlid job_run_id", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Minute)
	defer cancel()

	client, err := db.Connect(ctx, config.Load().MongoURI)
	if err != nil {
		http.Error(w, "db error", http.StatusInternalServerError)
		return
	}

	jobRunRepo := repository.NewJobRunRepo(client.Database(config.Load().DBName))
	jobRepo := repository.NewJobRepo(client.Database(config.Load().DBName))

	// 1) Acquire lease
	jobRun, err := jobRunRepo.AcquireLease(ctx, runID)
	if err != nil || jobRun == nil {
		// Some other worker is executing
		w.WriteHeader(http.StatusOK)
		return
	}

	// 2) Load job definition
	job, err := jobRepo.FindJobDefinition(ctx, jobRun.JobID)
	if err != nil {
		_ = jobRunRepo.MarkFailure(ctx, runID, nil)
		w.WriteHeader(http.StatusOK)
		return
	}

	// 3) Execute HTTP request
	err = executeHTTP(ctx, job, jobRun)

	// 4) Report result
	if err == nil {
		_ = jobRunRepo.MarkSuccess(ctx, runID)
	} else {
		nextRetryAt := computeNextRetry(job, jobRun)
		_ = jobRunRepo.MarkFailure(ctx, runID, nextRetryAt)
	}

	w.WriteHeader(http.StatusOK)
}
