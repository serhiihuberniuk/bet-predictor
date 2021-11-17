package repository

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
	"github.com/serhiihuberniuk/bet-predictor/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func (r *Repository) CreateLeague(ctx context.Context, league *models.League) error {
	_, err := r.useLeagueCollection().InsertOne(ctx, league)
	if err != nil {
		return fmt.Errorf("error while inserting: %w", err)
	}

	return nil
}

func (r *Repository) GetLeagueByCountryAndName(ctx context.Context, countrySlug, slug string) (*models.League, error) {
	var league models.League

	if err := r.useLeagueCollection().FindOne(ctx, bson.M{"slug": slug, "country_slug": countrySlug}).Decode(&league); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, models.ErrNotFound
		}

		return nil, fmt.Errorf("error while getting from db: %w", err)
	}

	return &league, nil
}

func (r *Repository) DeleteLeague(ctx context.Context, leagueID string) error {

	if _, err := r.useLeagueCollection().DeleteOne(ctx, bson.M{"_id": leagueID}); err != nil {
		return fmt.Errorf("error while deleting from db: %w", err)
	}
	return nil
}

func (r *Repository) ListLeagues(ctx context.Context) ([]*models.League, error) {

	cursor, err := r.useLeagueCollection().Find(ctx, bson.M{})
	if err != nil {
		return nil, fmt.Errorf("error while finding leagues: %w", err)
	}

	var leagues []*models.League
	if err = cursor.All(ctx, &leagues); err != nil {
		if err != nil {
			return nil, fmt.Errorf("error while decoding: %w", err)
		}
	}

	return leagues, nil
}

func (r *Repository) useLeagueCollection() *mongo.Collection {
	return r.db.Collection("league")
}
