package worker

import (
	"time"

	"github.com/theweird-kid/blaze/internal/models"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func computeNextRetry(job *models.Job, run *models.JobRun) *bson.DateTime {
	if run.Attempt >= job.MaxAttempts {
		return nil
	}

	delay := time.Duration(job.BackoffSec) * time.Second
	next := time.Now().UTC().Add(delay * (1 << run.Attempt))

	dt := bson.NewDateTimeFromTime(next)
	return &dt
}
