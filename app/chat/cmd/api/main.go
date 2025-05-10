package main

import (
	"log/slog"
	"os"

	"github.com/monorepo/app/chat/cmd/api/chat"
	"github.com/spf13/cobra"
)

func run(_ []string) error {
	var rootCmd = &cobra.Command{
		Use:   "api",
		Short: "API service for chat",
		Long:  `API service for chat`,
		Run: func(cmd *cobra.Command, args []string) {
			slog.Info("Running API service")
		},
	}

	rootCmd.AddCommand(chat.Command)
	if err := rootCmd.Execute(); err != nil {
		return err
	}
	return nil
}

func main() {
	if err := run(os.Args); err != nil {
		slog.Error("run failed", "err", err)
	}
}
