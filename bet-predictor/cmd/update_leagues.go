package cmd

import (
	"fmt"

	"github.com/serhiihuberniuk/bet-predictor/models"
	"github.com/serhiihuberniuk/bet-predictor/scanner"
	"github.com/spf13/cobra"
)

var updateLeagueCmd = &cobra.Command{
	Use:   "update:leagues",
	Short: "Sync leagues information in db",
	Long:  "Long description of update leagues",
	RunE: func(cmd *cobra.Command, args []string) error {
		commandContext, err := commandInit()
		defer commandContext.DisconnectFn()
		defer commandContext.CancelFn()
		if err != nil {
			return fmt.Errorf("error while initialisation command: %w", err)
		}

		leaguesToDownload, leaguesToDelete, err := commandContext.Service.CompareLeaguesLists(commandContext.Ctx)
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

		Delete:

			for {
				answer := scanner.ScanWithMessage("Do you want to delete leagues (press y/n): ")
				switch answer {
				case "n":
					fmt.Println("leagues will not be deleted")

					break Delete

				case "y":
					for _, v := range leaguesToDelete {
						if err := commandContext.Service.DeleteLeague(commandContext.Ctx, v.ID); err != nil {
							return fmt.Errorf("error while deleting leagues: %w", err)
						}

						fmt.Println("leagues is deleted")

					}

					break Delete

				default:
					fmt.Println("unknown command: " + answer)
				}
			}
		}

		if len(leaguesToDownload) != 0 {
			fmt.Println("leagues to be downloaded: ")
			for _, v := range leaguesToDownload {
				fmt.Println(v.Country, v.Name)
			}

			fmt.Println()

		Download:
			for {
				answer := scanner.ScanWithMessage("Do you want to download leagues (press y/n): ")

				switch answer {
				case "n":
					fmt.Println("leagues will not be downloaded")

					break Download

				case "y":

					for _, v := range leaguesToDownload {
						if _, err := commandContext.Service.CreateLeague(commandContext.Ctx, models.CreateLeaguePayload{
							Name:            v.Name,
							Country:         v.Country,
							CurrentSeasonID: v.CurrentSeasonID,
						}); err != nil {
							return fmt.Errorf("error while creating leagues: %w", err)
						}
					}

					fmt.Println("leagues is downloaded")

					break Download

				default:
					fmt.Println("unknown command: " + answer)
				}
			}
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(updateLeagueCmd)
}
