package cache

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	redis "github.com/redis/go-redis/v9"
)

// RedisConfig ...
type RedisConfig struct {
	Addr       string `yaml:"addr" mapstructure:"addr"`
	Password   string `yaml:"password" mapstructure:"password"`
	DB         int    `yaml:"db" mapstructure:"db"`
	TLSEnabled bool   `yaml:"tls_enabled" mapstructure:"tls_enabled"`
}

var ErrNotFound = errors.New("not found")

type Cache interface {
	Set(ctx context.Context, key string, ttl time.Duration, target interface{}) error
	Get(ctx context.Context, key string, target interface{}) error
	Remove(ctx context.Context, key string) error
	RemoveByPattern(ctx context.Context, pattern string) error
}

type redisCache struct {
	client *redis.Client
}

func NewRedisCache(client *redis.Client) Cache {
	return &redisCache{
		client: client,
	}
}

func ConnectRedis(cfg *RedisConfig) *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.Addr,
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		fmt.Println("Error setting key: ", err)
	}

	return client
}

func (c *redisCache) Set(ctx context.Context, key string, ttl time.Duration, target interface{}) error {
	encoded, err := encodeAny(target)
	if err != nil {
		return err
	}
	return c.client.Set(ctx, key, encoded, ttl).Err()
}

func (c *redisCache) Get(ctx context.Context, key string, target interface{}) error {
	encoded, err := c.client.Get(ctx, key).Bytes()
	if errors.Is(err, redis.Nil) {
		return ErrNotFound
	}
	if err != nil {
		return err
	}

	return decodeAny(encoded, target)
}

func (c *redisCache) Remove(ctx context.Context, key string) error {
	return c.client.Del(ctx, key).Err()
}

func (c *redisCache) RemoveByPattern(ctx context.Context, pattern string) error {
	keys := c.client.Keys(ctx, pattern)
	if len(keys.Val()) > 0 {
		return c.client.Del(ctx, keys.Val()...).Err()
	}
	return nil
}

func MakeCacheKey(ss ...string) string {
	if len(ss) == 0 {
		return ""
	}
	var sb strings.Builder
	for _, s := range ss[:len(ss)-1] {
		sb.WriteString(s)
		sb.WriteString("|")
	}
	sb.WriteString(ss[len(ss)-1])
	return sb.String()
}

func encodeAny(any interface{}) ([]byte, error) {
	return json.Marshal(any)
}

func decodeAny(encoded []byte, any interface{}) error {
	return json.Unmarshal(encoded, any)
}
