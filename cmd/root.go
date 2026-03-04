package cmd

import (
	"log/slog"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "transactions-platform",
	Short: "transactions-platform app",
}

func Execute() {
	defer func() {
		err := recover()
		if err != nil {
			slog.Error("unexpected error while executing command %v", err)
		}
	}()

	err := rootCmd.Execute()
	if err != nil {
		slog.Error("error while executing command", "error", err.Error())
	}
}
