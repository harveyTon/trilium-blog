package blog

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
