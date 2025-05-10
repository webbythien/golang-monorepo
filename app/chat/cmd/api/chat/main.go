package chat

import (
	"context"
	"log/slog"
	"net/http"

	"connectrpc.com/connect"
	"github.com/monorepo/app/chat/config"
	"github.com/spf13/cobra"
	"github.com/webbythien/monorepo/pkg/l"
	"github.com/webbythien/monorepo/sdk/api/server"
	"github.com/webbythien/monorepo/api/chat/v1/chatv1connect"
	"github.com/monorepo/app/chat/internal/services"
	"github.com/webbythien/monorepo/sdk/must"
	"fmt"
)

var Command = &cobra.Command{
	Use:   "chat",
	Short: "Chat API",
	Long:  `Chat API`,
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

	var (
		db = must.ConnectPostgreSQL(cfg.PostgreSQL)
	)
	ll.Info("Init DB Connection: Done")
	fmt.Println("DB: ", db)
	srv := server.New(cfg.Server)
	err := srv.Register(func(mux *http.ServeMux) {
		opts := append(srv.WithRecommendedOptions(), connect.WithInterceptors(
		// Implement after
		))
		mux.Handle(chatv1connect.NewChatAPIHandler(services.NewChatAPI(), opts...))

	})
	if err != nil {
		return err
	}
	return srv.ListenAndServe(ctx)
}
