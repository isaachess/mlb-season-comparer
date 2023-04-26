/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"sort"
	"time"

	"github.com/spf13/cobra"
)

type teamStreaks struct {
	streaks     map[string]teamStreak
	bestStreaks map[string]teamStreak
}

func newTeamStreaks() *teamStreaks {
	return &teamStreaks{
		streaks:     make(map[string]teamStreak),
		bestStreaks: make(map[string]teamStreak),
	}
}

func (ts *teamStreaks) AddResult(game TeamGame) {
	currentStreak := ts.streaks[game.Franchise]
	newWins := currentStreak.wins
	newLosses := currentStreak.losses
	if game.Result == Win {
		newWins++
	}
	if game.Result == Loss {
		newLosses++
	}
	if currentStreak.games > 0 && newLosses >= newWins {
		oldBest := ts.bestStreaks[game.Franchise]
		if oldBest.games <= currentStreak.games {
			ts.bestStreaks[game.Franchise] = currentStreak
		}
		ts.streaks[game.Franchise] = teamStreak{}
		return
	}
	if currentStreak.games == 0 && game.Result == Loss {
		return
	}
	currentStreak.games++
	currentStreak.wins = newWins
	currentStreak.losses = newLosses
	currentStreak.end = game.Date
	currentStreak.endGame = game.TeamGameNumber
	if currentStreak.startGame == 0 && game.Result == Win {
		currentStreak.start = game.Date
		currentStreak.startGame = game.TeamGameNumber
	}
	ts.streaks[game.Franchise] = currentStreak
}

func (ts *teamStreaks) Flush() {
	for franchise, streak := range ts.streaks {
		oldBest := ts.bestStreaks[franchise]
		if oldBest.games <= streak.games {
			ts.bestStreaks[franchise] = streak
		}
	}
}

type teamStreak struct {
	franchise           string
	games, wins, losses int
	start               time.Time
	startGame           int
	end                 time.Time
	endGame             int
}

// longestOver500Cmd represents the longestOver500 command
var longestOver500Cmd = &cobra.Command{
	Use:   "longestOver500",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		teamsBySeason, err := GetTeamsBySeason("cmd/rs_data")
		if err != nil {
			return err
		}
		bss := teamsBySeason.BySortedSeason()
		streaks := newTeamStreaks()
		for _, seasons := range bss {
			for _, season := range seasons {
				for _, game := range season.Games {
					streaks.AddResult(game)
				}
			}
		}
		streaks.Flush()
		var allStreaks []teamStreak
		for franchise, streak := range streaks.bestStreaks {
			streak.franchise = franchise
			allStreaks = append(allStreaks, streak)
		}
		sort.Slice(allStreaks, func(i, j int) bool {
			return allStreaks[i].games < allStreaks[j].games
		})
		for _, streak := range allStreaks {
			fmt.Printf("Franchise: %s, Games: %d, Wins: %d, Losses: %d, start: %s, startGame: %d, end: %s, endGame: %d\n", streak.franchise, streak.games, streak.wins, streak.losses, streak.start.String(), streak.startGame, streak.end.String(), streak.endGame)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(longestOver500Cmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// longestOver500Cmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// longestOver500Cmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
