package main

import (
	"crypto/rand"
	"encoding/binary"

	"github.com/go-webauthn/webauthn/protocol"
	"github.com/go-webauthn/webauthn/webauthn"
	"log"
)

// A simple representation of a user
type User struct {
	ID          uint64                // A unique identifier, set by random
	Name        string                // The unique username
	DisplayName string                // The display of the user's name to them, meaningless to the system
	Credentials []webauthn.Credential // Credentials associated with the user
}

var UserAnonymous = &User{
	ID:          0,
	Name:        "",
	DisplayName: "Anonymous",
}

// Create a new user
func NewUser(name string, displayName string) *User {

	user := &User{}
	user.ID = randomUint64()
	user.Name = name
	user.DisplayName = displayName

	return user
}

// AddCredential associates the credential to the user
func (u *User) AddCredential(cred webauthn.Credential) {
	u.Credentials = append(u.Credentials, cred)
}

// CredentialExcludeList returns a CredentialDescriptor array filled
// with all the user's credentials
func (u User) CredentialExcludeList() []protocol.CredentialDescriptor {

	credentialExcludeList := []protocol.CredentialDescriptor{}
	for _, cred := range u.Credentials {
		descriptor := protocol.CredentialDescriptor{
			Type:         protocol.PublicKeyCredentialType,
			CredentialID: cred.ID,
		}
		credentialExcludeList = append(credentialExcludeList, descriptor)
	}

	return credentialExcludeList
}

func (u User) UpdateCredential(c *webauthn.Credential) {
	log.Println("This is the point where I should update the credential, but I don't know what that specifically means.")
}

// WebAuthnCredentials returns credentials owned by the user
func (u User) WebAuthnCredentials() []webauthn.Credential {
	return u.Credentials
}

// WebAuthnDisplayName returns the user's display name
func (u User) WebAuthnDisplayName() string {
	return u.DisplayName
}

// WebAuthnID returns the user's ID
func (u User) WebAuthnID() []byte {
	buf := make([]byte, binary.MaxVarintLen64)
	binary.PutUvarint(buf, uint64(u.ID))
	return buf
}

// WebAuthnName returns the user's username
func (u User) WebAuthnName() string {
	return u.Name
}

func randomUint64() uint64 {
	buf := make([]byte, 8)
	rand.Read(buf)
	return binary.LittleEndian.Uint64(buf)
}
