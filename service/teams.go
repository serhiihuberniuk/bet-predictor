package service

import (
	"context"
	"fmt"

	"github.com/serhiihuberniuk/bet-predictor/models"
)

func (s *Service) CreateTeam(ctx context.Context, payload models.CreateTeamPayload) (string, error) {
	t, err := models.NewTeam(payload)
	if err != nil {
		return "", fmt.Errorf("error while creating team: %w", err)
	}

	if err = s.repo.CreateTeam(ctx, t); err != nil {
		return "", fmt.Errorf("error while creating team in repository layer: %w", err)
	}

	return t.ID, nil
}

func (s *Service) DeleteTeam(ctx context.Context, id string) error {
	if err := s.repo.DeleteTeam(ctx, id); err != nil {
		return fmt.Errorf("error while delete teams from reposstory layer: %w", err)
	}

	return nil
}

func (s *Service) ListTeams(ctx context.Context) ([]*models.Team, error) {
	teams, err := s.repo.ListTeams(ctx)
	if err != nil {
		return nil, fmt.Errorf("error while getting list of teams from repository layer: %w", err)
	}

	return teams, nil
}

func (s *Service) CompareAllTeamsLists(ctx context.Context) (teamsToDownload []*models.Team, teamsToDelete []*models.Team, err error) {
	currentTeams, err := s.ListTeams(ctx)
	if err != nil {
		return nil, nil, fmt.Errorf("error while getting teams list: %w", err)
	}

	leagues, err := s.ListLeagues(ctx)
	if err != nil {
		return nil, nil, fmt.Errorf("error while getting leagues list: %w", err)
	}

	var remoteTeams []*models.Team
	for _, league := range leagues {

		teams, err := s.fetcher.GetTeamsBySeasonID(ctx, league.CurrentSeasonID)
		if err != nil {
			return nil, nil, fmt.Errorf("error while getting list of teams from remote: %w", err)
		}

		remoteTeams = append(remoteTeams, teams...)
	}

	teamsToDownload = findMissedTeamsInFirstList(currentTeams, remoteTeams)
	teamsToDelete = findMissedTeamsInFirstList(remoteTeams, currentTeams)

	return teamsToDownload, teamsToDelete, nil
}

func (s *Service) CompareTeamsListInLeague(ctx context.Context, league *models.League) (teamsToDownload []*models.Team, teamsToDelete []*models.Team, err error) {
	currentTeams, err := s.repo.ListTeams(ctx)
	if err != nil {
		return nil, nil, fmt.Errorf("error while getting teams list from db: %w", err)
	}

	remoteTeams, err := s.fetcher.GetTeamsBySeasonID(ctx, league.CurrentSeasonID)
	if err != nil {
		return nil, nil, fmt.Errorf("error while getting teams list from remote: %w", err)
	}

	teamsToDownload = findMissedTeamsInFirstList(currentTeams, remoteTeams)

	return teamsToDownload, nil, nil
}

func findMissedTeamsInFirstList(firstList []*models.Team, secondList []*models.Team) []*models.Team {
	teamsMap := make(map[string]struct{})
	for _, v := range firstList {
		teamsMap[v.ID] = struct{}{}
	}

	var missed []*models.Team
	duplicates := make(map[string]struct{})

	for _, v := range secondList {
		_, ok := teamsMap[v.ID]
		if !ok {
			_, ok = duplicates[v.ID]
			if !ok {
				missed = append(missed, v)
				duplicates[v.ID] = struct{}{}
			}
		}
	}

	return missed
}
