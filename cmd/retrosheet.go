package cmd

import (
	"encoding/csv"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"golang.org/x/sync/errgroup"
)

type Result int

const (
	Win Result = iota
	Loss
	Tie
)

type FranchiseConverter map[string]string

func (fc FranchiseConverter) Convert(team string) string {
	franchise, ok := fc[team]
	if ok {
		return franchise
	}
	return team
}

type TeamGame struct {
	Date               time.Time
	Team               string
	Franchise          string
	OpponentTeam       string
	OpponentFranchise  string
	OpponentGameNumber int
	TeamGameNumber     int
	OpponentScore      int
	TeamScore          int
	ForfeitInfo        string
	OpponentLineScore  string
	TeamLineScore      string
	Result             Result
}

type RetrosheetGame struct {
	Date               time.Time
	VisitingTeam       string
	VisitingGameNumber int
	HomeTeam           string
	HomeGameNumber     int
	VisitingScore      int
	HomeScore          int
	ForfeitInfo        string
	VisitingLineScore  string
	HomeLineScore      string
}

func (rg RetrosheetGame) GetHomeResult() Result {
	return calcRetrosheetResult(rg.HomeScore, rg.VisitingScore, true, rg.ForfeitInfo)
}

func (rg RetrosheetGame) GetVisitorResult() Result {
	return calcRetrosheetResult(rg.VisitingScore, rg.HomeScore, false, rg.ForfeitInfo)
}

func calcRetrosheetResult(teamScore, oppoScore int, isHome bool, forfeitInfo string) Result {
	if forfeitInfo != "" {
		if forfeitInfo == "T" {
			return Tie
		}
		if isHome {
			if forfeitInfo == "H" {
				return Win
			}
			if forfeitInfo == "V" {
				return Loss
			}
		} else {
			if forfeitInfo == "H" {
				return Loss
			}
			if forfeitInfo == "V" {
				return Win
			}
		}
	}

	if teamScore == oppoScore {
		return Tie
	}
	if teamScore > oppoScore {
		return Win
	}
	return Loss
}

func (rg RetrosheetGame) LineScoreProcessed(linescore string) ([]int, error) {
	var lookingForEnd bool
	var stringSoFar string
	scores := make([]int, 0, 9)
	for _, r := range linescore {
		if r == 'x' {
			continue
		}
		if r == '(' {
			lookingForEnd = true
			continue
		}
		if r == ')' {
			lookingForEnd = false
			score, err := strconv.Atoi(stringSoFar)
			if err != nil {
				return nil, err
			}
			scores = append(scores, score)
			stringSoFar = ""
			continue
		}
		if lookingForEnd {
			stringSoFar += string(r)
			continue
		}
		score, err := strconv.Atoi(string(r))
		if err != nil {
			return nil, err
		}
		scores = append(scores, score)
	}
	return scores, nil
}

type ByTeamsBySeason struct {
	mu                 sync.Mutex
	m                  map[string]map[int]*Season
	franchiseConverter FranchiseConverter
}

func newByTeamsBySeason(franchiseConverter FranchiseConverter) *ByTeamsBySeason {
	return &ByTeamsBySeason{
		m:                  make(map[string]map[int]*Season),
		franchiseConverter: franchiseConverter,
	}
}

func (btbs *ByTeamsBySeason) AddGame(game *RetrosheetGame) {
	addTeamResult := func(team string, game *RetrosheetGame) {
		franchise := btbs.franchiseConverter.Convert(team)
		year := game.Date.Year()
		teamMap, ok := btbs.m[franchise]
		if !ok {
			teamMap = make(map[int]*Season)
			btbs.m[franchise] = teamMap
		}
		season, ok := teamMap[year]
		if !ok {
			season = &Season{
				Franchise: franchise,
				Team:      team,
				Year:      year,
			}
			teamMap[year] = season
		}
		teamGame := retroGameToTeamGame(game, btbs.franchiseConverter, team == game.HomeTeam)
		season.Games = append(season.Games, teamGame)
	}
	btbs.mu.Lock()
	addTeamResult(game.HomeTeam, game)
	addTeamResult(game.VisitingTeam, game)
	btbs.mu.Unlock()
}

func (btbs *ByTeamsBySeason) SortGames() {
	btbs.mu.Lock()
	for _, seasonMap := range btbs.m {
		for _, season := range seasonMap {
			sort.Slice(season.Games, func(i, j int) bool {
				return season.Games[i].TeamGameNumber < season.Games[j].TeamGameNumber
			})
		}
	}
	btbs.mu.Unlock()
}

func (btbs *ByTeamsBySeason) BySortedSeason() map[string][]*Season {
	m := make(map[string][]*Season)
	btbs.mu.Lock()
	for team, seasonMap := range btbs.m {
		for _, season := range seasonMap {
			m[team] = append(m[team], season)
		}
	}
	btbs.mu.Unlock()
	for _, season := range m {
		sort.Slice(season, func(i, j int) bool { return season[i].Year < season[j].Year })
	}
	return m
}

type Season struct {
	Franchise string
	Team      string
	Year      int
	Games     []TeamGame // sorted by game number
}

func retroGameToTeamGame(rg *RetrosheetGame, franchiseConverter FranchiseConverter, isHome bool) TeamGame {
	tg := TeamGame{
		Date:        rg.Date,
		ForfeitInfo: rg.ForfeitInfo,
	}
	if isHome {
		tg.Team = rg.HomeTeam
		tg.Franchise = franchiseConverter.Convert(rg.HomeTeam)
		tg.OpponentTeam = rg.VisitingTeam
		tg.OpponentFranchise = franchiseConverter.Convert(rg.VisitingTeam)
		tg.OpponentGameNumber = rg.VisitingGameNumber
		tg.OpponentLineScore = rg.VisitingLineScore
		tg.OpponentScore = rg.VisitingScore
		tg.TeamGameNumber = rg.HomeGameNumber
		tg.TeamLineScore = rg.HomeLineScore
		tg.TeamScore = rg.HomeScore
		tg.Result = rg.GetHomeResult()
	} else {
		tg.Team = rg.VisitingTeam
		tg.Franchise = franchiseConverter.Convert(rg.VisitingTeam)
		tg.OpponentTeam = rg.HomeTeam
		tg.OpponentFranchise = franchiseConverter.Convert(rg.HomeTeam)
		tg.OpponentGameNumber = rg.HomeGameNumber
		tg.OpponentLineScore = rg.HomeLineScore
		tg.OpponentScore = rg.HomeScore
		tg.TeamGameNumber = rg.VisitingGameNumber
		tg.TeamLineScore = rg.VisitingLineScore
		tg.TeamScore = rg.VisitingScore
		tg.Result = rg.GetVisitorResult()
	}
	return tg
}

func GetTeamsBySeason(rsDataDir string) (*ByTeamsBySeason, error) {
	franchiseConverter, err := getFranchiseConverter(filepath.Join(rsDataDir, "misc/CurrentNames.csv"))
	if err != nil {
		return nil, err
	}
	btbs := newByTeamsBySeason(franchiseConverter)
	err = ByRetrosheetGame(rsDataDir, func(game *RetrosheetGame) error {
		btbs.AddGame(game)
		return nil
	})
	btbs.SortGames()
	return btbs, err
}

func getFranchiseConverter(path string) (FranchiseConverter, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	csvReader := csv.NewReader(f)
	records, err := csvReader.ReadAll()
	if err != nil {
		return nil, err
	}

	m := make(FranchiseConverter)
	for _, record := range records {
		team := record[1]
		franchise := record[0]
		m[team] = franchise
	}
	return m, nil
}

func ByRetrosheetGame(dir string, gameFunc func(*RetrosheetGame) error) error {
	gameDir := filepath.Join(dir, "games")
	files, err := os.ReadDir(gameDir)
	if err != nil {
		return err
	}

	var eg errgroup.Group

	for _, file := range files {
		if file.IsDir() {
			continue
		}
		if !strings.HasPrefix(file.Name(), "gl") {
			continue
		}
		f, err := os.Open(filepath.Join(gameDir, file.Name()))
		if err != nil {
			return err
		}
		defer f.Close()

		csvReader := csv.NewReader(f)
		records, err := csvReader.ReadAll()
		if err != nil {
			return err
		}

		for _, record := range records {
			r := record
			eg.Go(func() error {
				game, err := processRetrosheetGame(r)
				if err != nil {
					return err
				}
				return gameFunc(game)
			})
		}
	}
	return eg.Wait()
}

func processRetrosheetGame(record []string) (*RetrosheetGame, error) {
	visitingGameNumber, err := strconv.Atoi(record[5])
	if err != nil {
		return nil, err
	}
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
	date, err := time.Parse("20060102", record[0])
	if err != nil {
		return nil, err
	}

	return &RetrosheetGame{
		Date:               date,
		VisitingTeam:       record[3],
		VisitingGameNumber: visitingGameNumber,
		HomeTeam:           record[6],
		HomeGameNumber:     homeGameNumber,
		VisitingScore:      visitingScore,
		HomeScore:          homeScore,
		ForfeitInfo:        record[14],
		VisitingLineScore:  record[19],
		HomeLineScore:      record[20],
	}, nil
}
