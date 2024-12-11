package main

import (
	"bytes"
	"crypto/rand"
	"encoding/binary"
	"encoding/gob"
	"fmt"

	"github.com/go-webauthn/webauthn/protocol"
	"github.com/go-webauthn/webauthn/webauthn"
	"os"
	"sync"
)

// User represents the user model
type User struct {
	id          uint64
	name        string
	displayName string
	credentials []webauthn.Credential
}
type Userstore struct {
	mu    sync.RWMutex
	users map[string]*User
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

// NewUser creates and returns a new User
func NewUser(name string, displayName string) *User {

	user := &User{}
	user.id = randomUint64()
	user.name = name
	user.displayName = displayName
	// user.credentials = []webauthn.Credential{}

	return user
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

func randomUint64() uint64 {
	buf := make([]byte, 8)
	rand.Read(buf)
	return binary.LittleEndian.Uint64(buf)
}

// WebAuthnID returns the user's ID
func (u User) WebAuthnID() []byte {
	buf := make([]byte, binary.MaxVarintLen64)
	binary.PutUvarint(buf, uint64(u.id))
	return buf
}

// WebAuthnName returns the user's username
func (u User) WebAuthnName() string {
	return u.name
}

// WebAuthnDisplayName returns the user's display name
func (u User) WebAuthnDisplayName() string {
	return u.displayName
}

// WebAuthnIcon is not (yet) implemented
func (u User) WebAuthnIcon() string {
	return ""
}

// AddCredential associates the credential to the user
func (u *User) AddCredential(cred webauthn.Credential) {
	u.credentials = append(u.credentials, cred)
}

// WebAuthnCredentials returns credentials owned by the user
func (u User) WebAuthnCredentials() []webauthn.Credential {
	return u.credentials
}

// CredentialExcludeList returns a CredentialDescriptor array filled
// with all the user's credentials
func (u User) CredentialExcludeList() []protocol.CredentialDescriptor {

	credentialExcludeList := []protocol.CredentialDescriptor{}
	for _, cred := range u.credentials {
		descriptor := protocol.CredentialDescriptor{
			Type:         protocol.PublicKeyCredentialType,
			CredentialID: cred.ID,
		}
		credentialExcludeList = append(credentialExcludeList, descriptor)
	}

	return credentialExcludeList
}
