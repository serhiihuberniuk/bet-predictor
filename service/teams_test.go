package service_test

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/serhiihuberniuk/bet-predictor/models"
	"github.com/serhiihuberniuk/bet-predictor/service"
	"github.com/stretchr/testify/assert"
)

func TestService_CompareAllTeamsLists(t *testing.T) {
	t.Parallel()

	type RepoMockBehavior func(r *Mockrepository)
	type FetcherMockBehavior func(f *Mockfetcher, seasonID string)

	team1, _ := models.NewTeam(models.CreateTeamPayload{
		Name:    "team1",
		Country: "country1",
	})
	team2, _ := models.NewTeam(models.CreateTeamPayload{
		Name:    "team2",
		Country: "country2",
	})
	league1, _ := models.NewLeague(models.CreateLeaguePayload{
		Name:            "league1",
		Country:         "league1",
		CurrentSeasonID: "0",
	})
	league2, _ := models.NewLeague(models.CreateLeaguePayload{
		Name:            "league2",
		Country:         "league2",
		CurrentSeasonID: "1",
	})

	inCtx := context.Background()

	testCases := []struct {
		name                    string
		RepoMockBehavior        RepoMockBehavior
		FetcherMockBehavior     FetcherMockBehavior
		teamsToDownloadExpected []*models.Team
		teamsToDeleteExpected   []*models.Team
		errMessage              string
	}{
		{
			name: "OK",
			RepoMockBehavior: func(r *Mockrepository) {
				r.EXPECT().ListTeams(inCtx).Return([]*models.Team{team1}, nil)
				r.EXPECT().ListLeagues(inCtx).Return([]*models.League{league1, league2}, nil).MaxTimes(2)

			},
			FetcherMockBehavior: func(f *Mockfetcher, seasonID string) {
				f.EXPECT().GetTeamsBySeasonID(inCtx, seasonID).Return([]*models.Team{team1, team2}, nil).AnyTimes()
			},
			teamsToDeleteExpected:   nil,
			teamsToDownloadExpected: []*models.Team{team2},
			errMessage:              "",
		},
		{
			name: "Error while getting current teams",
			RepoMockBehavior: func(r *Mockrepository) {
				r.EXPECT().ListTeams(inCtx).Return(nil, errors.New("error"))
				r.EXPECT().ListLeagues(inCtx).Return(nil, nil).MaxTimes(2)
			},
			FetcherMockBehavior: func(f *Mockfetcher, seasonID string) {

			},
			errMessage: "error while getting teams list",
		},
		{
			name: "Error while getting list leagues",
			RepoMockBehavior: func(r *Mockrepository) {
				r.EXPECT().ListTeams(inCtx).Return([]*models.Team{team1}, nil)
				r.EXPECT().ListLeagues(inCtx).Return(nil, errors.New("error")).MaxTimes(2)
			},
			FetcherMockBehavior: func(f *Mockfetcher, seasonID string) {

			},
			errMessage: "error while getting leagues list",
		},
		{
			name: "Error while getting teams from remote",
			RepoMockBehavior: func(r *Mockrepository) {
				r.EXPECT().ListTeams(inCtx).Return([]*models.Team{team1}, nil)
				r.EXPECT().ListLeagues(inCtx).Return([]*models.League{league1, league2}, nil).MaxTimes(2)

			},
			FetcherMockBehavior: func(f *Mockfetcher, seasonID string) {
				f.EXPECT().GetTeamsBySeasonID(inCtx, seasonID).Return(nil, errors.New("error")).AnyTimes()
			},
			errMessage: "error while getting list of teams from remote",
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			repoMock := NewMockrepository(ctrl)
			fetcherMock := NewMockfetcher(ctrl)

			s := service.New(repoMock, fetcherMock)

			tc.RepoMockBehavior(repoMock)
			leagues, _ := s.ListLeagues(inCtx)
			for _, l := range leagues {
				tc.FetcherMockBehavior(fetcherMock, l.ESCurrentSeasonID)
			}

			teamsToDownload, teamsToDelete, err := s.CompareAllTeamsLists(inCtx)
			if tc.errMessage != "" {
				assert.Contains(t, err.Error(), tc.errMessage)

				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tc.teamsToDownloadExpected, teamsToDownload)
			assert.Equal(t, tc.teamsToDeleteExpected, teamsToDelete)

		})
	}
}
