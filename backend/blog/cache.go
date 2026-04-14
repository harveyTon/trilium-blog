package blog

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/harveyTon/trilium-blog/backend/pkg/logger"
	"github.com/redis/go-redis/v9"
)

type cachePolicy struct {
	Prefix         string
	Version        int
	TTLSeconds     int
	Preload        bool
	RefreshAhead   bool
	RefreshAtRatio float64
}

var (
	policyNotesList = cachePolicy{
		Prefix: "notes", Version: 1, TTLSeconds: 90,
		Preload: true, RefreshAhead: true, RefreshAtRatio: 0.3,
	}
	policyNote = cachePolicy{
		Prefix: "note", Version: 1, TTLSeconds: 300,
	}
	policyNoteContent = cachePolicy{
		Prefix: "note-content", Version: 1, TTLSeconds: 1800,
		Preload: true,
	}
	policyAttachmentMeta = cachePolicy{
		Prefix: "attachment-meta", Version: 1, TTLSeconds: 3600,
	}
	policyAttachmentData = cachePolicy{
		Prefix: "attachment-data", Version: 1, TTLSeconds: 3600,
	}
)

func (p cachePolicy) key(suffix string) string {
	return fmt.Sprintf("%s:v%d:%s", p.Prefix, p.Version, suffix)
}

func (p cachePolicy) ttl() time.Duration {
	return time.Duration(p.TTLSeconds) * time.Second
}

var ErrCacheMiss = &CacheError{Message: "cache miss"}

type CacheTypeStats struct {
	Name       string `json:"name"`
	KeyCount   int    `json:"keyCount"`
	MinTTL     string `json:"minTTL"`
	MaxTTL     string `json:"maxTTL"`
	TTLSeconds int    `json:"ttlSeconds"`
}

type CacheStats struct {
	RedisConnected bool             `json:"redisConnected"`
	Types          []CacheTypeStats `json:"types"`
}

var allPolicies = []cachePolicy{
	policyNotesList,
	policyNote,
	policyNoteContent,
	policyAttachmentMeta,
	policyAttachmentData,
}

func (c *cacheLayer) stats() CacheStats {
	result := CacheStats{}
	if c.store == nil || c.noop {
		return result
	}
	result.RedisConnected = true
	for _, p := range allPolicies {
		prefix := fmt.Sprintf("%s:v%d", p.Prefix, p.Version)
		pattern := fmt.Sprintf("%s:*", prefix)
		keys, err := c.store.Keys(pattern)
		if err != nil {
			result.Types = append(result.Types, CacheTypeStats{Name: p.Prefix, TTLSeconds: p.TTLSeconds})
			continue
		}
		ts := CacheTypeStats{
			Name:       p.Prefix,
			KeyCount:   len(keys),
			TTLSeconds: p.TTLSeconds,
		}
		if len(keys) > 0 {
			var minTTL, maxTTL time.Duration
			for i, k := range keys {
				ttl, err := c.store.TTL(k)
				if err != nil {
					continue
				}
				if i == 0 || ttl < minTTL {
					minTTL = ttl
				}
				if ttl > maxTTL {
					maxTTL = ttl
				}
			}
			if minTTL > 0 {
				ts.MinTTL = minTTL.Truncate(time.Second).String()
			}
			if maxTTL > 0 {
				ts.MaxTTL = maxTTL.Truncate(time.Second).String()
			}
		}
		result.Types = append(result.Types, ts)
	}
	return result
}

type CacheError struct {
	Message string
}

func (e *CacheError) Error() string {
	return e.Message
}

type Store interface {
	Get(key string) (string, error)
	Set(key string, value string, ttlSeconds int) error
	TTL(key string) (time.Duration, error)
	Del(keys ...string) error
	Keys(pattern string) ([]string, error)
}

type NoopStore struct{}

func (s *NoopStore) Get(key string) (string, error)            { return "", ErrCacheMiss }
func (s *NoopStore) Set(key string, value string, _ int) error { return nil }
func (s *NoopStore) TTL(key string) (time.Duration, error)     { return 0, ErrCacheMiss }
func (s *NoopStore) Del(keys ...string) error                  { return nil }
func (s *NoopStore) Keys(pattern string) ([]string, error)     { return nil, nil }

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

func (s *RedisStore) TTL(key string) (time.Duration, error) {
	d, err := s.client.TTL(s.requestContext, key).Result()
	if err != nil {
		return 0, err
	}
	if d < 0 {
		return 0, ErrCacheMiss
	}
	return d, nil
}

func (s *RedisStore) Del(keys ...string) error {
	return s.client.Del(s.requestContext, keys...).Err()
}

func (s *RedisStore) Keys(pattern string) ([]string, error) {
	return s.client.Keys(s.requestContext, pattern).Result()
}

func (s *RedisStore) Close() error {
	return s.client.Close()
}

type refreshGuard struct {
	mu       sync.Mutex
	inFlight map[string]struct{}
}

func newRefreshGuard() *refreshGuard {
	return &refreshGuard{inFlight: make(map[string]struct{})}
}

func (g *refreshGuard) tryStart(key string) bool {
	g.mu.Lock()
	defer g.mu.Unlock()
	if _, ok := g.inFlight[key]; ok {
		return false
	}
	g.inFlight[key] = struct{}{}
	return true
}

func (g *refreshGuard) done(key string) {
	g.mu.Lock()
	delete(g.inFlight, key)
	g.mu.Unlock()
}

type cacheLayer struct {
	store Store
	guard *refreshGuard
	noop  bool
}

func newCacheLayer(store Store) *cacheLayer {
	_, noop := store.(*NoopStore)
	return &cacheLayer{
		store: store,
		guard: newRefreshGuard(),
		noop:  noop,
	}
}

func (c *cacheLayer) get(policy cachePolicy, suffix string) (string, bool) {
	if c.store == nil {
		return "", false
	}
	val, err := c.store.Get(policy.key(suffix))
	if err != nil {
		logger.Debug(fmt.Sprintf("cache miss: %s", policy.key(suffix)))
		return "", false
	}
	logger.Debug(fmt.Sprintf("cache hit: %s", policy.key(suffix)))
	return val, true
}

func (c *cacheLayer) set(policy cachePolicy, suffix string, value string) {
	if c.store == nil {
		return
	}
	if err := c.store.Set(policy.key(suffix), value, policy.TTLSeconds); err != nil {
		logger.Error(fmt.Sprintf("cache set failed: %s", policy.key(suffix)), err)
	} else {
		logger.Debug(fmt.Sprintf("cache set: %s (ttl=%ds)", policy.key(suffix), policy.TTLSeconds))
	}
}

func (c *cacheLayer) maybeRefreshAhead(policy cachePolicy, suffix string, loader func()) {
	if !policy.RefreshAhead || c.store == nil {
		return
	}
	key := policy.key(suffix)
	ttl, err := c.store.TTL(key)
	if err != nil {
		return
	}
	threshold := time.Duration(float64(policy.TTLSeconds)*policy.RefreshAtRatio) * time.Second
	if ttl > threshold {
		return
	}
	if !c.guard.tryStart(key) {
		return
	}
	go func() {
		defer c.guard.done(key)
		logger.Debug(fmt.Sprintf("refresh-ahead: %s (remaining_ttl=%s)", key, ttl))
		loader()
	}()
}

func (c *cacheLayer) del(policy cachePolicy, suffix string) {
	if c.store == nil {
		return
	}
	key := policy.key(suffix)
	if err := c.store.Del(key); err != nil {
		logger.Error(fmt.Sprintf("cache del failed: %s", key), err)
	} else {
		logger.Debug(fmt.Sprintf("cache del: %s", key))
	}
}

func (c *cacheLayer) delByPrefix(prefix string) int {
	if c.store == nil {
		return 0
	}
	pattern := fmt.Sprintf("%s:*", prefix)
	keys, err := c.store.Keys(pattern)
	if err != nil {
		logger.Error(fmt.Sprintf("cache keys failed: %s", pattern), err)
		return 0
	}
	if len(keys) == 0 {
		return 0
	}
	if err := c.store.Del(keys...); err != nil {
		logger.Error(fmt.Sprintf("cache del-by-prefix failed: %s", prefix), err)
		return 0
	}
	logger.Debug(fmt.Sprintf("cache del-by-prefix: %s (%d keys)", prefix, len(keys)))
	return len(keys)
}

func (c *cacheLayer) readJSON(policy cachePolicy, suffix string, dest interface{}) bool {
	raw, ok := c.get(policy, suffix)
	if !ok {
		return false
	}
	return json.Unmarshal([]byte(raw), dest) == nil
}

func (c *cacheLayer) writeJSON(policy cachePolicy, suffix string, val interface{}) {
	data, err := json.Marshal(val)
	if err != nil {
		return
	}
	c.set(policy, suffix, string(data))
}
