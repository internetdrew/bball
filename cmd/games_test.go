package cmd

import (
	"testing"

	"github.com/internetdrew/bball/internal/nba"
)

func withFlags(live, final bool, fn func()) {
	oldLive, oldFinal := liveOnly, finalOnly
	liveOnly, finalOnly = live, final
	defer func() { liveOnly, finalOnly = oldLive, oldFinal }()
	fn()
}

func TestFilterGames_DefaultReturnsAll(t *testing.T) {
	games := []nba.Game{{GameStatus: 1}, {GameStatus: 2}, {GameStatus: 3}}
	withFlags(false, false, func() {
		out := filterGames(games)
		if len(out) != 3 {
			t.Fatalf("expected 3 games, got %d", len(out))
		}
	})
}

func TestFilterGames_LiveOnly(t *testing.T) {
	games := []nba.Game{{GameStatus: 1}, {GameStatus: 2}, {GameStatus: 2}, {GameStatus: 3}}
	withFlags(true, false, func() {
		out := filterGames(games)
		if len(out) != 2 {
			t.Fatalf("expected 2 live games, got %d", len(out))
		}
		for _, g := range out {
			if g.GameStatus != 2 {
				t.Fatalf("non-live game in result")
			}
		}
	})
}

func TestFilterGames_FinalOnly(t *testing.T) {
	games := []nba.Game{{GameStatus: 1}, {GameStatus: 3}, {GameStatus: 3}}
	withFlags(false, true, func() {
		out := filterGames(games)
		if len(out) != 2 {
			t.Fatalf("expected 2 final games, got %d", len(out))
		}
		for _, g := range out {
			if g.GameStatus != 3 {
				t.Fatalf("non-final game in result")
			}
		}
	})
}
