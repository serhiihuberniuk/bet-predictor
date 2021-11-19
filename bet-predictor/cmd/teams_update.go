package cmd

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/serhiihuberniuk/bet-predictor/models"
	"github.com/serhiihuberniuk/bet-predictor/scanner"
	"github.com/serhiihuberniuk/bet-predictor/service"
	"github.com/spf13/cobra"
)

var leagueSlag string

var updateTeamsCmd = &cobra.Command{
	Use:   "teams update",
	Short: "Sync teams in your DB",
	Long:  "Long description of update:teams",
	RunE: func(cmd *cobra.Command, args []string) error {
		wg := &sync.WaitGroup{}
		defer wg.Wait()

		if err := updateTeams(wg); err != nil {
			wg.Wait()

			return fmt.Errorf("error while updating teams: %w", err)
		}
		wg.Wait()

		return nil
	},
}

func deleteTeams(ctx context.Context, service *service.Service, teamsToDelete []*models.Team) error {
	for {
		answer := scanner.ScanWithMessage("Do you want to delete teams (press y/n): ")
		switch answer {
		case "n":
			fmt.Println("teams will not be deleted")

			return nil

		case "y":
			for _, v := range teamsToDelete {
				if err := service.DeleteTeam(ctx, v.ID); err != nil {
					return fmt.Errorf("error while deleting teams: %w", err)
				}
			}

			fmt.Println("teams is deleted")

			return nil

		default:
			fmt.Println("unknown command: " + answer)
		}
	}
}

func downloadTeams(ctx context.Context, service *service.Service, teamsToDownload []*models.Team) error {
	for {
		answer := scanner.ScanWithMessage("Do you want to download teams (press y/n): ")

		switch answer {
		case "n":
			fmt.Println("teams will not be downloaded")

			return nil

		case "y":

			for _, v := range teamsToDownload {
				if _, err := service.CreateTeam(ctx, models.CreateTeamPayload{
					Name:    v.Name,
					Country: v.Country,
				}); err != nil {
					return fmt.Errorf("error while creating teams: %w", err)
				}
			}

			fmt.Println("teams is downloaded")

			return nil

		default:
			fmt.Println("unknown command: " + answer)
		}
	}
}

func updateTeams(wg *sync.WaitGroup) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	s, err := commandInit(ctx, wg)
	if err != nil {
		return fmt.Errorf("error while initialisation command: %w", err)
	}

	teamsToDownload := make([]*models.Team, 0)
	teamsToDelete := make([]*models.Team, 0)

	if leagueSlag != "" {
		league, err := s.GetLeagueBySlug(ctx, leagueSlag)
		if err != nil {
			if errors.Is(err, models.ErrNotFound) {
				fmt.Println("league you specified in flag --league-slug does not exist in DB")

				return nil
			}
			return fmt.Errorf("error while getting leagues specified by flags: %w", err)
		}

		teamsToDownload, teamsToDelete, err = s.CompareTeamsListInLeague(ctx, league)
		if err != nil {
			return fmt.Errorf("error while compering lists of teams: %w", err)
		}

	} else {
		teamsToDownload, teamsToDelete, err = s.CompareAllTeamsLists(ctx)
		if err != nil {
			return fmt.Errorf("error while compering lists of teams: %w", err)
		}
	}

	if len(teamsToDownload) == 0 && len(teamsToDelete) == 0 {
		fmt.Println("teams in db is already is up-to-dated")

		return nil
	}

	if len(teamsToDelete) != 0 {
		fmt.Println("teams to be deleted: ")
		for _, v := range teamsToDelete {
			fmt.Println(v.Country, v.Name)
		}

		fmt.Println()

		if err := deleteTeams(ctx, s, teamsToDownload); err != nil {
			return fmt.Errorf("error while deleting teams: %w", err)
		}
	}

	if len(teamsToDownload) != 0 {
		fmt.Println("teams to be downloaded: ")
		for _, v := range teamsToDownload {
			fmt.Println(v.Country, v.Name)
		}

		fmt.Println()

		if err := downloadTeams(ctx, s, teamsToDownload); err != nil {
			return fmt.Errorf("error while downloading teams: %w", err)
		}
	}

	return nil
}
