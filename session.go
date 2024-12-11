package main

import (
	"fmt"
	"github.com/google/uuid"
	"log"
	"sync"
)

type SessionUser struct {
	id   uuid.UUID
	user *User
}

type Sessionstore struct {
	mu       sync.RWMutex
	sessions map[string]*SessionUser
	users    map[string]*User
}

func CreateSessionstore() Sessionstore {
	return Sessionstore{
		sessions: make(map[string]*SessionUser),
		users:    make(map[string]*User, 1),
	}
}

func (d *Sessionstore) GetUser(username string) (*User, error) {
	d.mu.Lock()
	defer d.mu.Unlock()
	user, ok := d.users[username]
	if !ok {
		return &User{}, fmt.Errorf("error getting user '%s': does not exist", username)
	}
	return user, nil
}

func (db *Sessionstore) GetSession(sessionID string) (*SessionUser, error) {
	db.mu.Lock()
	defer db.mu.Unlock()

	session, ok := db.sessions[sessionID]
	log.Println("Getting session", sessionID, ok, session)
	if !ok {
		return nil, fmt.Errorf("error getting session '%s': does not exist", sessionID)
	}
	return session, nil
}

/*
	func (d *Sessionstore) SaveSession(s *webauthn.SessionData, userid string) {
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
func (d *Sessionstore) SaveUser(u *User) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.users[u.name] = u
	fmt.Println("Saved user", u.name)
}

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

func (db *Sessionstore) DeleteSession(u string) {
	db.mu.Lock()
	defer db.mu.Unlock()

	delete(db.sessions, u)
}
