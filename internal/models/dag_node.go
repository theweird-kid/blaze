package models

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type DAGNodeStatus string

const (
	DAGNodePending DAGNodeStatus = "PENDING"
	DAGNodeRunning DAGNodeStatus = "RUNNING"
	DAGNodeSuccess DAGNodeStatus = "SUCCESS"
	DAGNodeFailed  DAGNodeStatus = "FAILED"
)

type DAGNode struct {
	ID    bson.ObjectID `bson:"_id,omitempty" json:"id"`
	DAGID bson.ObjectID `bson:"dag_id" json:"dag_id"`

	Name string `bson:"name" json:"name"`

	JobID bson.ObjectID `bson:"job_id" json:"job_id"`

	Parents  []bson.ObjectID `bson:"parents" json:"parents"`
	Children []bson.ObjectID `bson:"children" json:"children"`

	RemainingDeps int `bson:"remaining_deps" json:"remaining_deps"`

	Status DAGNodeStatus `bson:"status" json:"status"`

	CreatedAt time.Time `bson:"created_at" json:"created_at"`
}
