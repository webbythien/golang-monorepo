package watcher

import (
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const tearDownTimeout = 10 * time.Second

func WatchKillSignal(cancel func()) {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	slog.Warn("Force exit after timeout...", slog.String("Timeout", tearDownTimeout.String()))
	time.AfterFunc(tearDownTimeout, func() {
		os.Exit(1)
	})
	go func() {
		<-sigChan
		slog.Info("Received signal again, force exit...")
		os.Exit(1)
	}()
	cancel()
}

type CleanFunc func()

var cleanFuncs []CleanFunc

func RegisterCleanFunc(closeFunc CleanFunc) {
	cleanFuncs = append(cleanFuncs, closeFunc)
}

func CleanUp() {
	for _, fun := range cleanFuncs {
		fun()
	}
	slog.Info("Closed all connections")
}
