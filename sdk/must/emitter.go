package must

import (
	"context"
	"crypto/tls"

	"github.com/monorepo/pkg/emitter"
	"github.com/monorepo/pkg/l"
	"github.com/monorepo/pkg/watcher"
	"github.com/redis/go-redis/extra/redisotel/v9"
	"github.com/redis/go-redis/v9"
)

func ConnectEmitter(cfg *emitter.EmitterConfig) *emitter.Emitter {
	ll := l.New()

	connectOptions := &redis.Options{
		Addr:     cfg.RedisAddr,
		Password: cfg.RedisPassword,
		DB:       cfg.RedisDB,
	}

	if cfg.TLSEnabled {
		connectOptions.TLSConfig = &tls.Config{
			InsecureSkipVerify: true, // #nosec G402
		}
	}

	client := redis.NewClient(connectOptions)

	if err := redisotel.InstrumentTracing(client); err != nil {
		ll.Fatal("Error instrument tracing redis for Emitter", l.Error(err))
	}

	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		ll.Fatal("Error ping redis for Emitter", l.Error(err))
	}

	watcher.RegisterCleanFunc(func() {
		ll.Info("Closing Emitter Redis connection")
		if err := client.Close(); err != nil {
			ll.Info("Error disconnecting Emitter Redis", l.Error(err))
		}
	})

	prefix := cfg.Prefix
	if prefix == "" {
		prefix = emitter.DefaultRedisPrefix
	}

	uid := cfg.Uid
	if uid == "" {
		uid = emitter.DefaultUid
	}

	nsp := cfg.Nsp
	if nsp == "" {
		nsp = emitter.DefaultNsp
	}

	eventType := cfg.EventType
	if eventType == 0 {
		eventType = emitter.NormalEvent
	}

	// Tạo một instance Emitter mới
	em := emitter.Emitter{
		Redis:     client,
		Prefix:    prefix,
		EventType: eventType,
		Nsp:       nsp,
		Uid:       uid,
		Flags:     make(map[string]interface{}),
	}

	return &em
}
