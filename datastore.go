package main

import (
	"encoding/gob"
	"errors"
	"github.com/go-webauthn/webauthn/webauthn"
	"github.com/google/uuid"
	"os"
)

type Datastore struct {
	sessions []webauthn.SessionData
}

func (d *Datastore) GetOrCreateUser(username string, displayname string) (User, error) {
	return User{
		displayname: displayname,
		name:        username,
		id:          uuid.New().String(),
	}, nil
}

func (d *Datastore) GetUser(username string) (User, error) {
	return User{}, errors.New("No such user")
}

func (d *Datastore) GetSession() webauthn.SessionData {
	return webauthn.SessionData{}
}
func (d *Datastore) SaveSession(u *webauthn.SessionData) {
	d.sessions = append(d.sessions, *u)

	f, err := os.Open("datastore.glob")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	enc := gob.NewEncoder(f)
	err = enc.Encode(d)
}
func (d *Datastore) SaveUser(u User) {
}
