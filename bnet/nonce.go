package bnet

import (
	"sync"
	"time"
)

type NonceMap struct {
	sync.Mutex
	nonces map[string]time.Time
}

// Creates a new NonceMap and starts a goroutine to clean up expired nonces at the given interval.
func NewNonceMap(d time.Duration) *NonceMap {
	m := &NonceMap{
		nonces: make(map[string]time.Time),
	}
	go m.trimForever(d)
	return m
}

func (m *NonceMap) Add(nonce string, validDuration time.Duration) {
	m.Lock()
	m.nonces[nonce] = time.Now().Add(validDuration)
	defer m.Unlock()
}

func (m *NonceMap) Remove(nonce string) bool {
	m.Lock()
	defer m.Unlock()

	if expiration, ok := m.nonces[nonce]; ok {
		delete(m.nonces, nonce)
		return time.Now().Before(expiration)
	}
	return false
}

// trim removes all expired nonces from the map.
func (m *NonceMap) trim() {
	m.Lock()
	defer m.Unlock()

	now := time.Now()
	for nonce, expiration := range m.nonces {
		if now.After(expiration) {
			delete(m.nonces, nonce)
		}
	}
}

func (m *NonceMap) trimForever(d time.Duration) {
	t := time.NewTicker(d)
	for range t.C {
		m.trim()
	}
}
