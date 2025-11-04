package util

import (
	"strings"
	"testing"

	"github.com/internetdrew/bball/internal/nba"
)

func TestFormatGameSummary_ContainsKeyFields(t *testing.T) {
	summary := nba.GameSummary{
		Game: nba.Game{
			HomeTeam:       nba.Team{Name: "Boston Celtics", Tricode: "BOS", Score: 100},
			AwayTeam:       nba.Team{Name: "New York Knicks", Tricode: "NYK", Score: 98},
			GameStatusText: "Final",
			GameTimeUTC:    "2023-11-04T19:00:00Z",
		},
		TopPerformers: []nba.PlayerStats{{PlayerName: "Tatum", TeamCode: "BOS", Points: 30, Rebounds: 8, Assists: 5}},
		LastUpdated:   "just now",
	}
	out := FormatGameSummary(summary)
	if !strings.Contains(out, "Boston Celtics") || !strings.Contains(out, "New York Knicks") {
		t.Fatalf("expected team names in summary: %s", out)
	}
	if !strings.Contains(out, "Final") {
		t.Fatalf("expected status text in summary: %s", out)
	}
}

func TestFormatGamesList_ScheduledHidesScores(t *testing.T) {
	games := []nba.Game{
		{GameStatus: 1, GameStatusText: "7:30 PM ET", GameTimeUTC: "2023-11-04T23:30:00Z", HomeTeam: nba.Team{Tricode: "BOS", Score: 0}, AwayTeam: nba.Team{Tricode: "NYK", Score: 0}},
	}
	out := FormatGamesList(games)
	if strings.Contains(out, "0 vs 0") {
		t.Fatalf("scheduled game should not show scores: %s", out)
	}
}

func TestFormatTeamSchedule_Basic(t *testing.T) {
	games := []nba.Game{
		{GameStatus: 1, GameStatusText: "7:30 PM ET", HomeTeam: nba.Team{Name: "Boston Celtics", Tricode: "BOS"}, AwayTeam: nba.Team{Name: "New York Knicks", Tricode: "NYK"}},
	}
	out := FormatTeamSchedule(games, "bos", "Upcoming Games")
	if !strings.Contains(out, "vs NYK") && !strings.Contains(out, "@ BOS") { // either formatting depending on home/away
		t.Fatalf("expected opponent tricode in schedule: %s", out)
	}
}
