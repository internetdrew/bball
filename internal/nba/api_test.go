package nba_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/internetdrew/bball/internal/nba"
)

func TestFetchScoreboard_Success(t *testing.T) {
	// Mock server
	games := []nba.Game{{HomeTeam: nba.Team{Name: "Knicks", Tricode: "NYK"}, AwayTeam: nba.Team{Name: "Nets", Tricode: "BKN"}}}
	resp := nba.Scoreboard{}
	resp.Scoreboard.Games = games
	resp.Meta.Time = "now"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	// Patch scoreboardURL
	oldURL := nba.ScoreboardURL
	nba.ScoreboardURL = server.URL
	defer func() { nba.ScoreboardURL = oldURL }()

	result, err := nba.FetchScoreboard()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if len(result.Scoreboard.Games) != 1 {
		t.Errorf("Expected 1 game, got %d", len(result.Scoreboard.Games))
	}
}

func TestFindTeamGame_NotFound(t *testing.T) {
	// Mock server returns empty games
	resp := nba.Scoreboard{}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	oldURL := nba.ScoreboardURL
	nba.ScoreboardURL = server.URL
	defer func() { nba.ScoreboardURL = oldURL }()

	_, err := nba.FindTeamGame("Lakers")
	if err == nil || err.Error() != "team not found" {
		t.Errorf("Expected 'team not found' error, got %v", err)
	}
}

func TestFetchScoreboard_InvalidJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "{invalid json")
	}))
	defer server.Close()

	oldURL := nba.ScoreboardURL
	nba.ScoreboardURL = server.URL
	defer func() { nba.ScoreboardURL = oldURL }()

	_, err := nba.FetchScoreboard()
	if err == nil {
		t.Fatalf("expected error for invalid JSON, got nil")
	}
}

func TestFetchTeamSchedule_Upcoming(t *testing.T) {
	// Build a league schedule with two future scheduled games for BOS
	now := time.Now().UTC()
	games := []nba.Game{
		{HomeTeam: nba.Team{Name: "Boston Celtics", Tricode: "BOS"}, AwayTeam: nba.Team{Name: "New York Knicks", Tricode: "NYK"}, GameStatus: 1, GameDateTimeUTC: now.Add(2 * time.Hour).Format(time.RFC3339)},
		{HomeTeam: nba.Team{Name: "Miami Heat", Tricode: "MIA"}, AwayTeam: nba.Team{Name: "Boston Celtics", Tricode: "BOS"}, GameStatus: 1, GameDateTimeUTC: now.Add(4 * time.Hour).Format(time.RFC3339)},
		{HomeTeam: nba.Team{Name: "Chicago Bulls", Tricode: "CHI"}, AwayTeam: nba.Team{Name: "Boston Celtics", Tricode: "BOS"}, GameStatus: 3}, // completed, should be ignored in upcoming
	}
	payload := struct {
		LeagueSchedule struct {
			GameDates []struct {
				GameDate string     `json:"gameDate"`
				Games    []nba.Game `json:"games"`
			} `json:"gameDates"`
		} `json:"leagueSchedule"`
	}{}
	payload.LeagueSchedule.GameDates = []struct {
		GameDate string     `json:"gameDate"`
		Games    []nba.Game `json:"games"`
	}{
		{GameDate: now.Format("2006-01-02"), Games: games},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(payload)
	}))
	defer server.Close()

	oldURL := nba.LeagueScheduleURL
	nba.LeagueScheduleURL = server.URL
	defer func() { nba.LeagueScheduleURL = oldURL }()

	result, err := nba.FetchTeamSchedule("BOS", "upcoming")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 2 {
		t.Fatalf("expected 2 upcoming games, got %d", len(result))
	}
}

func TestFetchTeamSchedule_Recent(t *testing.T) {
	games := []nba.Game{
		{HomeTeam: nba.Team{Name: "Boston Celtics", Tricode: "BOS"}, AwayTeam: nba.Team{Name: "NY Knicks", Tricode: "NYK"}, GameStatus: 3},
		{HomeTeam: nba.Team{Name: "Boston Celtics", Tricode: "BOS"}, AwayTeam: nba.Team{Name: "Miami Heat", Tricode: "MIA"}, GameStatus: 3},
		{HomeTeam: nba.Team{Name: "Boston Celtics", Tricode: "BOS"}, AwayTeam: nba.Team{Name: "Chicago Bulls", Tricode: "CHI"}, GameStatus: 1}, // scheduled, ignored
	}
	payload := struct {
		LeagueSchedule struct {
			GameDates []struct {
				GameDate string     `json:"gameDate"`
				Games    []nba.Game `json:"games"`
			} `json:"gameDates"`
		} `json:"leagueSchedule"`
	}{}
	payload.LeagueSchedule.GameDates = []struct {
		GameDate string     `json:"gameDate"`
		Games    []nba.Game `json:"games"`
	}{{GameDate: time.Now().Format("2006-01-02"), Games: games}}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(payload)
	}))
	defer server.Close()

	oldURL := nba.LeagueScheduleURL
	nba.LeagueScheduleURL = server.URL
	defer func() { nba.LeagueScheduleURL = oldURL }()

	result, err := nba.FetchTeamSchedule("BOS", "recent")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 2 {
		t.Fatalf("expected 2 recent games, got %d", len(result))
	}
}

func TestFetchTeamSchedule_TeamNotFound(t *testing.T) {
	// Schedule contains no games for requested team
	payload := struct {
		LeagueSchedule struct {
			GameDates []struct {
				GameDate string     `json:"gameDate"`
				Games    []nba.Game `json:"games"`
			} `json:"gameDates"`
		} `json:"leagueSchedule"`
	}{}
	payload.LeagueSchedule.GameDates = []struct {
		GameDate string     `json:"gameDate"`
		Games    []nba.Game `json:"games"`
	}{{GameDate: time.Now().Format("2006-01-02"), Games: []nba.Game{}}}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(payload)
	}))
	defer server.Close()

	oldURL := nba.LeagueScheduleURL
	nba.LeagueScheduleURL = server.URL
	defer func() { nba.LeagueScheduleURL = oldURL }()

	_, err := nba.FetchTeamSchedule("BOS", "upcoming")
	if err == nil || err.Error() != "team not found" {
		t.Fatalf("expected 'team not found' error, got %v", err)
	}
}

func TestFetchTeamSchedule_HTTPError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusServiceUnavailable)
		fmt.Fprint(w, "unavailable")
	}))
	defer server.Close()

	oldURL := nba.LeagueScheduleURL
	nba.LeagueScheduleURL = server.URL
	defer func() { nba.LeagueScheduleURL = oldURL }()

	_, err := nba.FetchTeamSchedule("BOS", "upcoming")
	if err == nil {
		t.Fatalf("expected error on non-200 response, got nil")
	}
}
