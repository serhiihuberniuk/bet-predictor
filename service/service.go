package service

import (
	"context"

	"github.com/serhiihuberniuk/bet-predictor/models"
)

type repository interface {
	CreateLeague(ctx context.Context, league *models.League) error
	GetLeagueByID(ctx context.Context, id string) (*models.League, error)
	DeleteLeague(ctx context.Context, id string) error
	ListLeagues(ctx context.Context) ([]*models.League, error)
}

type fetcher interface {
	AllLeaguesList(ctx context.Context) ([]*models.League, error)
}

type Service struct {
	repo    repository
	fetcher fetcher
}

func New(r repository, f fetcher) *Service {
	return &Service{
		repo:    r,
		fetcher: f,
	}
}
