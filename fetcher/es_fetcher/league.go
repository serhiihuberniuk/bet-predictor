package fetcher

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/serhiihuberniuk/bet-predictor/fetcher/es_fetcher/es_models"
	"github.com/serhiihuberniuk/bet-predictor/models"
)

func (f *ESFetcher) AllLeaguesList(ctx context.Context) ([]*models.League, error) {

	remoteLeagues := make([]esmodels.League, 0)
	page := 10

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
	for _, v := range remoteLeagues {
		if len(v.Expand.CurrentSeason) == 0 {
			continue
		}

		league, err := models.NewLeague(models.CreateLeaguePayload{
			Name:            v.Name,
			Country:         v.CountryName,
			CurrentSeasonID: v.Expand.CurrentSeason[0].ID,
		})
		if err != nil {
			return nil, fmt.Errorf("error while creating league: %w", err)
		}

		leagues = append(leagues, league)
	}

	return leagues, nil
}

func (f *ESFetcher) getPageOfLeagues(ctx context.Context, page int) (*esmodels.Leagues, error) {
	req, err := http.NewRequest(http.MethodGet, "https://football.elenasport.io/v2/leagues?expand=current_season&page="+strconv.Itoa(page), nil)
	if err != nil {
		return nil, fmt.Errorf("error while creating request:  %w", err)
	}

	var leaguesResponse esmodels.Leagues

	decodeFn := func(reader io.Reader) error {
		if err := json.NewDecoder(reader).Decode(&leaguesResponse); err != nil {
			return fmt.Errorf("error while decoding response: %w", err)
		}

		return nil
	}

	if err = f.doRequest(ctx, req, decodeFn); err != nil {
		return nil, fmt.Errorf("error while doing request: %w", err)
	}

	return &leaguesResponse, nil
}
