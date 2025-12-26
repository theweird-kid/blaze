package worker

import (
	"context"
	"encoding/json"
	"log/slog"
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
	start := time.Now()

	var req executeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		slog.Warn("invalid request body", "error", err)
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	logger := slog.With("job_run_id", req.JobRunID)
	// logger.Info("received execute request")

	runID, err := bson.ObjectIDFromHex(req.JobRunID)
	if err != nil {
		logger.Warn("invalid job_run_id format", "error", err)
		http.Error(w, "inavlid job_run_id", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Minute)
	defer cancel()

	client, err := db.Connect(ctx, config.Load().MongoURI)
	if err != nil {
		logger.Error("db connection failed", "error", err)
		http.Error(w, "db error", http.StatusInternalServerError)
		return
	}
	defer client.Disconnect(ctx)

	jobRunRepo := repository.NewJobRunRepo(client.Database(config.Load().DBName))
	jobRepo := repository.NewJobRepo(client.Database(config.Load().DBName))

	// 1) Acquire lease
	jobRun, err := jobRunRepo.AcquireLease(ctx, runID)
	if err != nil || jobRun == nil {
		// Some other worker is executing
		// logger.Info("failed to acquire lease (already running or leased)")
		w.WriteHeader(http.StatusOK)
		return
	}

	logger = logger.With("job_id", jobRun.JobID.Hex())
	logger.Info("executing job")

	// 2) Load job definition
	job, err := jobRepo.FindJobDefinition(ctx, jobRun.JobID)
	if err != nil {
		logger.Error("failed to load job definition", "error", err)
		_ = jobRunRepo.MarkFailure(ctx, runID, nil)
		w.WriteHeader(http.StatusOK)
		return
	}

	// 3) Execute HTTP request
	err = executeHTTP(ctx, job, jobRun)
	duration := time.Since(start)

	// 4) Report result
	if err == nil {
		logger.Info("job finished successfully", "duration_ms", duration.Milliseconds())
		_ = jobRunRepo.MarkSuccess(ctx, runID)
	} else {
		logger.Error("job failed", "error", err, "duration_ms", duration.Milliseconds())
		nextRetryAt := computeNextRetry(job, jobRun)
		_ = jobRunRepo.MarkFailure(ctx, runID, nextRetryAt)
	}

	w.WriteHeader(http.StatusOK)
}
