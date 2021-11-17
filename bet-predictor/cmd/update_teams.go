package cmd

import (
	"errors"
	"fmt"
	"strings"

	"github.com/serhiihuberniuk/bet-predictor/models"
	"github.com/serhiihuberniuk/bet-predictor/scanner"
	"github.com/spf13/cobra"
)

var leagueFlag string
var countryFlag string

var updateTeamsCmd = &cobra.Command{
	Use:   "update:teams",
	Short: "Sync teams in your DB",
	Long:  "Long description of update:teams",
	RunE: func(cmd *cobra.Command, args []string) error {
		commandContext, err := commandInit()
		defer commandContext.DisconnectFn()
		defer commandContext.CancelFn()
		if err != nil {
			return fmt.Errorf("error while initialisation command: %w", err)
		}

		teamsToDownload := make([]*models.Team, 0)
		teamsToDelete := make([]*models.Team, 0)

		if leagueFlag != "" {
			league, err := commandContext.Service.GetLeagueByCountryAndName(commandContext.Ctx,
				strings.Trim(countryFlag, " "), strings.Trim(leagueFlag, " "))
			if err != nil {
				if errors.Is(err, models.ErrNotFound) {
					fmt.Println("league you specified in flag --league does not exist in country you specified by flag --country")

					return nil
				}
				return fmt.Errorf("error while getting leagues specified by flags: %w", err)
			}

			teamsToDownload, teamsToDelete, err = commandContext.Service.CompareTeamsListInLeague(commandContext.Ctx, league)
			if err != nil {
				return fmt.Errorf("error while compering lists of teams: %w", err)
			}

		} else {
			teamsToDownload, teamsToDelete, err = commandContext.Service.CompareAllTeamsLists(commandContext.Ctx)
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

		Delete:

			for {
				answer := scanner.ScanWithMessage("Do you want to delete teams (press y/n): ")
				switch answer {
				case "n":
					fmt.Println("teams will not be deleted")

					break Delete

				case "y":
					for _, v := range teamsToDelete {
						if err := commandContext.Service.DeleteTeam(commandContext.Ctx, v.ID); err != nil {
							return fmt.Errorf("error while deleting teams: %w", err)
						}

						fmt.Println("teams is deleted")

					}

					break Delete

				default:
					fmt.Println("unknown command: " + answer)
				}
			}
		}

		if len(teamsToDownload) != 0 {
			fmt.Println("teams to be downloaded: ")
			for _, v := range teamsToDownload {
				fmt.Println(v.Country, v.Name)
			}

			fmt.Println()

		Download:
			for {
				answer := scanner.ScanWithMessage("Do you want to download teams (press y/n): ")

				switch answer {
				case "n":
					fmt.Println("teams will not be downloaded")

					break Download

				case "y":

					for _, v := range teamsToDownload {
						if _, err := commandContext.Service.CreateTeam(commandContext.Ctx, models.CreateTeamPayload{
							Name:    v.Name,
							Country: v.Country,
						}); err != nil {
							return fmt.Errorf("error while creating teams: %w", err)
						}
					}

					fmt.Println("teams is downloaded")

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
	rootCmd.AddCommand(updateTeamsCmd)
	updateTeamsCmd.PersistentFlags().StringVarP(&leagueFlag, "league", "l", "", "Set slug-name of the league in which you want to update the teams")
	updateTeamsCmd.PersistentFlags().StringVarP(&countryFlag, "country", "c", "", "Set slug-name of the country in which the league you want to update is located")
}
