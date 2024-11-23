package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-webauthn/webauthn/webauthn"
	"log"
	"net/http"
)

type User struct {
}

func (u User) AddCredential(c *webauthn.Credential) {
}
func (u User) WebAuthnID() []byte {
	return make([]byte, 0)
}
func (u User) WebAuthnName() string {
	return ""
}
func (u User) WebAuthnDisplayName() string {
	return ""
}
func (u User) WebAuthnCredentials() []webauthn.Credential {
	return make([]webauthn.Credential, 0)
}

type Datastore struct {
}

func (d *Datastore) GetUser() User {
	return User{}
}
func (d *Datastore) GetSession() webauthn.SessionData {
	return webauthn.SessionData{}
}
func (d *Datastore) SaveUser(u User) {
}

var datastore Datastore

/*func BeginRegistration(w http.ResponseWriter, r *http.Request) error {
	user := datastore.GetUser() // Find or create the new user
	options, session, err := webAuthn.BeginRegistration(user)
	if err != nil {
		return err
	}
	log.Println("Got session", session)
	// handle errors if present
	// store the sessionData values
	JSONResponse(w, options, http.StatusOK) // return the options generated
	// options.publicKey contain our registration options
	return nil
}

func FinishRegistration(w http.ResponseWriter, r *http.Request) {
	user := datastore.GetUser() // Get the user

	// Get the session data stored from the function above
	session := datastore.GetSession()

	credential, err := webAuthn.FinishRegistration(user, session, r)
	if err != nil {
		// Handle Error and return.

		return
	}

	// If creation was successful, store the credential object
	// Pseudocode to add the user credential.
	user.AddCredential(credential)
	datastore.SaveUser(user)

	JSONResponse(w, "Registration Success", http.StatusOK) // Handle next steps
}*/

func main() {
	wconfig := &webauthn.Config{
		RPDisplayName: "sovr.io",                         // Display Name for your site
		RPID:          "sovr.io",                         // Generally the FQDN for your site
		RPOrigins:     []string{"https://login.sovr.io"}, // The origin URLs allowed for WebAuthn requests
	}
	var webAuthn *webauthn.WebAuthn
	var err error
	if webAuthn, err = webauthn.New(wconfig); err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Got webAuthn", webAuthn)

	r := gin.Default()
	r.StaticFile("/", "./static/index.html")
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	r.GET("/hello", getHello)
	r.Static("/static", "./static")
	r.GET("/register/begin", func(c *gin.Context) {
		user := datastore.GetUser() // Find or create the new user
		options, session, err := webAuthn.BeginRegistration(user)
		if err != nil {
			c.JSON(500, gin.H{
				"error": err,
			})
			return
		}
		log.Println("Got session", session)
		// handle errors if present
		// store the sessionData values
		//JSONResponse(w, options, http.StatusOK) // return the options generated
		// options.publicKey contain our registration options
		c.JSON(200, options)
	})
	r.GET("/register/finish", func(c *gin.Context) {
		user := datastore.GetUser() // Get the user

		// Get the session data stored from the function above
		session := datastore.GetSession()

		credential, err := webAuthn.FinishRegistration(user, session, c.Request)
		if err != nil {
			// Handle Error and return.

			return
		}

		// If creation was successful, store the credential object
		// Pseudocode to add the user credential.
		user.AddCredential(credential)
		datastore.SaveUser(user)

		//JSONResponse(w, "Registration Success", http.StatusOK) // Handle next steps
		c.JSON(200, "Registration Success")
	})
	_ = r.Run()
}

func getHello(c *gin.Context) {
	c.String(http.StatusOK, "Hello World!")
}
