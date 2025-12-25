package repository

import "go.mongodb.org/mongo-driver/v2/mongo"

type MongoRepo struct {
	DB *mongo.Database
}

func NewMongoRepo(db *mongo.Database) *MongoRepo {
	return &MongoRepo{DB: db}
}
