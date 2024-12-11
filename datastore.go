package main

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"github.com/go-webauthn/webauthn/webauthn"
	"github.com/google/uuid"
	"log"
	"sync"
)

type SessionUser struct {
	id   uuid.UUID
	user *User
}

type Datastore struct {
	mu            sync.RWMutex
	sessions_auth map[string]*webauthn.SessionData
	sessions_user map[string]*SessionUser
	users         map[string]*User
}

func random(length int) (string, error) {
	randomData := make([]byte, length)
	_, err := rand.Read(randomData)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(randomData), nil
}
func CreateDatastore() Datastore {
	return Datastore{
		sessions_auth: make(map[string]*webauthn.SessionData),
		sessions_user: make(map[string]*SessionUser),
		users:         make(map[string]*User, 1),
	}
}

func (d *Datastore) GetUser(username string) (*User, error) {
	d.mu.Lock()
	defer d.mu.Unlock()
	user, ok := d.users[username]
	if !ok {
		return &User{}, fmt.Errorf("error getting user '%s': does not exist", username)
	}
	return user, nil
}

func (db *Datastore) GetSessionAuth(sessionID string) (*webauthn.SessionData, error) {
	db.mu.Lock()
	defer db.mu.Unlock()

	session, ok := db.sessions_auth[sessionID]
	if !ok {
		return nil, fmt.Errorf("error getting session '%s': does not exist", sessionID)
	}

	return session, nil
}
func (db *Datastore) GetSessionUser(sessionID string) (*SessionUser, error) {
	db.mu.Lock()
	defer db.mu.Unlock()

	session, ok := db.sessions_user[sessionID]
	log.Println("Getting session", sessionID, ok, session)
	if !ok {
		return nil, fmt.Errorf("error getting session '%s': does not exist", sessionID)
	}
	return session, nil
}

/*
	func (d *Datastore) SaveSession(s *webauthn.SessionData, userid string) {
		d.sessions[userid] = *s

		f, err := os.Create("datastore.glob")
		if err != nil {
			panic(err)
		}
		defer f.Close()
		enc := gob.NewEncoder(f)
		err = enc.Encode(d)
	}
*/
func (d *Datastore) SaveUser(u *User) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.users[u.name] = u
	fmt.Println("Saved user", u.name)
}

// PutUser stores a new user by the user's username
func (db *Datastore) StartSessionAuth(data *webauthn.SessionData) string {

	db.mu.Lock()
	defer db.mu.Unlock()

	sessionId, _ := random(32)
	db.sessions_auth[sessionId] = data

	return sessionId
}

func (db *Datastore) StartSessionUser(u *User) string {

	db.mu.Lock()
	defer db.mu.Unlock()

	id := uuid.New()
	db.sessions_user[id.String()] = &SessionUser{
		id,
		u,
	}
	log.Println("Started user session for", u.name, id)
	return id.String()
}

func (db *Datastore) DeleteSessionAuth(sessionID string) {
	db.mu.Lock()
	defer db.mu.Unlock()

	delete(db.sessions_auth, sessionID)
}

func (db *Datastore) DeleteSessionUser(u string) {
	db.mu.Lock()
	defer db.mu.Unlock()

	delete(db.sessions_user, u)
}
