package cmd

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"
	"github.com/transactions-platform/internal/app"
)

var apiCommand = &cobra.Command{
	Use:   "api",
	Short: "Initializes product api",
	RunE:  apiExecute,
}

func init() {
	rootCmd.AddCommand(apiCommand)
}

func apiExecute(cmd *cobra.Command, args []string) error {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	a, err := app.Build(ctx)
	if err != nil {
		return err
	}
	defer a.Close(ctx)

	return a.Run(ctx)
}

