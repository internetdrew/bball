package util

import (
	"testing"

	"github.com/internetdrew/bball/internal/nba"
)

func TestFormatGameDate(t *testing.T) {
	date := "2023-11-04T19:00:00Z"
	formatted := FormatGameDate(date)
	if formatted == date || len(formatted) < 10 {
		t.Errorf("Expected formatted date, got %s", formatted)
	}
}

func TestFormatScore(t *testing.T) {
	game := nba.Game{
		HomeTeam: nba.Team{Tricode: "BKN", Score: 109},
		AwayTeam: nba.Team{Tricode: "MIN", Score: 125},
	}
	result := FormatScore(game)
	if result != "BKN 109 - MIN 125" {
		t.Errorf("Unexpected score format: %s", result)
	}
}

func TestFormatTopPerformers(t *testing.T) {
	players := []nba.PlayerStats{
		{PlayerName: "Player1", TeamCode: "BKN", Points: 30, Rebounds: 10, Assists: 5},
		{PlayerName: "Player2", TeamCode: "MIN", Points: 25, Rebounds: 8, Assists: 7},
	}
	result := FormatTopPerformers(players)
	if len(result) == 0 || result == "" {
		t.Error("Expected non-empty top performers string")
	}
}
