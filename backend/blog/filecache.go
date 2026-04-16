package blog

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/harveyTon/trilium-blog/backend/pkg/logger"
)

type FileStore struct {
	dir string
	mu  sync.RWMutex
}

type fileEntry struct {
	Value     string `json:"value"`
	ExpiresAt int64  `json:"expires_at"`
}

func NewFileStore(dir string) (*FileStore, error) {
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("create cache dir: %w", err)
	}
	return &FileStore{dir: dir}, nil
}

func (s *FileStore) filePath(key string) string {
	safe := strings.ReplaceAll(key, ":", "_")
	safe = strings.ReplaceAll(safe, "/", "_")
	return filepath.Join(s.dir, safe+".json")
}

func (s *FileStore) Get(key string) (string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	data, err := os.ReadFile(s.filePath(key))
	if err != nil {
		return "", ErrCacheMiss
	}
	var entry fileEntry
	if err := json.Unmarshal(data, &entry); err != nil {
		return "", ErrCacheMiss
	}
	if entry.ExpiresAt > 0 && time.Now().UnixMilli() > entry.ExpiresAt {
		_ = os.Remove(s.filePath(key))
		return "", ErrCacheMiss
	}
	return entry.Value, nil
}

func (s *FileStore) Set(key string, value string, ttlSeconds int) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	entry := fileEntry{Value: value}
	if ttlSeconds > 0 {
		entry.ExpiresAt = time.Now().Add(time.Duration(ttlSeconds) * time.Second).UnixMilli()
	}
	data, err := json.Marshal(entry)
	if err != nil {
		return err
	}
	return os.WriteFile(s.filePath(key), data, 0644)
}

func (s *FileStore) TTL(key string) (time.Duration, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	data, err := os.ReadFile(s.filePath(key))
	if err != nil {
		return 0, ErrCacheMiss
	}
	var entry fileEntry
	if err := json.Unmarshal(data, &entry); err != nil {
		return 0, ErrCacheMiss
	}
	if entry.ExpiresAt <= 0 {
		return 0, ErrCacheMiss
	}
	remaining := time.Until(time.UnixMilli(entry.ExpiresAt))
	if remaining <= 0 {
		return 0, ErrCacheMiss
	}
	return remaining, nil
}

func (s *FileStore) Del(keys ...string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, key := range keys {
		_ = os.Remove(s.filePath(key))
	}
	return nil
}

func (s *FileStore) Keys(pattern string) ([]string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var matches []string
	entries, err := os.ReadDir(s.dir)
	if err != nil {
		return nil, err
	}
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".json") {
			continue
		}
		storedKey := strings.TrimSuffix(entry.Name(), ".json")
		storedKey = strings.ReplaceAll(storedKey, "_", ":")
		if matched, _ := filepath.Match(pattern, storedKey); matched {
			matches = append(matches, storedKey)
			continue
		}
		if strings.Contains(pattern, "*") {
			prefix := strings.SplitN(pattern, "*", 2)[0]
			prefixNorm := strings.ReplaceAll(prefix, ":", "_")
			if strings.HasPrefix(entry.Name(), prefixNorm) {
				matches = append(matches, storedKey)
			}
		}
	}
	return matches, nil
}

func InitFileCache(dataDir string) Store {
	cacheDir := filepath.Join(dataDir, "cache")
	store, err := NewFileStore(cacheDir)
	if err != nil {
		logger.Error("Failed to initialize file cache; running without cache", err)
		return &NoopStore{}
	}
	logger.Info(fmt.Sprintf("Using file cache at %s", cacheDir))
	return store
}
