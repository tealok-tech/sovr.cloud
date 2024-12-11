package main

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"github.com/go-webauthn/webauthn/webauthn"
	"sync"
)

type Authstore struct {
	mu       sync.RWMutex
	sessions map[string]*webauthn.SessionData
	users    map[string]*User
}

func random(length int) (string, error) {
	randomData := make([]byte, length)
	_, err := rand.Read(randomData)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(randomData), nil
}

func CreateAuthstore() Authstore {
	return Authstore{
		sessions: make(map[string]*webauthn.SessionData),
	}
}

func (db *Authstore) GetSession(sessionID string) (*webauthn.SessionData, error) {
	db.mu.Lock()
	defer db.mu.Unlock()

	session, ok := db.sessions[sessionID]
	if !ok {
		return nil, fmt.Errorf("error getting session '%s': does not exist", sessionID)
	}

	return session, nil
}

// PutUser stores a new user by the user's username
func (db *Authstore) StartSession(data *webauthn.SessionData) string {

	db.mu.Lock()
	defer db.mu.Unlock()

	sessionId, _ := random(32)
	db.sessions[sessionId] = data

	return sessionId
}

func (db *Authstore) DeleteSession(sessionID string) {
	db.mu.Lock()
	defer db.mu.Unlock()

	delete(db.sessions, sessionID)
}
