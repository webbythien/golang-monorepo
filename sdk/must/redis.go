package must

import (
	"context"
	"crypto/tls"

	"github.com/redis/go-redis/extra/redisotel/v9"

	"github.com/monorepo/pkg/watcher"

	"github.com/monorepo/pkg/cache"
	"github.com/monorepo/pkg/l"
	"github.com/redis/go-redis/v9"
)

func ConnectRedis(cfg *cache.RedisConfig) *redis.Client {
	ll := l.New()
	connectOptions := &redis.Options{
		Addr:     cfg.Addr,
		Password: cfg.Password,
		DB:       cfg.DB,
	}
	if cfg.TLSEnabled {
		connectOptions.TLSConfig = &tls.Config{
			InsecureSkipVerify: true, // #nosec G402
		}
	}
	client := redis.NewClient(connectOptions)

	if err := redisotel.InstrumentTracing(client); err != nil {
		ll.Fatal("Error instrument tracing redis", l.Error(err))
	}

	_, err := client.Ping(context.TODO()).Result()
	if err != nil {
		ll.Fatal("Error ping redis", l.Error(err))
	}
	watcher.RegisterCleanFunc(func() {
		ll.Info("Closing redis connection")
		if err := client.Close(); err != nil {
			ll.Info("Error disconnect redis", l.Error(err))
		}
	})

	return client
}
