package service

import (
	"context"
	"fmt"

	"github.com/serhiihuberniuk/bet-predictor/models"
)

func (s *Service) CreateLeague(ctx context.Context, payload models.CreateLeaguePayload) (string, error) {
	league, err := models.NewLeague(payload)
	if err != nil {
		return "", fmt.Errorf("error while creating league: %w", err)
	}

	if err := s.repo.CreateLeague(ctx, league); err != nil {
		return "", fmt.Errorf("error while creating league in db: %w", err)
	}

	return league.ID, nil
}

func (s *Service) GetLeagueByCountryAndName(ctx context.Context, countrySlug, slug string) (*models.League, error) {
	league, err := s.repo.GetLeagueByCountryAndName(ctx, countrySlug, slug)
	if err != nil {
		return nil, fmt.Errorf("error while getting league from repository layer: %w", err)
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

func (s *Service) CompareLeaguesLists(ctx context.Context) (leaguesToDownload []*models.League,
	leaguesToDelete []*models.League, err error) {
	currentLeaguesList, err := s.ListLeagues(ctx)
	if err != nil {
		return nil, nil, fmt.Errorf("error while getting current leagues list: %w", err)
	}

	remoteLeaguesList, err := s.fetcher.AllLeaguesList(ctx)
	if err != nil {
		return nil, nil, fmt.Errorf("error while getting leagues liist from remote: %w", err)
	}

	leaguesToDelete = findMissedLeaguesInFirstList(remoteLeaguesList, currentLeaguesList)
	leaguesToDownload = findMissedLeaguesInFirstList(currentLeaguesList, remoteLeaguesList)

	return leaguesToDownload, leaguesToDelete, nil
}

func findMissedLeaguesInFirstList(firstList []*models.League, secondList []*models.League) []*models.League {
	leaguesMap := make(map[string]struct{})
	for _, v := range firstList {
		leaguesMap[v.ID] = struct{}{}
	}

	var missed []*models.League
	for _, v := range secondList {
		_, ok := leaguesMap[v.ID]
		if !ok {
			missed = append(missed, v)
		}
	}

	return missed
}
