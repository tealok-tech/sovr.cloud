package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"log"
	"net/http"
	"sync"
)

const SESSION_COOKIE_NAME = "session"

// A user's session which is used after authentication to continue to identify the user to the system.
type SessionUser struct {
	id   uuid.UUID
	user *User
}

// The store of all sessions. This is just stored in memory of the application and will therefore invalidate all sessions on process restart.
type Sessionstore struct {
	mu       sync.RWMutex
	sessions map[string]*SessionUser
}

// Create a new store of sessions
func CreateSessionstore() Sessionstore {
	return Sessionstore{
		sessions: make(map[string]*SessionUser),
	}
}

// Get session by session UUID
func (db *Sessionstore) GetSession(c *gin.Context) (*SessionUser, error) {
	db.mu.Lock()
	defer db.mu.Unlock()

	cookie, err := c.Request.Cookie(SESSION_COOKIE_NAME)
	if err != nil {
		c.HTML(http.StatusOK, "hello.tmpl", gin.H{
			"displayname": "anonymous",
		})
		return nil, err
	}
	session, ok := db.sessions[cookie.Value]
	if !ok {
		return nil, fmt.Errorf("error getting session '%s': does not exist", cookie.Value)
	}
	return session, nil
}

// Start a new session for the given user and return the UUID as a string for storing in a cookie
func (db *Sessionstore) StartSession(c *gin.Context, u *User) string {

	db.mu.Lock()
	defer db.mu.Unlock()

	id := uuid.New()
	db.sessions[id.String()] = &SessionUser{
		id,
		u,
	}
	log.Println("Started user session for", u.name, id)
	c.SetCookie(
		SESSION_COOKIE_NAME,
		id.String(),
		60*60*24*14, // Session cookie, closes when the browser window closes
		"/",         // Valid for all paths
		"",
		false, // HTTPS only
		true,  // allow JavaScript access to the cookie
	)

	return id.String()
}
