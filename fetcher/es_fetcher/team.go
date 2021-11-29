package fetcher

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	esmodels "github.com/serhiihuberniuk/bet-predictor/fetcher/es_fetcher/es_models"
	"github.com/serhiihuberniuk/bet-predictor/models"
)

func (f *ESFetcher) GetTeamsBySeasonID(ctx context.Context, season string) ([]*models.Team, error) {
	page := 1
	remoteTeams := make([]esmodels.Team, 0)

	for {
		teamsResponse, err := f.getPageOfTeamsBySeason(ctx, season, page)
		if err != nil {
			return nil, fmt.Errorf("error while getting teams from page %v: %w", page, err)
		}

		remoteTeams = append(remoteTeams, teamsResponse.Data...)
		if !teamsResponse.Pagination.HasNextPage {
			break
		}
	}

	teams, err := createTeamsListFromResponse(remoteTeams)
	if err != nil {
		return nil, fmt.Errorf("error while creating list of teams from response: %w", err)
	}

	return teams, nil
}

func createTeamsListFromResponse(remoteTeams []esmodels.Team) ([]*models.Team, error) {
	teams := make([]*models.Team, 0, len(remoteTeams))
	for _, v := range remoteTeams {
		if v.FullName == "" {
			continue
		}

		team, err := models.NewTeam(models.CreateTeamPayload{
			Name:    v.FullName,
			Country: v.Country,
		})
		if err != nil {
			return nil, fmt.Errorf("error while creating new team: %w", err)
		}

		teams = append(teams, team)
	}
	return teams, nil
}

func (f *ESFetcher) getPageOfTeamsBySeason(ctx context.Context, season string, page int) (*esmodels.Teams, error) {
	req, err := http.NewRequest(http.MethodGet,
		fmt.Sprintf(hostNameES+"/v2/seasons/%v/teams?page=%v", season, page), nil)
	if err != nil {
		return nil, fmt.Errorf("error while creating request: %w", err)
	}

	var teams esmodels.Teams
	decodeFn := func(reader io.Reader) error {
		if err = json.NewDecoder(reader).Decode(&teams); err != nil {
			return fmt.Errorf("error while decoding from json: %w", err)
		}

		return nil
	}

	if err = f.doRequest(ctx, req, decodeFn); err != nil {
		return nil, fmt.Errorf("error while doing request: %w", err)
	}

	return &teams, nil
}
