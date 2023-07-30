package main

import (
	"fmt"
	"html/template"
	"net/http"
	"time"
)

var sessions = map[string]session{}

type session struct {
	user   User
	expiry time.Time
}

func (s session) isExpired() bool {
	return s.expiry.Before(time.Now())
}

type Credentials struct {
	Username string
	Email    string
	Password string
}

func auth(w http.ResponseWriter, r *http.Request) {
	var creds Credentials
	// Parse form to credentials
	creds.Username = r.FormValue("username")
	creds.Email = r.FormValue("email")
	creds.Password = r.FormValue("password")

	// Initialize Message for login validation
	msg := &Message{}

	if creds.Username != "" {
		// If the username field is not empty, use it for login
		msg.UsernameLogin = creds.Username
		msg.PasswordLogin = creds.Password
	} else if creds.Email != "" {
		// If the email field is not empty, use it for login
		msg.EmailLogin = creds.Email
		msg.PasswordLogin = creds.Password
	} else {
		// If both fields are empty, return an error message
		fmt.Println("Please enter a username or an email to login")
		return
	}

	// Validate login credentials
	if !msg.ValidateLogin() {
		fmt.Println("Invalid login credentials")
		// Prepare data to send to the template
		data := Data{LoggedIn: false, User: User{}, Message: msg, Posts: fetchAllPosts(database), Post: Post{}, Threads: fetchAllThreads(database), SigninModalOpen: "true", ScrollTo: ""}
		data.Posts = fillPosts(&data, fetchAllPosts(database))
		reverse(data.Posts)
		tmpl, _ := template.ParseFiles("static/template/index.html", "static/template/base.html")
		tmpl.Execute(w, data)
		return
	}

	fmt.Println("Logged in, preparing token...")
	setSessionToken(w, creds)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}


func registration(w http.ResponseWriter, r *http.Request) {
	var creds Credentials // users input
	// Parse form to credentials
	creds.Email = r.FormValue("email")
	creds.Username = r.FormValue("username")
	creds.Password = r.FormValue("password")
	msg := &Message{
		UsernameRegister: creds.Username,
		EmailRegister:    creds.Email,
		PasswordRegister: creds.Password,
	}
	if !msg.ValidateRegistration() {
		data := Data{LoggedIn: false, User: User{}, Message: msg, Posts: fetchAllPosts(database), Post: Post{}, Threads: fetchAllThreads(database), SignupModalOpen: "true", ScrollTo: ""}
		reverse(data.Posts)
		tmpl, _ := template.ParseFiles("static/template/index.html", "static/template/base.html")
		tmpl.Execute(w, data)
		return
	} else {
		p, _ := hashPassword(creds.Password)
		addUser(database, creds.Username, creds.Email, p)
		fmt.Println(creds.Username, creds.Email, p)
		setSessionToken(w, creds)
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}
func logout(w http.ResponseWriter, r *http.Request) {
	c, err := r.Cookie("session_token")
	if err != nil {
		if err == http.ErrNoCookie {
			// If the cookie is not set, return an unauthorized status
			fmt.Println("Unauthorized")
			http.Redirect(w, r, "/", http.StatusSeeOther)
		}
		// For any other type of error, return a bad request status
		fmt.Println("Bad Request")
	} else {
		sessionToken := c.Value
		delete(sessions, sessionToken)
		http.SetCookie(w, &http.Cookie{
			Name:    "session_token",
			Value:   "",
			Expires: time.Now(),
		})
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}
