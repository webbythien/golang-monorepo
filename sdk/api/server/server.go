package server

import (
	"context"
	"errors"
	"log"
	"log/slog"
	"net/http"
	"time"

	"connectrpc.com/connect"
	"connectrpc.com/otelconnect"
	"connectrpc.com/validate"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/prometheus"
	"go.opentelemetry.io/otel/sdk/metric"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Server struct {
	*http.Server
}

func New(cfg Config) *Server {
	return &Server{
		Server: &http.Server{
			Addr:         cfg.HTTP.String(),
			ReadTimeout:  10 * time.Second,
			WriteTimeout: 10 * time.Second,
		},
	}
}

func (s *Server) Register(f func(mux *http.ServeMux)) error {
	mux := http.NewServeMux()

	mux.Handle("/metrics", promhttp.Handler())
	f(mux)
	s.Handler = h2c.NewHandler(mux, &http2.Server{}) // Use h2c so we can serve HTTP/2 without TLS
	return nil
}

func (s *Server) ListenAndServe(ctx context.Context) error {
	go func(ctx context.Context) {
		<-ctx.Done()

		shutdownCtx, shutdownRelease := context.WithTimeout(ctx, 10*time.Second)
		defer shutdownRelease()

		if err := s.Server.Shutdown(shutdownCtx); err != nil {
			log.Fatalf("HTTP shutdown error: %v", err)
		}
		slog.Info("Graceful shutdown complete.")
	}(ctx)
	if err := s.Server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	return nil
}

var otelPromExporter *prometheus.Exporter

func init() {
	var err error
	otelPromExporter, err = prometheus.New()
	if err != nil {
		panic(err)
	} else {
		otel.SetMeterProvider(metric.NewMeterProvider(metric.WithReader(otelPromExporter)))
	}
}

func (s *Server) WithRecommendedOptions() []connect.HandlerOption {
	var opts = []connect.HandlerOption{
		connect.WithCodec(&protoJSONCodec{name: "json"}),
	}
	validateInterceptor, err := validate.NewInterceptor()
	if err != nil {
		slog.Error("failed to create validate interceptor", "err", err)
	} else {
		opts = append(opts, connect.WithInterceptors(validateInterceptor))
	}
	// loggingInterceptor := interceptor2.NewLoggingInterceptor()
	// opts = append(opts, connect.WithInterceptors(loggingInterceptor))

	if otelPromExporter != nil {
		otelInterceptor, err := otelconnect.NewInterceptor(
			otelconnect.WithoutServerPeerAttributes(),
		)
		if err != nil {
			slog.Error("failed to create otel interceptor", "err", err)
		} else {
			opts = append(opts, connect.WithInterceptors(otelInterceptor))
		}
	}

	// i18nInterceptor := i18n.Unary()
	// bizErrInterceptor := bizerr.NewInterceptor()
	// opts = append(opts, connect.WithInterceptors(i18nInterceptor, bizErrInterceptor))
	return opts
}
