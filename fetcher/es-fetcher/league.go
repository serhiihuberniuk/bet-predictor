package fetcher

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/serhiihuberniuk/bet-predictor/fetcher/es-fetcher/es-models"
	"github.com/serhiihuberniuk/bet-predictor/models"
)

func (f *ESFetcher) AllLeaguesList(ctx context.Context) ([]*models.League, error) {

	remoteLeagues := make([]esmodels.League, 0)
	page := 1

	for {
		leaguesResponse, err := f.getPageOfLeagues(ctx, page)
		if err != nil {
			return nil, fmt.Errorf("error while getting leagues from page %v: %w", page, err)
		}

		remoteLeagues = append(remoteLeagues, leaguesResponse.Data...)

		if !leaguesResponse.Pagination.HasNextPage {
			break
		}

		page++
	}

	var leagues []*models.League
	for _, league := range remoteLeagues {
		leagues = append(leagues, &models.League{
			Name:    league.Name,
			Country: league.CountryName,
		})
	}

	return leagues, nil
}

func (f *ESFetcher) getPageOfLeagues(ctx context.Context, page int) (*esmodels.Leagues, error) {
	req, err := http.NewRequest(http.MethodGet, "https://football.elenasport.io/v2/leagues?page="+strconv.Itoa(page), nil)
	if err != nil {
		return nil, fmt.Errorf("error while creating request:  %w", err)
	}

	req.Header.Add("Authorization", "Bearer "+f.token)

	resp, err := f.client.Do(req.WithContext(ctx))
	if err != nil {
		return nil, fmt.Errorf("error while sending GET request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("response with status: " + resp.Status)
	}

	var leaguesResponse esmodels.Leagues
	if err := json.NewDecoder(resp.Body).Decode(&leaguesResponse); err != nil {
		return nil, fmt.Errorf("error while decoding response: %w", err)
	}

	return &leaguesResponse, nil
}
