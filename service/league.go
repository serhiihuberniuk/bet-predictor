package service

import (
	"context"
	"fmt"

	"github.com/serhiihuberniuk/bet-predictor/models"
)

func (s *Service) CreateLeague(ctx context.Context, payload *models.CreateLeaguePayload) (string, error) {
	league := &models.League{
		Name:    payload.Name,
		Country: payload.Country,
	}

	league.SetCountrySlug()
	if err := league.SetSlug(); err != nil {
		return "", fmt.Errorf("error while setting slug: %w", err)
	}

	if err := league.SetID(); err != nil {
		return "", fmt.Errorf("error while setting ID: %w", err)
	}

	if err := s.repo.CreateLeague(ctx, league); err != nil {
		return "", fmt.Errorf("error while creating league in db: %w", err)
	}

	return league.ID, nil
}

func (s *Service) GetLeagueByID(ctx context.Context, id string) (*models.League, error) {

	league, err := s.repo.GetLeagueByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("error while getting league from db: %w", err)
	}
	return league, nil
}

func (s *Service) DeleteLeague(ctx context.Context, id string) error {
	if err := s.repo.DeleteLeague(ctx, id); err != nil {
		return fmt.Errorf("error while deleting league from db: %w", err)
	}
	return nil
}

func (s *Service) ListLeagues(ctx context.Context) ([]*models.League, error) {
	leagues, err := s.repo.ListLeagues(ctx)
	if err != nil {
		return nil, fmt.Errorf("error while getting list of leagues from db: %w", err)
	}
	return leagues, nil
}

func (s *Service) CompareLeaguesLists(ctx context.Context,
	remoteLeaguesList []*models.League) ([]*models.League, []*models.League, error) {
	currentLeaguesList, err := s.ListLeagues(ctx)
	if err != nil {
		return nil, nil, fmt.Errorf("error while getting current leagues list: %w", err)
	}

	leaguesToDelete := findMissedLeaguesInFirstList(remoteLeaguesList, currentLeaguesList)
	leaguesToDownload := findMissedLeaguesInFirstList(currentLeaguesList, remoteLeaguesList)

	return leaguesToDownload, leaguesToDelete, nil
}

func findMissedLeaguesInFirstList(firstList []*models.League, secondList []*models.League) []*models.League {
	leaguesMap := make(map[string]struct{})
	for _, v := range firstList {
		leaguesMap[v.Country+v.Name] = struct{}{}
	}

	var missed []*models.League
	for _, v := range secondList {
		_, ok := leaguesMap[v.Country+v.Name]
		if !ok {
			missed = append(missed, v)
		}
	}

	return missed
}
