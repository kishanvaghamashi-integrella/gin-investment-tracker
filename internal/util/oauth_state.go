package util

import (
	"errors"
	"sync"
	"time"
)

const (
	oauthStateExpiry  = 5 * time.Minute
	oauthStateMaxSize = 500
)

var oauthStateStore = struct {
	mu    sync.Mutex
	store map[string]struct{}
}{
	store: make(map[string]struct{}),
}

// StoreOAuthState stores a state string that auto-expires after 5 minutes.
// Returns an error if the store is at capacity (DoS protection).
func StoreOAuthState(state string) error {
	s := &oauthStateStore
	s.mu.Lock()
	defer s.mu.Unlock()

	if len(s.store) >= oauthStateMaxSize {
		return errors.New("oauth state store is full, try again later")
	}

	s.store[state] = struct{}{}

	time.AfterFunc(oauthStateExpiry, func() {
		s.mu.Lock()
		defer s.mu.Unlock()
		delete(s.store, state)
	})

	return nil
}

// ValidateAndConsumeOAuthState returns true if the state exists (not yet
// expired or used), then removes it so it cannot be reused.
func ValidateAndConsumeOAuthState(state string) bool {
	s := &oauthStateStore
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.store[state]; !ok {
		return false
	}

	delete(s.store, state)
	return true
}
