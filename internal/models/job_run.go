package models

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type JobRunStatus string

const (
	JobRunPending JobRunStatus = "PENDING"
	JobRunRunning JobRunStatus = "RUNNING"
	JobRunSuccess JobRunStatus = "SUCCESS"
	JobRunFailed  JobRunStatus = "FAILED"
	JobRunDead    JobRunStatus = "DEAD"
)

type JobRun struct {
	ID    bson.ObjectID `bson:"_id,omitempty" json:"id"`
	JobID bson.ObjectID `bson:"job_id" json:"job_id"`

	Status JobRunStatus `bson:"status" json:"status"`

	Attempt int `bson:"attempt" json:"attempt"`

	// Leasing
	LeaseUntil *time.Time `bson:"lease_until,omitempty" json:"lease_until,omitempty"`

	// Retry scheduling
	NextRetryAt *time.Time `bson:"next_retry_at,omitempty" json:"next_retry_at,omitempty"`

	// Idempotency
	IdempotencyKey string `bson:"idempotency_key" json:"idempotency_key"`

	CreatedAt time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time `bson:"updated_at" json:"updated_at"`
}
