package models

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type DAGStatus string

const (
	DAGRunning DAGStatus = "RUNNING"
	DAGFailed  DAGStatus = "FAILED"
	DAGSuccess DAGStatus = "SUCCESS"
)

type DAG struct {
	ID        bson.ObjectID `bson:"_id,omitempty" json:"id"`
	Name      string        `bson:"name" json:"name"`
	Status    DAGStatus     `bson:"status" json:"status"`
	CreatedAt time.Time     `bson:"created_at" json:"created_at"`
}
