package util

import (
	"strings"
	"testing"

	"github.com/internetdrew/bball/internal/nba"
)

func TestFormatTeamSchedule_IncludesDateForScheduled(t *testing.T) {
	games := []nba.Game{
		{
			GameStatus:     1,
			GameStatusText: "7:30 PM ET",
			GameTimeUTC:    "2023-11-04T23:30:00Z", // Sat Nov 4 in ET
			HomeTeam:       nba.Team{Name: "New York Knicks", Tricode: "NYK"},
			AwayTeam:       nba.Team{Name: "Minnesota Timberwolves", Tricode: "MIN"},
		},
	}

	out := FormatTeamSchedule(games, "nyk", "Upcoming Games")

	// Expect a short date like "Sat Nov 4" present in the output
	if !strings.Contains(out, "Sat Nov 4") {
		t.Fatalf("expected date to be included in schedule output, got: %s", out)
	}
}
