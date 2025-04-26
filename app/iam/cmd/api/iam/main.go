package iam

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/monorepo/app/iam/config"
	"github.com/spf13/cobra"
	"github.com/webbythien/monorepo/pkg/l"
	"github.com/webbythien/monorepo/sdk/must"
)

var Command = &cobra.Command{
	Use:   "iam",
	Short: "Iam API",
	Long:  `Iam API`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := run(cmd.Context(), args); err != nil {
			slog.Error("Run failed", "err", err)
		}
	},
}

func run(ctx context.Context, _ []string) error {
	var cfg = config.Load()
	var ll = l.New()
	ll.Debug("Config loaded", l.Object("configs", cfg))

	// ================== Initiate all dependencies ==================
	var (
		db = must.ConnectPostgreSQL(cfg.PostgreSQL)
	)
	ll.Info("Init DB Connection: Done")

	// Implemt API
	fmt.Println("DB Connection: Done: ", db)
	return nil
}
