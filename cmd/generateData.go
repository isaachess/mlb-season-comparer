/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"encoding/csv"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"time"

	"github.com/spf13/cobra"
)

// generateDataCmd represents the generateData command
var generateDataCmd = &cobra.Command{
	Use:   "generateData",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		const (
			numTeams    = 30
			numSeasons  = 150
			numMinGames = 150
			numMaxGames = 162
		)

		outFilePath, err := cmd.Flags().GetString("out-file")
		if err != nil {
			return err
		}

		f, err := os.Create(outFilePath)
		if err != nil {
			return err
		}
		defer f.Close()

		csvWriter := csv.NewWriter(f)
		headers := []string{"Year", "Team"}
		for i := 0; i < numMaxGames; i++ {
			headers = append(headers, fmt.Sprintf("Game %d", i+1))
		}
		if err := csvWriter.Write(headers); err != nil {
			return err
		}

		rand.Seed(time.Now().Unix())

		diffGames := numMaxGames - numMinGames

		for season := 1; season <= numSeasons; season++ {
			for team := 1; team <= numTeams; team++ {
				modifier := rand.Intn(diffGames + 1)
				games := numMinGames + modifier
				data := make([]string, 2, games+2)
				data[0] = strconv.Itoa(season)
				data[1] = strconv.Itoa(team)
				for game := 1; game <= numMaxGames; game++ {
					data = append(data, getGameResult(game, games))
				}
				if err := csvWriter.Write(data); err != nil {
					return fmt.Errorf("failed to write data for season %d and team %d", season, team)
				}
			}
		}

		csvWriter.Flush()

		return csvWriter.Error()
	},
}

func init() {
	rootCmd.AddCommand(generateDataCmd)
	generateDataCmd.Flags().String("out-file", "", "path to output CSV")
	generateDataCmd.MarkFlagRequired("out-file")
}

func getGameResult(game, games int) string {
	if game > games {
		return ""
	}
	win := rand.Int() % 2
	if win == 1 {
		return "W"
	}
	return "L"
}
