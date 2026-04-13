package blog

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type NoopStore struct{}

func (s *NoopStore) Get(key string) (string, error) {
	return "", ErrCacheMiss
}

func (s *NoopStore) Set(key string, value string, ttlSeconds int) error {
	return nil
}

var ErrCacheMiss = &CacheError{Message: "cache miss"}

type CacheError struct {
	Message string
}

func (e *CacheError) Error() string {
	return e.Message
}

type RedisStore struct {
	client         *redis.Client
	defaultTTL     time.Duration
	requestContext context.Context
}

func NewRedisStore(addr, password string, db int, ttlSeconds int) (*RedisStore, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	ctx := context.Background()
	if err := client.Ping(ctx).Err(); err != nil {
		_ = client.Close()
		return nil, err
	}

	ttl := time.Duration(ttlSeconds) * time.Second
	if ttlSeconds <= 0 {
		ttl = 5 * time.Minute
	}

	return &RedisStore{
		client:         client,
		defaultTTL:     ttl,
		requestContext: ctx,
	}, nil
}

func (s *RedisStore) Get(key string) (string, error) {
	value, err := s.client.Get(s.requestContext, key).Result()
	if err == redis.Nil {
		return "", ErrCacheMiss
	}
	return value, err
}

func (s *RedisStore) Set(key string, value string, ttlSeconds int) error {
	ttl := s.defaultTTL
	if ttlSeconds > 0 {
		ttl = time.Duration(ttlSeconds) * time.Second
	}
	return s.client.Set(s.requestContext, key, value, ttl).Err()
}

func (s *RedisStore) Close() error {
	return s.client.Close()
}
