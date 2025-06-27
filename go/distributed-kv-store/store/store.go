package store

import (
	"errors"
	"sync"
)

// Store represents a simple in-memory key-value store.
type Store struct {
	mu    sync.RWMutex
	data  map[string]string
}

// NewStore creates and returns a new Store instance.
func NewStore() *Store {
	return &Store{
		data: make(map[string]string),
	}
}

// Get retrieves the value associated with a key.
func (s *Store) Get(key string) (string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	val, ok := s.data[key]
	if !ok {
		return "", errors.New("key not found")
	}
	return val, nil
}

// Put stores a key-value pair.
func (s *Store) Put(key, value string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.data[key] = value
	return nil
}
