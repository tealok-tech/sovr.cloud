package main

import (
	"crypto/rand"
	"encoding/binary"

	"github.com/go-webauthn/webauthn/protocol"
	"github.com/go-webauthn/webauthn/webauthn"
)

// A simple representation of a user
type User struct {
	id          uint64                // A unique identifier, set by random
	name        string                // The unique username
	displayName string                // The display of the user's name to them, meaningless to the system
	credentials []webauthn.Credential // Credentials associated with the user
}

var UserAnonymous = &User{
	id:          0,
	name:        "",
	displayName: "Anonymous",
}

// Create a new user
func NewUser(name string, displayName string) *User {

	user := &User{}
	user.id = randomUint64()
	user.name = name
	user.displayName = displayName

	return user
}

// AddCredential associates the credential to the user
func (u *User) AddCredential(cred webauthn.Credential) {
	u.credentials = append(u.credentials, cred)
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

// WebAuthnCredentials returns credentials owned by the user
func (u User) WebAuthnCredentials() []webauthn.Credential {
	return u.credentials
}

// WebAuthnDisplayName returns the user's display name
func (u User) WebAuthnDisplayName() string {
	return u.displayName
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

func randomUint64() uint64 {
	buf := make([]byte, 8)
	rand.Read(buf)
	return binary.LittleEndian.Uint64(buf)
}
