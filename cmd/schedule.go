package cmd

import (
	"fmt"
	"strings"

	"github.com/internetdrew/bball/internal/nba"
	"github.com/internetdrew/bball/internal/util"
	"github.com/spf13/cobra"
)

var (
	upcoming bool
	recent   bool
)

var scheduleCmd = &cobra.Command{
	Use:   "schedule [team]",
	Short: "View a team's schedule",
	Long:  "Display upcoming games or recent games for a specific NBA team.\n\nUse --upcoming (-u) to see future games or --recent (-r) to see past games.\nDefaults to upcoming games if no flag is specified.",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		team := strings.TrimSpace(strings.ToLower(args[0]))
		if team == "" {
			return fmt.Errorf("please specify a team name")
		}

		// Default to upcoming if neither flag is set
		if !upcoming && !recent {
			upcoming = true
		}

		var games []nba.Game
		var err error

		if upcoming {
			games, err = nba.FetchTeamSchedule(team, "upcoming")
		} else {
			games, err = nba.FetchTeamSchedule(team, "recent")
		}

		if err != nil {
			return err
		}

		if len(games) == 0 {
			if upcoming {
				fmt.Println("No upcoming games found for that team.")
			} else {
				fmt.Println("No recent games found for that team.")
			}
			return nil
		}

		title := "Upcoming Games"
		if recent {
			title = "Recent Games"
		}

		fmt.Println(util.FormatTeamSchedule(games, team, title))
		return nil
	},
}

func init() {
	rootCmd.AddCommand(scheduleCmd)
	scheduleCmd.Flags().BoolVarP(&upcoming, "upcoming", "u", false, "Show upcoming games")
	scheduleCmd.Flags().BoolVarP(&recent, "recent", "r", false, "Show recent games")
}
