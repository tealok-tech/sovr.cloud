package main

import (
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/go-webauthn/webauthn/protocol"
	"github.com/go-webauthn/webauthn/webauthn"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
)

func main() {
	config := CreateConfig()
	wconfig := &webauthn.Config{
		RPDisplayName: config.RelyingPartyDisplayName, // Display Name for your site
		RPID:          config.RelyingPartyID,          // Generally the FQDN for your site
		RPOrigins:     config.RelyingPartyOrigins,     // The origin URLs allowed for WebAuthn requests
	}
	var webAuthn *webauthn.WebAuthn
	var err error
	if webAuthn, err = webauthn.New(wconfig); err != nil {
		log.Println(err)
		return
	}

	authstore := CreateAuthstore()
	userstore := CreateUserstore()

	r := gin.Default()
	r.ForwardedByClientIP = true
	r.SetTrustedProxies(config.TrustedProxies)
	store := cookie.NewStore([]byte(config.SessionSecret))
	r.Use(sessions.Sessions("session", store))

	//r.LoadHTMLGlob("templates/*")
	templates, err := loadTemplates(config.TemplatesDir)
	if err != nil {
		panic(err)
	}
	r.Static("/static", "./static")
	r.GET("/", func(c *gin.Context) {
		session := sessions.Default(c)
		username := session.Get("username")
		var user *User
		if username == nil {
			c.Redirect(http.StatusFound, "/login")
			return
		}
		user = userstore.GetUser(username.(string))
		if user == nil {
			c.Redirect(http.StatusFound, "/login")
			return
		}
		/*c.HTML(http.StatusOK, "index.tmpl", gin.H{
			"User": user,
		})*/
		err = renderTemplate(c.Writer, templates, "index.tmpl", gin.H{
			"User": user,
		})
		if err != nil {
			log.Println("Failed to render template:", err)
		}
	})
	r.GET("/login", func(c *gin.Context) {
		err = renderTemplate(c.Writer, templates, "login.tmpl", gin.H{})
		if err != nil {
			log.Println("Failed to render template:", err)
		}
	})
	r.GET("/login/begin", func(c *gin.Context) {
		username := c.Query("username")
		log.Println("Start login for:", username)
		user := userstore.GetUser(username)
		if user == nil {
			c.JSON(404, gin.H{
				"error": "User does not exist",
			})
			return
		}

		options, session, err := webAuthn.BeginLogin(user)
		if err != nil {
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

		// options.publicKey contain our registration options
		c.JSON(200, options)
	})
	r.POST("/login/finish", func(c *gin.Context) {
		username := c.Query("username")
		log.Println("Finish login for:", username)
		user := userstore.GetUser(username)
		if user == nil {
			c.JSON(400, gin.H{
				"error": "User does not exist",
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
			log.Println("Error on login finish:", err)
			c.JSON(400, gin.H{
				"error": err.Error(),
			})
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
		user_session := sessions.Default(c)
		user_session.Set("username", user.Name)
		err = user_session.Save()
		if err != nil {
			log.Println("Failed to save session", err)
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err,
			})
			return
		}
		c.Header("Location", "/")
		c.Writer.WriteHeader(http.StatusNoContent)
	})
	r.POST("/logout", func(c *gin.Context) {
		user_session := sessions.Default(c)
		user_session.Clear()
		err = user_session.Save()
		if err != nil {
			log.Println("Failed to save session", err)
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err,
			})
			return
		}
		c.Redirect(http.StatusFound, "/")
	})
	r.GET("/register/begin", func(c *gin.Context) {
		username := c.Query("username")
		displayname := c.Query("displayname")
		log.Println("Beginning registration for:", username)
		user := userstore.GetUser(username)
		if user == nil {
			// User doesn't exist, create new user
			user = NewUser(username, displayname)
			userstore.SaveUser(user)
		}

		registerOptions := func(credCreationOpts *protocol.PublicKeyCredentialCreationOptions) {
			credCreationOpts.CredentialExcludeList = user.CredentialExcludeList()
		}

		log.Println("Generating registration options and session data for:", username)
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
		user := userstore.GetUser(username) // Get the user
		if user == nil {
			log.Println("Unable to find user to finish registration", username)
			c.JSON(400, gin.H{
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
		user_session := sessions.Default(c)
		user_session.Set("username", user.Name)
		err = user_session.Save()
		if err != nil {
			log.Println("Failed to save session", err)
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err,
			})
		}
		c.Header("Location", "/")
		c.Writer.WriteHeader(http.StatusNoContent)
	})
	r.GET("/trace", func(c *gin.Context) {
		forwarded := c.ClientIP()
		fmt.Println(c.Request.Host, c.Request.RemoteAddr, c.Request.RequestURI)
		c.JSON(200, gin.H{
			"address": forwarded,
		})
		return
	})
	_ = r.Run()
}

func loadTemplates(templatesDir string) (map[string]*template.Template, error) {
	log.Println("Using templates dir", templatesDir)
	result := make(map[string]*template.Template)
	pattern := templatesDir + "/layouts/*.tmpl"
	layouts, err := filepath.Glob(pattern)
	log.Println("Found", len(layouts), "layouts in", pattern)
	if err != nil {
		log.Fatal(err)
		return result, err
	}
	pattern = templatesDir + "/includes/*.tmpl"
	includes, err := filepath.Glob(pattern)
	log.Println("Found", len(includes), "includes in", pattern)
	if err != nil {
		log.Fatal(err)
		return result, err
	}

	for _, layout := range layouts {
		files := append(includes, layout)
		base := filepath.Base(layout)
		result[base] = template.Must(template.ParseFiles(files...))
		log.Printf("Added template '%s' with files %s", base, files)
	}
	return result, nil
}

func renderTemplate(w http.ResponseWriter, templates map[string]*template.Template, name string, data map[string]interface{}) error {
	tmpl, ok := templates[name]
	if !ok {
		return fmt.Errorf("The template %s does not exist.", name)
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	return tmpl.ExecuteTemplate(w, "base", data)
}
