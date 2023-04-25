/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"encoding/csv"
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/spf13/cobra"
	"golang.org/x/sync/errgroup"
)

type seasonDetails struct {
	season    string
	team      string
	length    int
	gameStart int
	gameEnd   int
}

type gameCombos struct {
	mu     sync.Mutex
	combos map[string][]seasonDetails

	// matches is a set of hashes in combos that have matches
	matches map[string]struct{}
}

func newGameCombos() *gameCombos {
	return &gameCombos{
		combos:  map[string][]seasonDetails{},
		matches: map[string]struct{}{},
	}
}

func (gc *gameCombos) Add(key string, val seasonDetails) {
	gc.mu.Lock()
	defer gc.mu.Unlock()
	details, ok := gc.combos[key]
	details = append(details, val)
	gc.combos[key] = details
	if ok {
		gc.matches[key] = struct{}{}
	}
}

// compareCmd represents the compare command
var compareCmd = &cobra.Command{
	Use:   "compare",
	Short: "Run the comparison",
	Long: `Using the input CSV file provided, compare will run season W/L streak comparisons for your dataset. It will print out matches at the end.

Inputs:

in-file: The path the CSV containing the data
min-game-window: The lower bound of game-streak to look for.
max-game-window: The upper bound of game-streak to look for.

For example, if you the min is 30 and the max is 35, the script will find all instances where two seasons matched exactly for 30, 31, 32, 33, 34, and 35 games.
`,
	RunE: func(cmd *cobra.Command, args []string) error {
		inFilePath, err := cmd.Flags().GetString("in-file")
		if err != nil {
			return err
		}

		minGameWindow, err := cmd.Flags().GetInt("min-game-window")
		if err != nil {
			return err
		}
		maxGameWindow, err := cmd.Flags().GetInt("max-game-window")
		if err != nil {
			return err
		}

		f, err := os.Open(inFilePath)
		if err != nil {
			return err
		}
		defer f.Close()

		csvReader := csv.NewReader(f)
		records, err := csvReader.ReadAll()
		if err != nil {
			return err
		}

		combos := newGameCombos()

		if err := findMatches(records, combos, minGameWindow, maxGameWindow); err != nil {
			return err
		}

		for match := range combos.matches {
			details := combos.combos[match]
			fmt.Println("Match Found: ", match)
			for _, detail := range details {
				fmt.Printf("%+v\n", detail)
			}
		}

		return nil
	},
}

func findMatches(records [][]string, combos *gameCombos, minGameWindow, maxGameWindow int) error {
	var eg errgroup.Group

	for i, record := range records {
		if i == 0 {
			// Header row, skip
			continue
		}
		r := record

		// eg.Go(func() error {
		calculateHashes(r, combos, minGameWindow, maxGameWindow)
		// return nil
		// })
	}

	return eg.Wait()
}

func calculateHashes(record []string, combos *gameCombos, minGameWindow, maxGameWindow int) {
	season := record[0]
	team := record[1]
	gameOffset := 2
	maxGames := len(record) - gameOffset
	for gameWindow := minGameWindow; gameWindow <= maxGameWindow; gameWindow++ {
	A:
		for i := 0; i < maxGames; i++ {
			end := i + gameWindow + gameOffset
			if len(record) < end {
				continue A
			}
			results := record[i+gameOffset : end]
			if results[len(results)-1] == "" {
				continue A
			}
			combined := strings.Join(results, "")
			combos.Add(combined, seasonDetails{
				team:      team,
				season:    season,
				length:    gameWindow,
				gameStart: i + 1,
				gameEnd:   i + gameWindow,
			})
		}
	}
}

func init() {
	rootCmd.AddCommand(compareCmd)

	compareCmd.Flags().String("in-file", "", "path to input CSV")
	compareCmd.Flags().Int("min-game-window", 0, "lower bound of game window to compare")
	compareCmd.Flags().Int("max-game-window", 0, "upper bound of game window to compare")
	compareCmd.MarkFlagRequired("in-file")
	compareCmd.MarkFlagRequired("min-game-window")
	compareCmd.MarkFlagRequired("max-game-window")
}
