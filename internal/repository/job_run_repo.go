package repository

import (
	"context"
	"time"

	"github.com/theweird-kid/blaze/internal/models"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type JobRunRepo struct {
	*MongoRepo
	col *mongo.Collection
}

func NewJobRunRepo(db *mongo.Database) *JobRunRepo {
	return &JobRunRepo{
		MongoRepo: NewMongoRepo(db),
		col:       db.Collection("job_runs"),
	}
}

func (r *JobRunRepo) Create(ctx context.Context, run *models.JobRun) error {
	now := time.Now().UTC()
	run.Status = models.JobRunPending
	run.Attempt = 0
	run.CreatedAt = now
	run.UpdatedAt = now

	_, err := r.col.InsertOne(ctx, run)
	return err
}

func (r *JobRunRepo) AcquireLease(ctx context.Context, runID bson.ObjectID) (*models.JobRun, error) {
	leaseUntil := time.Now().UTC().Add(30 * time.Second)

	filter := bson.M{
		"_id":    runID,
		"status": models.JobRunPending,
	}

	update := bson.M{
		"$set": bson.M{
			"status":      models.JobRunRunning,
			"lease_until": leaseUntil,
			"updated_at":  time.Now().UTC(),
		},
		"$inc": bson.M{
			"attempt": 1,
		},
	}

	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)

	var run models.JobRun
	err := r.col.FindOneAndUpdate(ctx, filter, update, opts).Decode(&run)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &run, nil
}

func (r *JobRunRepo) MarkSuccess(ctx context.Context, runID bson.ObjectID) error {
	filter := bson.M{
		"_id":         runID,
		"status":      models.JobRunRunning,
		"lease_until": bson.M{"$gt": time.Now().UTC()},
	}

	update := bson.M{
		"$set": bson.M{
			"status":     models.JobRunSuccess,
			"updated_at": time.Now().UTC(),
		},
	}

	res, err := r.col.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	if res.MatchedCount == 0 {
		return nil // lease expired or stolen
	}
	return nil
}

func (r *JobRunRepo) MarkFailure(ctx context.Context, runID bson.ObjectID, nextRetryAt bson.DateTime) error {
	update := bson.M{
		"$set": bson.M{
			"status":        models.JobRunFailed,
			"next_retry_at": nextRetryAt,
			"updated_at":    time.Now().UTC(),
		},
	}

	_, err := r.col.UpdateOne(ctx, runID, update)
	return err
}

func (r *JobRunRepo) FindExpiredLeases(ctx context.Context) ([]models.JobRun, error) {
	filter := bson.M{
		"status":      models.JobRunFailed,
		"lease_until": bson.M{"$lt": time.Now().UTC()},
	}

	cursor, err := r.col.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var runs []models.JobRun
	if err := cursor.All(ctx, &runs); err != nil {
		return nil, err
	}

	return runs, nil
}
