package service

import (
	"context"

	"github.com/serhiihuberniuk/bet-predictor/models"
)

type repository interface {
	CreateLeague(ctx context.Context, league *models.League) error
	GetLeagueByCountryAndName(ctx context.Context, countrySlug, slug string) (*models.League, error)
	DeleteLeague(ctx context.Context, id string) error
	ListLeagues(ctx context.Context) ([]*models.League, error)

	CreateTeam(ctx context.Context, team *models.Team) error
	DeleteTeam(ctx context.Context, teamID string) error
	ListTeams(ctx context.Context) ([]*models.Team, error)
}

type fetcher interface {
	AllLeaguesList(ctx context.Context) ([]*models.League, error)
	GetTeamsBySeasonID(ctx context.Context, seasonID int) ([]*models.Team, error)
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
