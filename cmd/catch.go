package cmd

import (
	"errors"
	"fmt"
	"strings"

	"github.com/internetdrew/bball/internal/nba"
	"github.com/internetdrew/bball/internal/util"
	"github.com/spf13/cobra"
)

var catchCmd = &cobra.Command{
	Use:   "catch [team]",
	Short: "Catch up on a current (live) game for a team",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return fmt.Errorf("please specify a team name")
		}

		team := strings.TrimSpace(strings.ToLower(args[0]))
		if team == "" {
			return fmt.Errorf("please specify a team name")
		}

		game, err := nba.FindTeamGame(team)
		if err != nil {
			if errors.Is(err, nba.ErrTeamNotFound) {
				fmt.Println("No current game found for that team.")
				return nil
			}
			return err
		}

		if game == nil {
			fmt.Println("No current game found for that team.")
			return nil
		}

		summary := nba.GameSummary{
			Game: *game,
			TopPerformers: []nba.PlayerStats{
				{
					PlayerName: game.GameLeaders.HomeLeaders.Name,
					TeamCode:   game.GameLeaders.HomeLeaders.TeamTricode,
					Points:     game.GameLeaders.HomeLeaders.Points,
					Rebounds:   game.GameLeaders.HomeLeaders.Rebounds,
					Assists:    game.GameLeaders.HomeLeaders.Assists,
				},
				{
					PlayerName: game.GameLeaders.AwayLeaders.Name,
					TeamCode:   game.GameLeaders.AwayLeaders.TeamTricode,
					Points:     game.GameLeaders.AwayLeaders.Points,
					Rebounds:   game.GameLeaders.AwayLeaders.Rebounds,
					Assists:    game.GameLeaders.AwayLeaders.Assists,
				},
			},
			LastUpdated: "just now",
		}

		fmt.Print(util.FormatGameSummary(summary))

		return nil
	},
}

func init() {
	rootCmd.AddCommand(catchCmd)
}
