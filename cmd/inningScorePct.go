/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"sort"
	"sync"

	"github.com/spf13/cobra"
)

type inningOutscorePerSeason struct {
	season     int
	weirdGames int
	totalGames int
}

// inningScorePctCmd represents the inningScorePct command
var inningScorePctCmd = &cobra.Command{
	Use:   "inningScorePct",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		var mu sync.Mutex
		allGames := map[int]inningOutscorePerSeason{}
		err := ByRetrosheetGame("cmd/rs_data", func(game *RetrosheetGame) error {
			homeLineScore, err := game.LineScoreProcessed(game.HomeLineScore)
			if err != nil {
				return err
			}
			visitingLineScore, err := game.LineScoreProcessed(game.VisitingLineScore)
			if err != nil {
				return err
			}
			mu.Lock()
			season := allGames[game.Date.Year()]
			season.totalGames++
			if isWeirdGame(homeLineScore, game.VisitingScore) || isWeirdGame(visitingLineScore, game.HomeScore) {
				fmt.Println("It's weird!", game.HomeLineScore, game.VisitingLineScore, game.HomeScore, game.VisitingScore)
				season.weirdGames++
			}
			allGames[game.Date.Year()] = season
			mu.Unlock()
			return nil
		})
		if err != nil {
			return err
		}

		var seasonsList []inningOutscorePerSeason
		for year, details := range allGames {
			details.season = year
			seasonsList = append(seasonsList, details)
		}
		sort.Slice(seasonsList, func(i, j int) bool {
			return seasonsList[i].season < seasonsList[j].season
		})
		for _, season := range seasonsList {
			fmt.Printf("%d\t%d\n", season.season, season.weirdGames*100/season.totalGames)
		}
		return nil
	},
}

func isWeirdGame(lineScore []int, oppoScore int) bool {
	for _, score := range lineScore {
		if score > oppoScore {
			return true
		}
	}
	return false
}

func init() {
	rootCmd.AddCommand(inningScorePctCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// inningScorePctCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// inningScorePctCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
