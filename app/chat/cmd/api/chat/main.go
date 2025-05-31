package chat

import (
	"context"
	"log/slog"
	"net/http"

	"connectrpc.com/connect"
	"github.com/monorepo/api/chat/v1/chatv1connect"
	"github.com/monorepo/app/chat/config"
	"github.com/monorepo/app/chat/internal/repositories"
	"github.com/monorepo/app/chat/internal/services"
	"github.com/monorepo/pkg/emitter"
	"github.com/monorepo/pkg/l"
	"github.com/monorepo/sdk/api/server"
	"github.com/monorepo/sdk/must"
	"github.com/spf13/cobra"
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
	cfg := config.Load()
	ll := l.New()
	ll.Debug("Config loaded", l.Object("configs", cfg))

	db := must.ConnectPostgreSQL(cfg.PostgreSQL)
	ll.Info("Init DB Connection: Done")

	emitter := must.ConnectEmitter(&emitter.EmitterConfig{
		RedisAddr:     cfg.Redis.Addr,
		RedisPassword: cfg.Redis.Password,
		RedisDB:       cfg.Redis.DB,
	})
	ll.Info("Init Emitter Connection: Done")

	meetingStore := repositories.NewMeetingStore(db)

	srv := server.New(cfg.Server)
	err := srv.Register(func(mux *http.ServeMux) {
		opts := append(srv.WithRecommendedOptions(), connect.WithInterceptors(
		// Implement after
		))
		mux.Handle(chatv1connect.NewChatAPIHandler(services.NewChatAPI(meetingStore, emitter), opts...))
	})
	if err != nil {
		return err
	}
	return srv.ListenAndServe(ctx)
}
