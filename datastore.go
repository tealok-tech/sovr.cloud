package main

import (
	"encoding/gob"
	"errors"
	"github.com/go-webauthn/webauthn/webauthn"
	"github.com/google/uuid"
	"os"
)

type Datastore struct {
	sessions map[string]webauthn.SessionData
	users    []User
}

func CreateDatastore() Datastore {
	return Datastore{
		sessions: make(map[string]webauthn.SessionData),
		users:    make([]User, 1),
	}
}

func (d *Datastore) GetOrCreateUser(username string, displayname string) (User, error) {
	user := User{
		displayname: displayname,
		name:        username,
		id:          uuid.New().String(),
	}
	d.users = append(d.users, user)
	return user, nil
}

func (d *Datastore) GetUser(username string) (User, error) {
	for _, u := range d.users {
		if u.name == username {
			return u, nil
		}
	}
	return User{}, errors.New("No such user")
}

func (d *Datastore) GetSession(userid string) webauthn.SessionData {
	return d.sessions[userid]
}
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
func (d *Datastore) SaveUser(u User) {
}
