package nba

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"
)

var (
	ScoreboardURL     = "https://cdn.nba.com/static/json/liveData/scoreboard/todaysScoreboard_00.json"
	LeagueScheduleURL = "https://cdn.nba.com/static/json/staticData/scheduleLeagueV2.json"
)

// ErrTeamNotFound indicates no current game for the provided team was found
// in today's scoreboard data.
var ErrTeamNotFound = errors.New("team not found")

type Scoreboard struct {
	Scoreboard struct {
		Games []Game `json:"games"`
	} `json:"scoreboard"`
	Meta struct {
		Time string `json:"time"`
	} `json:"meta"`
}

func FetchScoreboard() (*Scoreboard, error) {
	res, err := http.Get(ScoreboardURL)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var data Scoreboard
	if err := json.NewDecoder(res.Body).Decode(&data); err != nil {
		return nil, err
	}

	return &data, nil
}

func FindTeamGame(team string) (*Game, error) {
	board, err := FetchScoreboard()
	if err != nil {
		return nil, err
	}

	for _, g := range board.Scoreboard.Games {
		// Only consider live games for "catch"
		if g.GameStatus != 2 { // 1=Scheduled, 2=Live, 3=Final
			continue
		}
		if strings.Contains(strings.ToLower(g.HomeTeam.Name), team) ||
			strings.Contains(strings.ToLower(g.AwayTeam.Name), team) ||
			strings.EqualFold(g.HomeTeam.Tricode, team) ||
			strings.EqualFold(g.AwayTeam.Tricode, team) {
			return &g, nil
		}
	}

	return nil, ErrTeamNotFound
}

// LeagueScheduleResponse represents the full season schedule
type LeagueScheduleResponse struct {
	LeagueSchedule struct {
		GameDates []struct {
			GameDate string `json:"gameDate"`
			Games    []Game `json:"games"`
		} `json:"gameDates"`
	} `json:"leagueSchedule"`
}

func FetchTeamSchedule(teamQuery string, mode string) ([]Game, error) {
	// Fetch the full league schedule
	resp, err := http.Get(LeagueScheduleURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch schedule: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("schedule not available (status: %d)", resp.StatusCode)
	}

	var scheduleResp LeagueScheduleResponse
	if err := json.NewDecoder(resp.Body).Decode(&scheduleResp); err != nil {
		return nil, fmt.Errorf("failed to parse schedule: %w", err)
	}

	// Flatten all games and filter by team
	var teamGames []Game
	teamLower := strings.ToLower(teamQuery)

	for _, gameDate := range scheduleResp.LeagueSchedule.GameDates {
		for _, game := range gameDate.Games {
			// Check if this game involves the requested team
			if strings.Contains(strings.ToLower(game.HomeTeam.Name), teamLower) ||
				strings.Contains(strings.ToLower(game.AwayTeam.Name), teamLower) ||
				strings.EqualFold(game.HomeTeam.Tricode, teamQuery) ||
				strings.EqualFold(game.AwayTeam.Tricode, teamQuery) {
				teamGames = append(teamGames, game)
			}
		}
	}

	if len(teamGames) == 0 {
		return nil, errors.New("team not found")
	}

	// Filter based on mode
	now := time.Now()
	var filtered []Game

	switch mode {
	case "upcoming":
		// Get next 5 scheduled games (status = 1 for scheduled)
		count := 0
		for _, game := range teamGames {
			if game.GameStatus == 1 && count < 5 {
				// Parse the game date/time using gameDateTimeUTC
				var gameTime time.Time
				var err error

				if game.GameDateTimeUTC != "" {
					gameTime, err = time.Parse(time.RFC3339, game.GameDateTimeUTC)
				}

				if err == nil && !gameTime.IsZero() && gameTime.After(now) {
					filtered = append(filtered, game)
					count++
				}
			}
		}
	case "recent":

		for _, game := range teamGames {
			if game.GameStatus == 3 {
				filtered = append(filtered, game)
			}
		}

		if len(filtered) > 5 {
			filtered = filtered[len(filtered)-5:]
		}
	}

	return filtered, nil
}
