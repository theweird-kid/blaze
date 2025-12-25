package repository

import (
	"context"

	"github.com/theweird-kid/blaze/internal/models"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type JobRepository interface {
	Create(ctx context.Context, job *models.Job) error
	FindRunnable(ctx context.Context) ([]models.Job, error)
}

type JobRunRepository interface {
	Create(ctx context.Context, run *models.JobRun) error
	AcquireLease(ctx context.Context, runID bson.ObjectID) (*models.JobRun, error)
	MarkSuccess(ctx context.Context, runID bson.ObjectID) error
	MarkFailure(ctx context.Context, runID bson.ObjectID, nextRetryAt *bson.DateTime) error
	FindExpiredLeases(ctx context.Context) ([]models.JobRun, error)
}
