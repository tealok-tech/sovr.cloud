package main

import (
	"fmt"
	"github.com/google/uuid"
	"log"
	"sync"
)

// A user's session which is used after authentication to continue to identify the user to the system.
type SessionUser struct {
	id   uuid.UUID
	user *User
}

// The store of all sessions. This is just stored in memory of the application and will therefore invalidate all sessions on process restart.
type Sessionstore struct {
	mu       sync.RWMutex
	sessions map[string]*SessionUser
}

// Create a new store of sessions
func CreateSessionstore() Sessionstore {
	return Sessionstore{
		sessions: make(map[string]*SessionUser),
	}
}

// Get session by session UUID
func (db *Sessionstore) GetSession(sessionID string) (*SessionUser, error) {
	db.mu.Lock()
	defer db.mu.Unlock()

	session, ok := db.sessions[sessionID]
	if !ok {
		return nil, fmt.Errorf("error getting session '%s': does not exist", sessionID)
	}
	return session, nil
}

// Start a new session for the given user and return the UUID as a string for storing in a cookie
func (db *Sessionstore) StartSession(u *User) string {

	db.mu.Lock()
	defer db.mu.Unlock()

	id := uuid.New()
	db.sessions[id.String()] = &SessionUser{
		id,
		u,
	}
	log.Println("Started user session for", u.name, id)
	return id.String()
}
