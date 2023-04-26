package cmd

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestTeamStreaks(t *testing.T) {
	tests := []struct {
		name         string
		games        []TeamGame
		expectedBest map[string]teamStreak
	}{
		{
			name: "hello",
			expectedBest: map[string]teamStreak{
				"F1": {
					games:     3,
					wins:      2,
					losses:    1,
					start:     time.Date(2020, time.June, 1, 0, 0, 0, 0, time.UTC),
					startGame: 1,
					end:       time.Date(2020, time.June, 3, 0, 0, 0, 0, time.UTC),
					endGame:   3,
				},
				"F2": {
					games:     5,
					wins:      3,
					losses:    2,
					start:     time.Date(2020, time.June, 4, 0, 0, 0, 0, time.UTC),
					startGame: 5,
					end:       time.Date(2021, time.June, 2, 0, 0, 0, 0, time.UTC),
					endGame:   1,
				},
			},
			games: []TeamGame{
				{
					Date:           time.Date(2020, time.June, 1, 0, 0, 0, 0, time.UTC),
					Team:           "T1",
					Franchise:      "F1",
					TeamGameNumber: 1,
					Result:         Win,
				},
				{
					Date:           time.Date(2020, time.June, 2, 0, 0, 0, 0, time.UTC),
					Team:           "T1",
					Franchise:      "F1",
					TeamGameNumber: 2,
					Result:         Win,
				},
				{
					Date:           time.Date(2020, time.June, 3, 0, 0, 0, 0, time.UTC),
					Team:           "T1",
					Franchise:      "F1",
					TeamGameNumber: 3,
					Result:         Loss,
				},
				{
					Date:           time.Date(2020, time.June, 4, 0, 0, 0, 0, time.UTC),
					Team:           "T1",
					Franchise:      "F1",
					TeamGameNumber: 4,
					Result:         Loss,
				},
				{
					Date:           time.Date(2020, time.June, 5, 0, 0, 0, 0, time.UTC),
					Team:           "T1",
					Franchise:      "F1",
					TeamGameNumber: 5,
					Result:         Win,
				},
				{
					Date:           time.Date(2020, time.June, 2, 0, 0, 0, 0, time.UTC),
					Team:           "T2",
					Franchise:      "F2",
					TeamGameNumber: 1,
					Result:         Loss,
				},
				{
					Date:           time.Date(2020, time.June, 2, 0, 0, 0, 0, time.UTC),
					Team:           "T2",
					Franchise:      "F2",
					TeamGameNumber: 2,
					Result:         Loss,
				},
				{
					Date:           time.Date(2020, time.June, 2, 0, 0, 0, 0, time.UTC),
					Team:           "T3",
					Franchise:      "F2",
					TeamGameNumber: 3,
					Result:         Win,
				},
				{
					Date:           time.Date(2020, time.June, 3, 0, 0, 0, 0, time.UTC),
					Team:           "T2",
					Franchise:      "F2",
					TeamGameNumber: 4,
					Result:         Loss,
				},
				{
					Date:           time.Date(2020, time.June, 4, 0, 0, 0, 0, time.UTC),
					Team:           "T4",
					Franchise:      "F2",
					TeamGameNumber: 5,
					Result:         Win,
				},
				{
					Date:           time.Date(2020, time.June, 5, 0, 0, 0, 0, time.UTC),
					Team:           "T2",
					Franchise:      "F2",
					TeamGameNumber: 6,
					Result:         Win,
				},
				{
					Date:           time.Date(2020, time.June, 6, 0, 0, 0, 0, time.UTC),
					Team:           "T2",
					Franchise:      "F2",
					TeamGameNumber: 7,
					Result:         Loss,
				},
				{
					Date:           time.Date(2020, time.June, 7, 0, 0, 0, 0, time.UTC),
					Team:           "T2",
					Franchise:      "F2",
					TeamGameNumber: 8,
					Result:         Win,
				},
				{
					Date:           time.Date(2021, time.June, 2, 0, 0, 0, 0, time.UTC),
					Team:           "T2",
					Franchise:      "F2",
					TeamGameNumber: 1,
					Result:         Loss,
				},
				{
					Date:           time.Date(2021, time.June, 3, 0, 0, 0, 0, time.UTC),
					Team:           "T2",
					Franchise:      "F2",
					TeamGameNumber: 2,
					Result:         Loss,
				},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			streaks := newTeamStreaks()
			for _, game := range test.games {
				streaks.AddResult(game)
			}
			streaks.Flush()
			assert.Equal(t, test.expectedBest, streaks.bestStreaks)
		})
	}
}
