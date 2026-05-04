package util

import (
	"sync"
	"time"
)

const oauthStateExpiry = 5 * time.Minute

type oauthStateEntry struct {
	expiresAt time.Time
}

var oauthStateStore = struct {
	mu    sync.Mutex
	store map[string]oauthStateEntry
}{
	store: make(map[string]oauthStateEntry),
}

// StoreOAuthState stores a state string that expires after 5 minutes.
func StoreOAuthState(state string) {
	s := &oauthStateStore
	s.mu.Lock()
	defer s.mu.Unlock()
	s.store[state] = oauthStateEntry{expiresAt: time.Now().Add(oauthStateExpiry)}
}

// ValidateAndConsumeOAuthState returns true if the state exists and has not
// expired, then removes it so it cannot be reused.
func ValidateAndConsumeOAuthState(state string) bool {
	s := &oauthStateStore
	s.mu.Lock()
	defer s.mu.Unlock()

	entry, ok := s.store[state]
	if !ok {
		return false
	}

	delete(s.store, state)

	return time.Now().Before(entry.expiresAt)
}
