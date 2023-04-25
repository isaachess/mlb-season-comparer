/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
)

type season struct {
	year  int
	games map[string][]game
}

type game struct {
	team       string
	result     string
	gameNumber int
}

// transformCmd represents the transform command
var transformCmd = &cobra.Command{
	Use:   "transform",
	Short: "Transform Retrosheet data",
	RunE: func(cmd *cobra.Command, args []string) error {
		inDirPath, err := cmd.Flags().GetString("in-dir")
		if err != nil {
			return err
		}
		outFilePath, err := cmd.Flags().GetString("out-file")
		if err != nil {
			return err
		}

		files, err := os.ReadDir(inDirPath)
		if err != nil {
			return err
		}

		var seasons []season

		for _, file := range files {
			if file.IsDir() {
				continue
			}
			f, err := os.Open(filepath.Join(inDirPath, file.Name()))
			if err != nil {
				return err
			}
			defer f.Close()

			csvReader := csv.NewReader(f)
			records, err := csvReader.ReadAll()
			if err != nil {
				return err
			}

			seasonNumber, err := getSeasonFromName(file.Name())
			if err != nil {
				fmt.Printf("Skipping file %s due to error %s\n", file.Name(), err.Error())
				continue
			}
			games, err := processGames(records)
			if err != nil {
				return err
			}
			seasons = append(seasons, season{
				year:  seasonNumber,
				games: games,
			})
		}

		outFile, err := os.Create(outFilePath)
		if err != nil {
			return err
		}
		defer outFile.Close()

		csvWriter := csv.NewWriter(outFile)

		maxGames := getMaxGames(seasons)

		headers := []string{"Year", "Team"}
		for i := 0; i < maxGames; i++ {
			headers = append(headers, fmt.Sprintf("Game%d", i+1))
		}

		if err := csvWriter.Write(headers); err != nil {
			return err
		}

		for _, season := range seasons {
			for team, games := range season.games {
				vals := getCSV(season.year, team, games, maxGames)
				if err := csvWriter.Write(vals); err != nil {
					return err
				}
			}
		}
		csvWriter.Flush()
		return csvWriter.Error()
	},
}

func getCSV(year int, team string, games []game, maxGames int) []string {
	csv := []string{strconv.Itoa(year), team}
	for i := 0; i < maxGames; i++ {
		if i+1 > len(games) {
			csv = append(csv, "")
		} else {
			csv = append(csv, games[i].result)
		}
	}
	return csv
}

func getSeasonFromName(name string) (int, error) {
	return strconv.Atoi(strings.TrimSuffix(name, filepath.Ext(name))[2:])
}

func getMaxGames(seasons []season) int {
	var max int
	for _, season := range seasons {
		for _, games := range season.games {
			for _, game := range games {
				if game.gameNumber > max {
					max = game.gameNumber
				}
			}
		}
	}
	return max
}

func processGames(records [][]string) (map[string][]game, error) {
	allGames := make(map[string][]game)
	for _, record := range records {
		games, err := processGame(record)
		if err != nil {
			return nil, err
		}
		for _, game := range games {
			allGames[game.team] = append(allGames[game.team], game)
		}
	}
	for team, games := range allGames {
		sort.Slice(games, func(i, j int) bool {
			return games[i].gameNumber < games[j].gameNumber
		})
		allGames[team] = games
	}
	return allGames, nil
}

func processGame(record []string) ([]game, error) {
	visitingTeam := record[3]
	visitingGameNumber, err := strconv.Atoi(record[5])
	if err != nil {
		return nil, err
	}
	homeTeam := record[6]
	homeGameNumber, err := strconv.Atoi(record[8])
	if err != nil {
		return nil, err
	}
	visitingScore, err := strconv.Atoi(record[9])
	if err != nil {
		return nil, err
	}
	homeScore, err := strconv.Atoi(record[10])
	if err != nil {
		return nil, err
	}
	forfeitInfo := record[14]

	homeResult := calcResult(homeScore, visitingScore, true, forfeitInfo)
	visitingResult := calcResult(visitingScore, homeScore, false, forfeitInfo)

	return []game{
		{
			team:       homeTeam,
			result:     homeResult,
			gameNumber: homeGameNumber,
		},
		{
			team:       visitingTeam,
			result:     visitingResult,
			gameNumber: visitingGameNumber,
		},
	}, nil
}

func calcResult(teamScore, oppoScore int, isHome bool, forfeitInfo string) string {
	if forfeitInfo != "" {
		if forfeitInfo == "T" {
			return "T"
		}
		if isHome {
			if forfeitInfo == "H" {
				return "W"
			}
			if forfeitInfo == "V" {
				return "L"
			}
		} else {
			if forfeitInfo == "H" {
				return "L"
			}
			if forfeitInfo == "V" {
				return "W"
			}
		}
	}

	if teamScore == oppoScore {
		return "T"
	}
	if teamScore > oppoScore {
		return "W"
	}
	return "L"
}

func init() {
	rootCmd.AddCommand(transformCmd)
	transformCmd.Flags().String("in-dir", "", "path to directory containing retrosheet data")
	transformCmd.MarkFlagRequired("in-dir")
	transformCmd.Flags().String("out-file", "", "path to output CSV")
	transformCmd.MarkFlagRequired("out-file")
}
