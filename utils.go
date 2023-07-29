package main

import (
	"fmt"
	"net/http"
	"regexp"
	"time"

	"github.com/gofrs/uuid"
	"golang.org/x/crypto/bcrypt"
	"reflect"
)

func setSessionToken(w http.ResponseWriter, creds Credentials) {
	// Create a new random session token
	uuid, _ := uuid.NewV4()
	sessionToken := (uuid).String()
	expiresAt := time.Now().Add(15 * 60 * time.Second)
	// Set the token in the session map, along with the session information
	dropOpenSession(fetchUserByEmail(database, creds.Email))
	sessions[sessionToken] = session{
		user:   fetchUserByEmail(database, creds.Email),
		expiry: expiresAt,
	}
	// http.SetCookie(w, &http.Cookie{
	// 	Name:    "session_token",
	// 	Value:   sessionToken,
	// 	Path:    "/",
	// 	Expires: expiresAt,
	// })
	setSessionCookie(w, sessionToken, expiresAt)
}

func dropOpenSession(user User) {
	for k, u := range sessions {
		if u.user == user {
			delete(sessions, k)
		}
	}
}

func setLastPage(w http.ResponseWriter, url string) {
	expiresAt := time.Now().Add(15 * 60 * time.Second)
	http.SetCookie(w, &http.Cookie{
		Name:    "last_page",
		Value:   url,
		Path:    "/",
		Expires: expiresAt,
	})
}

func setSessionCookie(w http.ResponseWriter, sessionToken string, expiresAt time.Time) {
	// expiresAt := time.Now().Add(15 * 60 * time.Second)

	http.SetCookie(w, &http.Cookie{
		Name:    "session_token",
		Value:   sessionToken,
		Path:    "/",
		Expires: expiresAt,
	})
}

func welcome(w http.ResponseWriter, r *http.Request) Data {

	output := Data{LoggedIn: false, User: User{}, Threads: fetchAllThreads(database)}

	c, err := r.Cookie("session_token")
	if err != nil {
		if err == http.ErrNoCookie {
			// If the cookie is not set, return an unauthorized status
			fmt.Println("Unauthorized")
			return output
		}
		// For any other type of error, return a bad request status
		fmt.Println("Bad Request")
		return output
	}
	sessionToken := c.Value
	userSession, exists := sessions[sessionToken]
	if !exists {
		// If the session token is not present in session map, return an unauthorized error
		fmt.Println("Unauthorized")
		return output
	}
	if userSession.isExpired() {
		_, ok := sessions[sessionToken]
		if ok {
			delete(sessions, sessionToken)
		}
		fmt.Println("Unauthorized")
		return output
	}
	// If the session is valid, return the welcome message to the user
	output.LoggedIn = true
	output.User = userSession.user
	output.Notifications = fetchActiveNotificationsByUserId(database, output.User.Id)
	// disable(output.Notifications)
	reverse(output.Notifications)
	fmt.Printf("\nWelcome %s!\n", userSession.user.Username)
	// // maybe think how to make session token change less frequently
	refresh(w, r)
	setSessionToken(w, Credentials{userSession.user.Username, userSession.user.Email,userSession.user.Password})
	fmt.Println("refreshing token")
	return output
}

func refresh(w http.ResponseWriter, r *http.Request) {
	c, err := r.Cookie("session_token")
	if err != nil {
		if err == http.ErrNoCookie {
			// If the cookie is not set, return an unauthorized status
			fmt.Println("Unauthorized")
		}
		// For any other type of error, return a bad request status
		fmt.Println("Bad Request")
	}
	sessionToken := c.Value
	userSession, exists := sessions[sessionToken]
	if !exists {
		// If the session token is not present in session map, return an unauthorized error
		fmt.Println("Unauthorized")
	}
	if userSession.isExpired() {
		delete(sessions, sessionToken)
		fmt.Println("Unauthorized")
	}
	// If the previous session is valid, create a new session token for the current user
	uuid, _ := uuid.NewV4()
	newSessionToken := (uuid).String()
	expiresAt := time.Now().Add(120 * time.Second)
	// Set the token in the session map, along with the user whom it represents
	sessions[newSessionToken] = session{
		user:   userSession.user,
		expiry: expiresAt,
	}
	// Delete the older session token
	delete(sessions, sessionToken)
	// Set the new token as the users `session_token` cookie
	// http.SetCookie(w, &http.Cookie{
	// 	Name:    "session_token",
	// 	Value:   newSessionToken,
	// 	Path:    "/",
	// 	Expires: time.Now().Add(15 * 60 * time.Second),
	// })
	setSessionCookie(w, sessionToken, expiresAt)
}

func isUnique(p Post, posts []Post) bool {
	for _, v := range posts {
		if p.Id == v.Id {
			return false
		}
	}
	return true
}

func filterByThread(posts []Post, thread string) []Post {
	var filtered []Post
	r, _ := regexp.Compile(thread + `($|,)`)
	for _, p := range posts {
		if r.MatchString(p.Thread) {
			filtered = append(filtered, p)
		}
	}
	return filtered
}

// add refresh func before every action

func fillPosts(data *Data, posts []Post) []Post {
	for i := 0; i < len(posts); i++ {
		posts[i].User = fetchUserById(database, posts[i].UserId)
		posts[i].Comments = fetchCommentsByPost(database, posts[i].Id)
		if data.LoggedIn {
			posts[i].UserReaction = fetchReactionByUserAndId(database, "postsReactions", data.User.Id, posts[i].Id).Value
		}
	}
	return posts
}

func reverse(s interface{}) {
	n := reflect.ValueOf(s).Len()
    swap := reflect.Swapper(s)
	for i, j := 0,  n-1; i < j; i, j = i+1, j-1 {
		swap(i, j)
	}
}

func contains(strings []string, s string) bool {
	for _, v := range strings {
		if v == s {
			return true
		}
	}
	return false
}

func hashPassword(password string) (string, error) {
    bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
    return string(bytes), err
}

func checkPasswordHash(password, hash string) bool {
    err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
    return err == nil
}

func disable(all []Notification) {
	for _,n := range all {
		err := disableNotificationByID(database, n.Id)
		fmt.Println("disable err", err)
	}
}

func deleteIndexFromSlice(s []Comment, i int) []Comment {
	copy(s[i:], s[i+1:])
	s[len(s)-1] = Comment{}
	return s[:len(s)-1]
}

func isValidCommentID(id int, comments []Comment) bool {
	for _, comment := range comments {
		if comment.Id == id {return true}
	}
	return false
}


func isValidPostID(id int, posts []Post) bool {
	for _, post := range posts {
		if post.Id == id {return true}
	}
	return false
}