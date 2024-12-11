package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-webauthn/webauthn/protocol"
	"github.com/go-webauthn/webauthn/webauthn"
	"log"
	"net/http"
	"os"
)

func main() {
	wconfig := &webauthn.Config{
		RPDisplayName: "localhost",                       // Display Name for your site
		RPID:          "localhost",                       // Generally the FQDN for your site
		RPOrigins:     []string{"http://localhost:8080"}, // The origin URLs allowed for WebAuthn requests
	}
	var webAuthn *webauthn.WebAuthn
	var err error
	if webAuthn, err = webauthn.New(wconfig); err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Got webAuthn", webAuthn)

	authstore := CreateAuthstore()
	sessionstore := CreateSessionstore()
	userstore := CreateUserstore()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	r := gin.Default()
	r.LoadHTMLGlob("templates/*")
	r.StaticFile("/", "./static/index.html")
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	r.GET("/hello", func(c *gin.Context) {
		session, _ := sessionstore.GetSession(c)
		var user *User
		if session == nil {
			user = UserAnonymous
		} else {
			user = session.user
		}
		c.HTML(http.StatusOK, "hello.tmpl", gin.H{
			"displayname": user.displayName,
		})
	})
	r.Static("/static", "./static")
	r.GET("/login/begin", func(c *gin.Context) {
		username := c.Query("username")
		log.Println("Start login for '%s'", username)
		user, err := userstore.GetUser(username) // Find the user
		if err != nil {
			fmt.Println("Error on GetUser", err)
			c.JSON(400, gin.H{
				"error": err.Error(),
			})
		}

		options, session, err := webAuthn.BeginLogin(user)
		if err != nil {
			// Handle Error and return.
			log.Println("Error on login begin: %v", err)
			c.JSON(400, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.SetCookie(
			"authentication",
			authstore.StartSession(session),
			60*60*24*14, // Expires in 2 weeks
			"/",         // Valid for all paths
			"",
			false, // HTTPS only
			false, // allow JavaScript access to the cookie
		)

		//JSONResponse(w, options, http.StatusOK) // return the options generated
		c.JSON(200, options)
		// options.publicKey contain our registration options
	})
	r.POST("/login/finish", func(c *gin.Context) {
		username := c.Query("username")
		log.Println("Finish login for '%s'", username)
		user, err := userstore.GetUser(username) // Get the user
		if err != nil {
			log.Println("Failed to get user: %v", err)
			c.JSON(400, gin.H{
				"error": err.Error(),
			})
			return
		}
		cookie, err := c.Request.Cookie("authentication")
		// Get the session data stored from the function above
		if err != nil {
			log.Println("Failed to get authentication cookie: %v", err)
			c.JSON(400, gin.H{
				"error": err.Error(),
			})
			return
		}
		session, err := authstore.GetSession(cookie.Value)
		if err != nil {
			log.Println("Failed to get session: %v", err)
			c.JSON(400, gin.H{
				"error": err.Error(),
			})
			return
		}

		// in an actual implementation, we should perform additional checks on
		// the returned 'credential', i.e. check 'credential.Authenticator.CloneWarning'
		// and then increment the credentials counter
		credential, err := webAuthn.FinishLogin(user, *session, c.Request)
		if err != nil {
			// Handle Error and return.
			log.Println("Error on login finish: %v", err)
			return
		}

		if credential.Authenticator.CloneWarning {
			log.Println("cloned key detected")
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "cloned key detected",
			})
			return
		}

		// If login was successful, update the credential object
		// Pseudocode to update the user credential.
		user.UpdateCredential(credential)
		authstore.DeleteSession(cookie.Value)
		userstore.SaveUser(user)
		sessionstore.StartSession(c, user)

		c.JSON(200, "login success")
	})
	r.GET("/register/begin", func(c *gin.Context) {
		username := c.Query("username")
		displayname := c.Query("displayname")
		log.Println("Beginning registration for: ", username)
		// Get user
		user, err := userstore.GetUser(username)
		if err != nil {
			// User doesn't exist, create new user
			user = NewUser(username, displayname)
			log.Println("Created new user", user)
			userstore.SaveUser(user)
		}

		registerOptions := func(credCreationOpts *protocol.PublicKeyCredentialCreationOptions) {
			credCreationOpts.CredentialExcludeList = user.CredentialExcludeList()
		}

		// generate PublicKeyCredentialCreationOptions, session data
		options, session, err := webAuthn.BeginRegistration(user, registerOptions)
		if err != nil {
			c.JSON(500, gin.H{
				"error": err,
			})
			return
		}
		c.SetCookie(
			"registration",
			authstore.StartSession(session),
			3600, // age of the cookie in seconds
			"/",  // Valid for all paths
			"",
			false, // HTTPS only
			false, // allow JavaScript access to the cookie
		)

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
		user, err := userstore.GetUser(username) // Get the user
		if user == nil {
			log.Println("Unable to find user", username)
			c.JSON(500, gin.H{
				"error": "No such user",
			})
			return

		}

		cookie, err := c.Cookie("registration")
		if cookie == "" {
			log.Println("Registration cookie is empty")
			c.JSON(500, gin.H{
				"error": "empty registration cookie",
			})
			return
		}
		// Get the session data stored from the function above
		session, err := authstore.GetSession(cookie)
		if err != nil {
			log.Println("Failed to get session: %v", err)
			c.JSON(500, gin.H{
				"error": err,
			})
			return
		}

		credential, err := webAuthn.FinishRegistration(user, *session, c.Request)
		if err != nil {
			log.Println("Failed to finish registration: %v", err)
			c.JSON(500, gin.H{
				"error": err,
			})
			return
		}

		// If creation was successful, store the credential object
		// Pseudocode to add the user credential.
		user.AddCredential(*credential)
		userstore.SaveUser(user)
		sessionstore.StartSession(c, user)
		c.JSON(200, "Registration Success")
	})
	_ = r.Run()
}
