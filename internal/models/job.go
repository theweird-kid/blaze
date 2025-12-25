package models

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type JobType string

const (
	JobTypeImmediate JobType = "immediate"
	JobTypeDelayed   JobType = "delayed"
	JobTypeCron      JobType = "cron"
)

type HTTPConfig struct {
	Method     string            `bson:"method" json:"method"`
	URL        string            `bson:"url" json:"url"`
	Headers    map[string]string `bson:"headers,omitempty" json:"headers,omitempty"`
	Body       []byte            `bson:"body,omitempty" json:"body,omitempty"`
	TimeoutSec int               `bson:"timeout_sec" json:"timeout_sec"`
}

type Job struct {
	ID   bson.ObjectID `bson:"_id,omitempty" json:"id"`
	Type JobType       `bson:"type" json:"type"`
	HTTP HTTPConfig    `bson:"http" json:"http"`

	// Scheduling
	RunAt     *time.Time `bson:"run_at,omitempty" json:"run_at,omitempty"`
	CronExpr  string     `bson:"cron_expr,omitempty" json:"cron_expr,omitempty"`
	NextRunAt *time.Time `bson:"next_run_at,omitempty" json:"next_run_at,omitempty"`

	// Retry
	MaxAttempts int `bson:"max_attempts" json:"max_attempts"`
	BackoffSec  int `bson:"backoff_sec" json:"backoff_sec"`

	CreatedAt time.Time `bson:"created_at" json:"created_at"`
}
