/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"sort"

	"github.com/spf13/cobra"
)

type SeasonSubset struct {
	Season     *Season
	Start, End int
}

// recordInSeasonCmd represents the recordInSeason command
var recordInSeasonCmd = &cobra.Command{
	Use:   "recordInSeason",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		wins, err := cmd.Flags().GetInt("wins")
		if err != nil {
			return err
		}
		losses, err := cmd.Flags().GetInt("losses")
		if err != nil {
			return err
		}
		since, err := cmd.Flags().GetInt("since")
		if err != nil {
			return err
		}

		teamsBySeason, err := GetTeamsBySeason("cmd/rs_data")
		if err != nil {
			return err
		}
		bss := teamsBySeason.BySortedSeason()
		var matchingSeasons []*SeasonSubset
		for _, seasons := range bss {
			for _, season := range seasons {
				if subset := seasonMatchesRecord(season, wins, losses, since); subset != nil {
					matchingSeasons = append(matchingSeasons, subset)
				}
			}
		}

		sort.Slice(matchingSeasons, func(i, j int) bool {
			return matchingSeasons[i].Season.GetSeasonRecord().Wins < matchingSeasons[j].Season.GetSeasonRecord().Wins
		})

		for _, subset := range matchingSeasons {
			fmt.Println(subset.Season.Franchise, subset.Season.Year, "Start", subset.Start, "End", subset.End, "Record", subset.Season.GetSeasonRecord().String())
		}

		return nil
	},
}

func seasonMatchesRecord(season *Season, wins, losses, since int) *SeasonSubset {
	if season.Year < since {
		return nil
	}
	gameWindow := wins + losses
	for i := 0; i < len(season.Games); i++ {
		end := i + gameWindow
		if len(season.Games) < end {
			return nil
		}
		var (
			winsInResults   int
			lossesInResults int
		)
		results := season.Games[i:end]
		for _, result := range results {
			if result.Result == Win {
				winsInResults++
			} else if result.Result == Loss {
				lossesInResults++
			}
		}
		if winsInResults == wins && lossesInResults == losses {
			return &SeasonSubset{
				Season: season,
				Start:  i + 1,
				End:    end,
			}
		}
	}
	return nil
}

func init() {
	rootCmd.AddCommand(recordInSeasonCmd)
	recordInSeasonCmd.Flags().Int("wins", 0, "wins within record")
	recordInSeasonCmd.Flags().Int("losses", 0, "losses within record")
	recordInSeasonCmd.Flags().Int("since", 0, "year to start tracking")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// recordInSeasonCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// recordInSeasonCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
