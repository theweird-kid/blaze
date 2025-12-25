package repository

import (
	"context"
	"time"

	"github.com/theweird-kid/blaze/internal/models"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type JobRepo struct {
	*MongoRepo
	col *mongo.Collection
}

func NewJobRepo(db *mongo.Database) *JobRepo {
	return &JobRepo{
		MongoRepo: NewMongoRepo(db),
		col:       db.Collection("jobs"),
	}
}

func (r *JobRepo) Create(ctx context.Context, job *models.Job) error {
	job.CreatedAt = time.Now().UTC()
	_, err := r.col.InsertOne(ctx, job)
	return err
}

func (r *JobRepo) FindRunnable(ctx context.Context) ([]models.Job, error) {
	now := time.Now().UTC()

	filter := bson.M{
		"$or": []bson.M{
			{"run_at": bson.M{"$lte": now}},
			{"next_run_at": bson.M{"$lte": now}},
		},
	}

	cursor, err := r.col.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var jobs []models.Job
	if err := cursor.All(ctx, &jobs); err != nil {
		return nil, err
	}
	return jobs, nil
}

func (r *JobRepo) FindJobDefinition(ctx context.Context, jobID bson.ObjectID) (*models.Job, error) {
	filter := bson.M{
		"_id": jobID,
	}

	var job models.Job
	err := r.col.FindOne(ctx, filter).Decode(&job)
	if err != nil {
		return nil, err
	}
	return &job, nil
}
