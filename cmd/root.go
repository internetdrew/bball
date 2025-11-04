package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "bball",
	Short: "Catch up on NBA games from your terminal",
	Long:  "A fast CLI tool for checking live NBA scores, today's games, and your favorite team's recent and upcoming matchups.",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
