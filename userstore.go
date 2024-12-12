package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"log"
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
	err := u.readStore()
	if err != nil {
		log.Println("Failed to read user store:", err)
	}
	return u
}

func (d *Userstore) GetUser(username string) *User {
	d.mu.Lock()
	defer d.mu.Unlock()
	user, ok := d.users[username]
	if !ok {
		return nil
	}
	return user
}

func (d Userstore) GobEncode() ([]byte, error) {
	// Create a buffer to store the encoded data
	var buf bytes.Buffer

	// Create a new encoder
	enc := gob.NewEncoder(&buf)

	// Create a temporary map to hold the dereferenced items
	items := make(map[string]User)
	for key, ptr := range d.users {
		if ptr != nil {
			items[key] = *ptr
			log.Println("Added", ptr.Name)
		}
	}
	log.Println("Added all users")

	// Encode the temporary map
	if err := enc.Encode(items); err != nil {
		return nil, fmt.Errorf("failed to encode items: %w", err)
	}
	log.Println("Encoded", items)

	return buf.Bytes(), nil
}

func (d *Userstore) GobDecode(b []byte) error {
	log.Println("Decoding user database")
	// Create a buffer with the input data
	buf := bytes.NewBuffer(b)

	// Create a new decoder
	dec := gob.NewDecoder(buf)

	// Create a temporary map to hold the decoded items
	var items map[string]User
	if err := dec.Decode(&items); err != nil {
		return fmt.Errorf("failed to decode items: %w", err)
	}

	// Convert the items back to pointers
	d.users = make(map[string]*User)
	for key, item := range items {
		itemCopy := item // Create a new variable to avoid pointer issues
		d.users[key] = &itemCopy
	}

	return nil
}

func (d *Userstore) SaveUser(u *User) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.users[u.Name] = u
	d.writeStore()
	fmt.Println("Saved user", u.Name)
}

func (d *Userstore) readStore() error {
	d.mu.Lock()
	defer d.mu.Unlock()

	log.Println("Reading user database")
	f, err := os.Open("userstore.glob")
	if err != nil {
		log.Println("Failed to open userstore.glob", err)
		return err
	}
	defer f.Close()

	dec := gob.NewDecoder(f)
	if err = dec.Decode(d); err != nil {
		log.Println("Failed to decode userstore.glob", err)
		return err
	}
	return nil
}

func (d *Userstore) writeStore() {
	f, err := os.Create("userstore.glob")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	enc := gob.NewEncoder(f)
	err = enc.Encode(*d)
	if err != nil {
		log.Println("Failed to write store", err)
	}
}
