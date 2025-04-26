package main

import (
	"log/slog"
	"os"

	"github.com/monorepo/app/iam/cmd/api/iam"
	"github.com/spf13/cobra"
)

func run(_ []string) error {
	var rootCmd = &cobra.Command{
		Use:   "api",
		Short: "API service for pod",
		Long:  `API service for pod`,
		Run: func(cmd *cobra.Command, args []string) {
			slog.Info("Running API service")
		},
	}

	rootCmd.AddCommand(iam.Command)
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
