package cmd

import (
	"bufio"
	"errors"
	"fmt"
	"os"

	"github.com/serhiihuberniuk/bet-predictor/models"
	"github.com/spf13/cobra"
)

var updateLeagueCmd = &cobra.Command{
	Use:   "update leagues",
	Short: "Sync leagues information in db",
	Long:  "Long description of update leagues",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, f, s, disconnectDbFunc, err := commandInit()
		if err != nil {
			if disconnectDbFunc != nil {
				disconnectDbFunc()
			}
			return fmt.Errorf("error while initialisation command: %w", err)
		}
		defer disconnectDbFunc()

		remoteLeagues, err := f.AllLeaguesList(ctx)
		if err != nil {
			return fmt.Errorf("error while getting leagues from remote: %w", err)
		}

		leaguesToDownload, leaguesToDelete, err := s.CompareLeaguesLists(ctx, remoteLeagues)
		if err != nil {
			return fmt.Errorf("error while compering list of leagues: %w", err)
		}

		if len(leaguesToDownload) == 0 && len(leaguesToDelete) == 0 {
			fmt.Println("leagues in db is already up-to-dated")

			return nil
		}

		if len(leaguesToDownload) != 0 {
			fmt.Println("leagues to be downloaded: ")
			for _, v := range leaguesToDownload {
				fmt.Println(v.Country, v.Name)
			}

			fmt.Println()
		}

		if len(leaguesToDelete) != 0 {
			fmt.Println("leagues to be deleted: ")
			for _, v := range leaguesToDelete {
				fmt.Println(v.Country, v.Name)
			}

			fmt.Println()
		}

		fmt.Print("Update leagues list? (press y/n): ")
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()

		switch scanner.Text() {
		case "n":
			fmt.Println("leagues will not be updated")
		case "y":
			for _, v := range leaguesToDelete {
				if err := s.DeleteLeague(ctx, v.ID); err != nil {
					return fmt.Errorf("error while deleting leagues: %w", err)
				}
			}

			for _, v := range leaguesToDownload {
				if _, err := s.CreateLeague(ctx, &models.CreateLeaguePayload{
					Name:    v.Name,
					Country: v.Country,
				}); err != nil {
					return fmt.Errorf("error while creating leagues: %w", err)
				}
			}

			fmt.Println("leagues is updated")

		default:
			return errors.New("unknown command: " + scanner.Text())
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(updateLeagueCmd)
}
