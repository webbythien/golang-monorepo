package iam

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	"connectrpc.com/connect"
	"github.com/monorepo/app/iam/config"
	"github.com/monorepo/app/iam/internal/services"
	"github.com/spf13/cobra"
	"github.com/webbythien/monorepo/api/iam/v1/iamv1connect"
	"github.com/webbythien/monorepo/pkg/l"
	"github.com/webbythien/monorepo/sdk/api/server"
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

	srv := server.New(cfg.Server)
	err := srv.Register(func(mux *http.ServeMux) {
		opts := append(srv.WithRecommendedOptions(), connect.WithInterceptors(
		// Implement after
		))
		// mux.Handle(iamv1connect.NewStaffAccessAPIHandler(staffAccessAPI, opts...))
		mux.Handle(iamv1connect.NewSecurityTokenAPIHandler(services.NewIamTest(), opts...))

	})
	if err != nil {
		return err
	}
	return srv.ListenAndServe(ctx)
}
