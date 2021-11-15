package cmd

import (
	"github.com/spf13/cobra"
)

var apiKeyFlag string
var cfgFileFlag string

var rootCmd = &cobra.Command{
	Use:   "bet-predictor",
	Short: "Predictor bets on football matches",
	Long:  `Long description of bet-predictor`,
}

func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	rootCmd.PersistentFlags().StringVar(&apiKeyFlag, "api-key", "", "Set api-key if it is different from default")
	rootCmd.PersistentFlags().StringVar(&cfgFileFlag, "config", "./../config.yaml", "Set path to config file")
}
