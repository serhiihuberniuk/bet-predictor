package cmd

import (
	"context"
	"fmt"
	"sync"

	"github.com/serhiihuberniuk/bet-predictor/models"
	"github.com/serhiihuberniuk/bet-predictor/scanner"
	"github.com/serhiihuberniuk/bet-predictor/service"
	"github.com/spf13/cobra"
)

var updateLeagueCmd = &cobra.Command{
	Use:   "leagues show",
	Short: "Sync leagues information in db",
	Long:  "Long description of update leagues",
	RunE: func(cmd *cobra.Command, args []string) error {
		wg := &sync.WaitGroup{}

		if err := updateLeague(wg); err != nil {
			wg.Wait()
			return fmt.Errorf("error while updating leagues: %w", err)
		}

		wg.Wait()
		return nil
	},
}

func deleteLeagues(ctx context.Context, service *service.Service, leaguesToDelete []*models.League) error {
	for {
		answer := scanner.ScanWithMessage("Do you want to delete leagues (press y/n): ")
		switch answer {
		case "n":
			fmt.Println("leagues will not be deleted")

			return nil

		case "y":
			for _, v := range leaguesToDelete {
				if err := service.DeleteLeague(ctx, v.ID); err != nil {
					return fmt.Errorf("error while deleting leagues: %w", err)
				}
			}

			fmt.Println("leagues is deleted")

			return nil

		default:
			fmt.Println("unknown command: " + answer)
		}
	}
}

func downloadLeagues(ctx context.Context, service *service.Service, leaguesToDownload []*models.League) error {
	for {
		answer := scanner.ScanWithMessage("Do you want to download leagues (press y/n): ")

		switch answer {
		case "n":
			fmt.Println("leagues will not be downloaded")

			return nil

		case "y":

			for _, v := range leaguesToDownload {
				if _, err := service.CreateLeague(ctx, models.CreateLeaguePayload{
					Name:            v.Name,
					Country:         v.Country,
					CurrentSeasonID: v.ESCurrentSeasonID,
				}); err != nil {
					return fmt.Errorf("error while creating leagues: %w", err)
				}
			}

			fmt.Println("leagues is downloaded")

			return nil

		default:
			fmt.Println("unknown command: " + answer)
		}
	}
}

func updateLeague(wg *sync.WaitGroup) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	s, err := commandInit(ctx, wg)
	if err != nil {
		return fmt.Errorf("error while initialisation command: %w", err)
	}

	leaguesToDownload, leaguesToDelete, err := s.CompareLeaguesLists(ctx)
	if err != nil {
		return fmt.Errorf("error while compering list of leagues: %w", err)
	}

	if len(leaguesToDownload) == 0 && len(leaguesToDelete) == 0 {
		fmt.Println("leagues in db is already up-to-dated")

		return nil
	}

	if len(leaguesToDelete) != 0 {
		fmt.Println("leagues to be deleted: ")
		for _, v := range leaguesToDelete {
			fmt.Println(v.Country, v.Name)
		}

		fmt.Println()

		if err = deleteLeagues(ctx, s, leaguesToDelete); err != nil {
			return fmt.Errorf("error while deleting leagues :%w", err)
		}
	}

	if len(leaguesToDownload) != 0 {
		fmt.Println("leagues to be downloaded: ")
		for _, v := range leaguesToDownload {
			fmt.Println(v.Country, v.Name)
		}

		fmt.Println()

		if err = downloadLeagues(ctx, s, leaguesToDownload); err != nil {
			return fmt.Errorf("error while downloading leagues: %w", err)
		}
	}

	return nil
}
