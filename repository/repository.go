package repository

import "go.mongodb.org/mongo-driver/mongo"

type Repository struct {
	db *mongo.Database
}

func New(db *mongo.Database) *Repository {
	return &Repository{
		db: db,
	}
}
