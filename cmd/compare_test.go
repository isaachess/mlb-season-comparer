package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFindMatches(t *testing.T) {
	tests := []struct {
		name            string
		records         [][]string
		minGameWindow   int
		maxGameWindow   int
		expectedMatches map[string]struct{}
		expectedCombos  map[string][]seasonDetails
		checkCombos     bool
	}{
		{
			name:        "basic check",
			checkCombos: true,
			records: [][]string{
				{"Season", "Team", "Game1", "Game2", "Game3", "Game4", "Game5"},
				{"1", "1", "W", "W", "W", "W", "W"},
				{"1", "3", "L", "W", "L", "W", "L"},
				{"1", "4", "W", "L", "W", "W", "W"},
				{"1", "5", "L", "W", "W", "W", ""},
				{"1", "6", "L", "W", "L", "W", "L"},
				{"1", "7", "W", "W", "W", "W", "W"},
			},
			minGameWindow: 4,
			maxGameWindow: 5,
			expectedMatches: map[string]struct{}{
				"LWLW":  {},
				"LWLWL": {},
				"LWWW":  {},
				"WLWL":  {},
				"WWWW":  {},
				"WWWWW": {},
			},
			expectedCombos: map[string][]seasonDetails{
				"LWLW":  {{team: "3", season: "1", length: 4, gameStart: 1, gameEnd: 4}, {team: "6", season: "1", length: 4, gameStart: 1, gameEnd: 4}},
				"LWLWL": {{team: "3", season: "1", length: 5, gameStart: 1, gameEnd: 5}, {team: "6", season: "1", length: 5, gameStart: 1, gameEnd: 5}},
				"LWWW":  {{team: "4", season: "1", length: 4, gameStart: 2, gameEnd: 5}, {team: "5", season: "1", length: 4, gameStart: 1, gameEnd: 4}},
				"WLWL":  {{team: "3", season: "1", length: 4, gameStart: 2, gameEnd: 5}, {team: "6", season: "1", length: 4, gameStart: 2, gameEnd: 5}},
				"WLWW":  {{team: "4", season: "1", length: 4, gameStart: 1, gameEnd: 4}},
				"WLWWW": {{team: "4", season: "1", length: 5, gameStart: 1, gameEnd: 5}},
				"WWWW":  {{team: "1", season: "1", length: 4, gameStart: 1, gameEnd: 4}, {team: "1", season: "1", length: 4, gameStart: 2, gameEnd: 5}, {team: "7", season: "1", length: 4, gameStart: 1, gameEnd: 4}, {team: "7", season: "1", length: 4, gameStart: 2, gameEnd: 5}},
				"WWWWW": {{team: "1", season: "1", length: 5, gameStart: 1, gameEnd: 5}, {team: "7", season: "1", length: 5, gameStart: 1, gameEnd: 5}},
			},
		},
		{
			name: "big check",
			records: [][]string{
				{"Season", "Team", "Game", "Game", "Game", "Game", "Game", "Game", "Game", "Game", "Game", "Game", "Game", "Game", "Game", "Game", "Game", "Game", "Game", "Game", "Game", "Game", "Game", "Game", "Game", "Game", "Game", "Game", "Game", "Game", "Game", "Game"},
				{"1", "1", "W", "W", "W", "W", "W", "W", "W", "W", "W", "W", "W", "W", "W", "W", "W", "W", "W", "W", "W", "W", "W", "W", "W", "W", "W", "W", "W", "W", "W", "W"},
				{"1", "2", "W", "W", "W", "W", "W", "W", "W", "W", "W", "W", "W", "W", "W", "W", "W", "W", "W", "W", "W", "W", "W", "W", "W", "W", "W", "W", "W", "L", "W", "W"},
				{"1", "3", "W", "W", "W", "W", "W", "W", "W", "W", "W", "W", "W", "W", "W", "W", "W", "W", "W", "W", "W", "W", "W", "W", "L", "W", "W", "W", "W", "W", "W", "W"},
				{"1", "4", "W", "W", "W", "W", "W", "W", "W", "W", "W", "W", "W", "W", "W", "W", "W", "W", "W", "W", "W", "W", "W", "W", "W", "W", "W", "W", "W", "L", "W", "W"},
			},
			minGameWindow: 29,
			maxGameWindow: 30,
			expectedMatches: map[string]struct{}{
				"WWWWWWWWWWWWWWWWWWWWWWWWWWLWW":  {},
				"WWWWWWWWWWWWWWWWWWWWWWWWWWWLW":  {},
				"WWWWWWWWWWWWWWWWWWWWWWWWWWWLWW": {},
				"WWWWWWWWWWWWWWWWWWWWWWWWWWWWW":  {},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			combos := newGameCombos()
			err := findMatches(test.records, combos, test.minGameWindow, test.maxGameWindow, 0)
			require.NoError(t, err)
			if test.checkCombos {
				assert.Equal(t, test.expectedCombos, combos.combos)
			}
			assert.Equal(t, test.expectedMatches, combos.matches)
		})
	}
}
