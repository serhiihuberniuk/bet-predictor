package repository

import (
	"context"
	"fmt"

	"github.com/serhiihuberniuk/bet-predictor/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func (r *Repository) CreateTeam(ctx context.Context, team *models.Team) error {
	if _, err := r.useTeamsCollection().InsertOne(ctx, team); err != nil {
		return fmt.Errorf("error while inserting team to db: %w", err)
	}

	return nil
}

func (r *Repository) DeleteTeam(ctx context.Context, teamID string) error {
	if _, err := r.useTeamsCollection().DeleteOne(ctx, bson.M{"_id": teamID}); err != nil {
		return fmt.Errorf("error while deleting from db: %w", err)
	}

	return nil
}
func (r *Repository) ListTeams(ctx context.Context) ([]*models.Team, error) {
	cursor, err := r.useTeamsCollection().Find(ctx, bson.M{})
	if err != nil {
		return nil, fmt.Errorf("error while finding teams in db: %w", err)
	}

	var teams []*models.Team

	if err := cursor.All(ctx, &teams); err != nil {
		return nil, fmt.Errorf("error while decoding teams; %w, err")
	}

	return teams, nil
}

func (r *Repository) useTeamsCollection() *mongo.Collection {
	return r.db.Collection("teams")
}
