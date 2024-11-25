package main

import (
	"errors"
	"github.com/go-webauthn/webauthn/webauthn"
)

type Datastore struct {
}

func (d *Datastore) GetUser(username string) (User, error) {
	return User{}, errors.New("No such user")
}
func (d *Datastore) GetSession() webauthn.SessionData {
	return webauthn.SessionData{}
}
func (d *Datastore) SaveSession(u *webauthn.SessionData) {
}
func (d *Datastore) SaveUser(u User) {
}
