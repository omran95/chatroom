package infrastructure

import (
	"context"
	"time"

	"github.com/omran95/chatroom/pkg/common"
	"github.com/omran95/chatroom/pkg/config"
	"github.com/redis/go-redis/extra/redisotel/v9"
	"github.com/redis/go-redis/v9"
)

var (
	RedisClient redis.UniversalClient

	expiration time.Duration
)

type RedisCache interface {
	HSet(ctx context.Context, key string, values ...interface{}) error
	HGet(ctx context.Context, key, field string) (string, error)
	HGetAll(ctx context.Context, key string) (map[string]string, error)
	HDel(ctx context.Context, key, field string) error
}

type RedisCacheImpl struct {
	client redis.UniversalClient
}

func NewRedisClient(config *config.Config) (redis.UniversalClient, error) {
	expiration = time.Duration(config.Redis.ExpirationHour) * time.Hour
	RedisClient = redis.NewClusterClient(&redis.ClusterOptions{
		Addrs:          common.GetServerAddrs(config.Redis.Addrs),
		Password:       config.Redis.Password,
		ReadOnly:       true,
		RouteByLatency: true,
		MinIdleConns:   config.Redis.MinIdleConn,
		PoolSize:       config.Redis.PoolSize,
		ReadTimeout:    time.Duration(config.Redis.ReadTimeoutMilliSecond) * time.Millisecond,
		WriteTimeout:   time.Duration(config.Redis.WriteTimeoutMilliSecond) * time.Millisecond,
		PoolTimeout:    5 * time.Second,
	})
	ctx := context.Background()
	_, err := RedisClient.Ping(ctx).Result()
	if err == redis.Nil || err != nil {
		return nil, err
	}
	if err = redisotel.InstrumentTracing(RedisClient); err != nil {
		return nil, err
	}
	return RedisClient, nil
}

func NewRedisCacheImpl(client redis.UniversalClient) *RedisCacheImpl {
	return &RedisCacheImpl{client}
}

func (rc *RedisCacheImpl) HGet(ctx context.Context, key, field string) (string, error) {
	return rc.client.HGet(ctx, key, field).Result()
}

func (rc *RedisCacheImpl) HGetAll(ctx context.Context, key string) (map[string]string, error) {
	return rc.client.HGetAll(ctx, key).Result()
}

func (rc *RedisCacheImpl) HSet(ctx context.Context, key string, values ...interface{}) error {
	return rc.client.HSet(ctx, key, values).Err()
}

func (rc *RedisCacheImpl) HDel(ctx context.Context, key, field string) error {
	return rc.client.HDel(ctx, key, field).Err()
}
