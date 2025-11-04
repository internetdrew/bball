package util

import (
	"fmt"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/internetdrew/bball/internal/nba"
)

// FormatGameDate converts ISO 8601 timestamp to human-readable format
func FormatGameDate(dateStr string) string {
	t, err := time.Parse(time.RFC3339, dateStr)
	if err != nil {
		return dateStr // fallback to original if parsing fails
	}

	// Convert to EST/EDT
	loc, err := time.LoadLocation("America/New_York")
	if err != nil {
		loc = time.UTC
	}
	t = t.In(loc)

	// Format as "Monday, January 2, 2006 at 3:04 PM EST"
	return t.Format("Monday, January 2, 2006 at 3:04 PM MST")
}

// FormatScore returns a concise string like "BKN 109 - MIN 125"
func FormatScore(game nba.Game) string {
	return fmt.Sprintf("%s %d - %s %d",
		game.HomeTeam.Tricode, game.HomeTeam.Score,
		game.AwayTeam.Tricode, game.AwayTeam.Score,
	)
}

func FormatTopPerformers(players []nba.PlayerStats) string {
	lines := make([]string, len(players))
	for i, p := range players {
		lines[i] = fmt.Sprintf("%s (%s) - %d PTS, %d REB, %d AST",
			p.PlayerName, p.TeamCode, p.Points, p.Rebounds, p.Assists)
	}
	return strings.Join(lines, "\n")
}

// FormatGameSummary returns a nicely formatted multi-line summary
func FormatGameSummary(summary nba.GameSummary) string {
	builder := strings.Builder{}
	builder.WriteString(fmt.Sprintf("🏀 %s vs %s — %s\n",
		summary.Game.HomeTeam.Name,
		summary.Game.AwayTeam.Name,
		summary.Game.GameStatusText,
	))
	builder.WriteString(fmt.Sprintf("📅 %s\n\n", FormatGameDate(summary.Game.GameTimeUTC)))
	builder.WriteString(FormatScore(summary.Game) + "\n\n")
	builder.WriteString("Top Performers:\n")
	builder.WriteString(FormatTopPerformers(summary.TopPerformers) + "\n\n")
	builder.WriteString(fmt.Sprintf("Last updated: %s\n", summary.LastUpdated))
	return builder.String()
}

func FormatGamesList(games []nba.Game) string {
	builder := strings.Builder{}

	// Colors
	green := color.New(color.FgGreen).SprintFunc()
	yellow := color.New(color.FgYellow).SprintFunc()
	cyan := color.New(color.FgCyan).SprintFunc()

	builder.WriteString(fmt.Sprintf("\n🏀 NBA Games - %d game(s)\n", len(games)))
	builder.WriteString(strings.Repeat("─", 60) + "\n\n")

	for i, game := range games {
		// Status indicator
		var status string
		switch game.GameStatus {
		case 1: // Scheduled
			status = cyan("⏰ Scheduled")
		case 2: // Live
			status = green(fmt.Sprintf("🔴 LIVE - Q%d %s", game.Period, game.GameClock))
		case 3: // Final
			status = yellow("✓ Final")
		default:
			status = game.GameStatusText
		}

		builder.WriteString(fmt.Sprintf("%s - %s\n", status, FormatGameDate(game.GameTimeUTC)))
		builder.WriteString(fmt.Sprintf("  %s %s vs %s %s\n",
			game.AwayTeam.Tricode,
			formatScore(game.AwayTeam.Score, game.GameStatus),
			game.HomeTeam.Tricode,
			formatScore(game.HomeTeam.Score, game.GameStatus),
		))

		// Show time for scheduled games
		if game.GameStatus == 1 {
			builder.WriteString(fmt.Sprintf("  %s\n", game.GameStatusText))
		}

		// Add separator between games (but not after the last one)
		if i < len(games)-1 {
			builder.WriteString("\n")
		}
	}

	builder.WriteString("\n")
	return builder.String()
}

// formatScore returns the score as a string, or empty if game hasn't started
func formatScore(score int, gameStatus int) string {
	if gameStatus == 1 { // Scheduled
		return ""
	}
	return fmt.Sprintf("%d", score)
}

// FormatTeamSchedule returns a formatted team schedule
func FormatTeamSchedule(games []nba.Game, teamQuery string, title string) string {
	builder := strings.Builder{}

	// Colors
	green := color.New(color.FgGreen).SprintFunc()
	yellow := color.New(color.FgYellow).SprintFunc()
	cyan := color.New(color.FgCyan).SprintFunc()
	bold := color.New(color.Bold).SprintFunc()

	builder.WriteString(fmt.Sprintf("\n📅 %s - %d game(s)\n", bold(title), len(games)))
	builder.WriteString(strings.Repeat("─", 60) + "\n\n")

	for i, game := range games {
		// Determine if team is home or away
		isHome := strings.Contains(strings.ToLower(game.HomeTeam.Name), teamQuery) ||
			strings.EqualFold(game.HomeTeam.Tricode, teamQuery)

		var opponent string
		var location string
		if isHome {
			opponent = game.AwayTeam.Tricode
			location = "vs"
		} else {
			opponent = game.HomeTeam.Tricode
			location = "@"
		}

		// Status and score
		var statusLine string
		switch game.GameStatus {
		case 1: // Scheduled
			statusLine = cyan(fmt.Sprintf("%s %s %s", location, opponent, game.GameStatusText))
		case 2: // Live
			var teamScore, oppScore int
			if isHome {
				teamScore = game.HomeTeam.Score
				oppScore = game.AwayTeam.Score
			} else {
				teamScore = game.AwayTeam.Score
				oppScore = game.HomeTeam.Score
			}
			statusLine = green(fmt.Sprintf("🔴 LIVE %s %s - %d vs %d (Q%d %s)",
				location, opponent, teamScore, oppScore, game.Period, game.GameClock))
		case 3: // Final
			var teamScore, oppScore int
			var result string
			if isHome {
				teamScore = game.HomeTeam.Score
				oppScore = game.AwayTeam.Score
			} else {
				teamScore = game.AwayTeam.Score
				oppScore = game.HomeTeam.Score
			}
			if teamScore > oppScore {
				result = green("W")
			} else {
				result = color.RedString("L")
			}
			statusLine = yellow(fmt.Sprintf("%s Final: %s %s - %d vs %d",
				result, location, opponent, teamScore, oppScore))
		default:
			statusLine = fmt.Sprintf("%s %s - %s", location, opponent, game.GameStatusText)
		}

		builder.WriteString(fmt.Sprintf("%s\n", statusLine))

		// Add separator between games (but not after the last one)
		if i < len(games)-1 {
			builder.WriteString("\n")
		}
	}

	builder.WriteString("\n")
	return builder.String()
}
