package service

import (
	"context"

	"github.com/serhiihuberniuk/bet-predictor/models"
)

type Service struct {
	repo repository
}

func New(r repository) *Service {
	return &Service{
		repo: r,
	}
}

type repository interface {
	CreateLeague(ctx context.Context, league *models.League) error
	GetLeagueByID(ctx context.Context, id string) (*models.League, error)
	DeleteLeague(ctx context.Context, id string) error
	ListLeagues(ctx context.Context) ([]*models.League, error)
}
