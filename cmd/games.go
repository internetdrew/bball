package cmd

import (
	"fmt"

	"github.com/internetdrew/bball/internal/nba"
	"github.com/internetdrew/bball/internal/util"
	"github.com/spf13/cobra"
)

var (
	liveOnly  bool
	finalOnly bool
)

var gamesCmd = &cobra.Command{
	Use:   "games",
	Short: "Show all NBA games for today",
	Long:  "Display a list of all NBA games scheduled for today, with optional filters for live or completed games.",
	RunE: func(cmd *cobra.Command, args []string) error {
		scoreboard, err := nba.FetchScoreboard()
		if err != nil {
			return fmt.Errorf("failed to fetch games: %w", err)
		}

		if len(scoreboard.Scoreboard.Games) == 0 {
			fmt.Println("No games scheduled for today.")
			return nil
		}

		// Filter games based on flags
		games := filterGames(scoreboard.Scoreboard.Games)

		if len(games) == 0 {
			if liveOnly {
				fmt.Println("No live games at the moment.")
			} else if finalOnly {
				fmt.Println("No completed games yet.")
			}
			return nil
		}

		fmt.Println(util.FormatGamesList(games))
		return nil
	},
}

func filterGames(games []nba.Game) []nba.Game {
	if !liveOnly && !finalOnly {
		return games
	}

	filtered := []nba.Game{}
	for _, g := range games {
		// GameStatus: 1 = Scheduled, 2 = Live, 3 = Final
		if liveOnly && g.GameStatus == 2 {
			filtered = append(filtered, g)
		} else if finalOnly && g.GameStatus == 3 {
			filtered = append(filtered, g)
		}
	}
	return filtered
}

func init() {
	rootCmd.AddCommand(gamesCmd)
	gamesCmd.Flags().BoolVarP(&liveOnly, "live", "l", false, "Show only live games")
	gamesCmd.Flags().BoolVarP(&finalOnly, "final", "f", false, "Show only completed games")
}
