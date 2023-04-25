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
	Short: "Generate test data.",
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

		teams := []string{
			"ARI",
			"ATL",
			"BAL",
			"BOS",
			"CHC",
			"CHW",
			"CIN",
			"CLE",
			"COL",
			"DET",
			"FLA",
			"HOU",
			"KAN",
			"LAA",
			"LAD",
			"MIL",
			"MIN",
			"NYM",
			"NYY",
			"OAK",
			"PHI",
			"PIT",
			"SD",
			"SF",
			"SEA",
			"STL",
			"TB",
			"TEX",
			"TOR",
			"WAS",
		}

		for season := 2022 - numSeasons; season <= 2022; season++ {
			for _, team := range teams {
				modifier := rand.Intn(diffGames + 1)
				games := numMinGames + modifier
				data := make([]string, 2, games+2)
				data[0] = strconv.Itoa(season)
				data[1] = team
				for game := 1; game <= numMaxGames; game++ {
					data = append(data, getGameResult(game, games))
				}
				if err := csvWriter.Write(data); err != nil {
					return fmt.Errorf("failed to write data for season %d and team %s", season, team)
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
