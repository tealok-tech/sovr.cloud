package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-webauthn/webauthn/webauthn"
	"log"
	"net/http"
)

type User struct {
	displayname string
	name        string
	id          string
}

func (u User) AddCredential(c *webauthn.Credential) {
}
func (u User) UpdateCredential(c *webauthn.Credential) {
}
func (u User) WebAuthnID() []byte {
	return []byte(u.id)
}
func (u User) WebAuthnName() string {
	return u.name
}
func (u User) WebAuthnDisplayName() string {
	return u.displayname
}
func (u User) WebAuthnCredentials() []webauthn.Credential {
	return make([]webauthn.Credential, 0)
}

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

	datastore := CreateDatastore()

	r := gin.Default()
	r.StaticFile("/", "./static/index.html")
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	r.GET("/hello", getHello)
	r.Static("/static", "./static")
	r.GET("/login/begin", func(c *gin.Context) {
		username := c.Query("username")
		user, err := datastore.GetUser(username) // Find the user
		if err != nil {
			fmt.Println("Error on GetUser", err)
			c.JSON(400, gin.H{
				"error": err.Error(),
			})
		}

		options, session, err := webAuthn.BeginLogin(user)
		if err != nil {
			// Handle Error and return.
			fmt.Println("Error on login begin", err)
			return
		}

		// store the session values
		datastore.SaveSession(session, user.id)

		//JSONResponse(w, options, http.StatusOK) // return the options generated
		c.JSON(200, options)
		// options.publicKey contain our registration options
	})
	r.GET("/login/finish", func(c *gin.Context) {
		username := c.Query("username")
		user, err := datastore.GetUser(username) // Get the user

		// Get the session data stored from the function above
		session := datastore.GetSession(user.id)

		credential, err := webAuthn.FinishLogin(user, session, c.Request)
		if err != nil {
			// Handle Error and return.
			fmt.Println("Error on login finish", err)
			return
		}

		// Handle credential.Authenticator.CloneWarning

		// If login was successful, update the credential object
		// Pseudocode to update the user credential.
		user.UpdateCredential(credential)
		datastore.SaveUser(user)

		//JSONResponse(w, "Login Success", http.StatusOK)
		c.JSON(200, "login success")
	})
	r.GET("/register/begin", func(c *gin.Context) {
		username := c.Query("username")
		displayname := c.Query("displayname")
		user, err := datastore.GetOrCreateUser(username, displayname) // Find or create the new user
		log.Println("Got user", user)
		options, session, err := webAuthn.BeginRegistration(user)
		if err != nil {
			c.JSON(500, gin.H{
				"error": err,
			})
			return
		}
		datastore.SaveSession(session, user.id)
		log.Println("Got session", session)
		log.Println("Got options", options)
		// handle errors if present
		// store the sessionData values
		//JSONResponse(w, options, http.StatusOK) // return the options generated
		// options.publicKey contain our registration options
		c.JSON(200, options)
	})
	r.POST("/register/finish", func(c *gin.Context) {
		username := c.Query("username")
		user, err := datastore.GetUser(username) // Get the user

		// Get the session data stored from the function above
		session := datastore.GetSession(username)

		credential, err := webAuthn.FinishRegistration(user, session, c.Request)
		if err != nil {
			// Handle Error and return.

			return
		}

		// If creation was successful, store the credential object
		// Pseudocode to add the user credential.
		user.AddCredential(credential)
		datastore.SaveUser(user)
		log.Println("Saved new user on registration")

		//JSONResponse(w, "Registration Success", http.StatusOK) // Handle next steps
		c.JSON(200, "Registration Success")
	})
	_ = r.Run()
}

func getHello(c *gin.Context) {
	c.String(http.StatusOK, "Hello World!")
}
