package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"os"
	"sync"
)

type Userstore struct {
	mu    sync.RWMutex
	users map[string]*User // Mapping of username to user
}

func CreateUserstore() Userstore {
	u := Userstore{
		users: make(map[string]*User, 1),
	}
	u.readStore()
	return u
}

func (d *Userstore) GetUser(username string) (*User, error) {
	d.mu.Lock()
	defer d.mu.Unlock()
	user, ok := d.users[username]
	if !ok {
		return &User{}, fmt.Errorf("error getting user '%s': does not exist", username)
	}
	return user, nil
}

func (d *Userstore) GobEncode() ([]byte, error) {
	var buf bytes.Buffer
	encoder := gob.NewEncoder(&buf)
	if err := encoder.Encode(len(d.users)); err != nil {
		return nil, err
	}
	for _, u := range d.users {
		if err := encoder.Encode(u); err != nil {
			return nil, err
		}
	}
	return buf.Bytes(), nil
}

func (d *Userstore) GobDecode(b []byte) error {
	buf := bytes.NewBuffer(b)
	decoder := gob.NewDecoder(buf)
	var users_len int
	if err := decoder.Decode(&users_len); err != nil {
		return err
	}
	users := make(map[string]*User, users_len)
	for i := 0; i < users_len; i++ {
		var user User
		if err := decoder.Decode(&user); err != nil {
			return err
		}
		users[user.name] = &user
	}
	d.users = users
	return nil
}

func (d *Userstore) SaveUser(u *User) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.users[u.name] = u
	d.writeStore()
	fmt.Println("Saved user", u.name)
}

func (d *Userstore) readStore() error {
	d.mu.Lock()
	defer d.mu.Unlock()

	f, err := os.Open("userstore.glob")
	if err != nil {
		return err
	}

	dec := gob.NewDecoder(f)
	if err = dec.Decode(d); err != nil {
		return err
	}
	return nil
}

func (d *Userstore) writeStore() {
	d.mu.Lock()
	defer d.mu.Unlock()

	f, err := os.Create("userstore.glob")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	enc := gob.NewEncoder(f)
	err = enc.Encode(d)
}
