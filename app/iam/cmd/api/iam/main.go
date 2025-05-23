package iam

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	"connectrpc.com/connect"
	"github.com/monorepo/api/iam/v1/iamv1connect"
	"github.com/monorepo/app/iam/config"
	"github.com/monorepo/app/iam/internal/services"
	"github.com/monorepo/pkg/l"
	"github.com/monorepo/sdk/api/server"
	"github.com/monorepo/sdk/must"
	"github.com/spf13/cobra"
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
	cfg := config.Load()
	ll := l.New()
	ll.Debug("Config loaded", l.Object("configs", cfg))

	// ================== Initiate all dependencies ==================

	db := must.ConnectPostgreSQL(cfg.PostgreSQL)

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
