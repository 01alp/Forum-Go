package main

import (
	"fmt"
	"mime/multipart"
	"regexp"
	"strings"
)

var rxEmail = regexp.MustCompile(`.+@.+\..+`)

type Message struct {
	EmailRegister    string
	UsernameRegister string
	PasswordRegister string

	EmailLogin       string
	PasswordLogin    string

	Threads          []string
	ImageHeader      *multipart.FileHeader
	Errors           map[string]string
}

func (msg *Message) ValidateLogin() bool {

	msg.Errors = make(map[string]string)

	user := fetchUserByEmail(database, msg.EmailLogin)
	if user == (User{}) {
		msg.Errors["EmailLogin"] = "No user with such email"
		return len(msg.Errors) == 0
	}
	
	if !checkPasswordHash(msg.PasswordLogin, user.Password) {
		// Handle wrong password
		msg.Errors["PasswordLogin"] = "Wrong password"
	}

	return len(msg.Errors) == 0
}

func (msg *Message) ValidateRegistration() bool {

	msg.Errors = make(map[string]string)

	// check if username is not present in database
	if fetchUserByUsername(database, msg.UsernameRegister) != (User{}) {
		msg.Errors["UsernameRegistration"] = "Username is present in database"
	}

	// check if email is correctly formated
	match := rxEmail.Match([]byte(msg.EmailRegister))

	if !match {
		msg.Errors["EmailRegistration"] = "Please enter a valid email address"
	}

	// check if email is present in database
	if fetchUserByEmail(database, msg.EmailRegister) != (User{}) {
		msg.Errors["EmailRegistration"] = "Email is present in database"
	}

	// check if username is correct length
	if len(msg.UsernameRegister) < 4 || len(msg.UsernameRegister) > 20 {
		msg.Errors["UsernameRegistration"] = "Invalid lenght of username"
		fmt.Println("Invalid lenght of Username")
	}

	// check if email is correct length
	if len(msg.EmailRegister) < 6 || len(msg.EmailRegister) > 500 {
		msg.Errors["EmailRegistration"] = "Invalid lenght of email"
	}

	// check if password is correct length
	if len(msg.PasswordRegister) < 4 || len(msg.PasswordRegister) > 20 {
		msg.Errors["PasswordRegistration"] = "Invalid lenght of password"
	}
	return len(msg.Errors) == 0
}

func (msg *Message) ValidateThreads() bool {

	msg.Errors = make(map[string]string)

	// check if at least one thread is chosen when creating new post
	if len(msg.Threads) == 0 {
		msg.Errors["Threads"] = "Choose at least one category"
		// fmt.Println("Smth went wrong")
		// fmt.Println(msg.Errors["Threads"])
	}
	return len(msg.Errors) == 0
}

func (msg *Message) ValidateComment() bool {

	msg.Errors = make(map[string]string)

	// if len(msg.Threads) == 0 {
	// 	msg.Errors["Threads"] = "Choose at least one category"
	// 	// fmt.Println("Smth went wrong")
	// 	// fmt.Println(msg.Errors["Threads"])
	// }
	return len(msg.Errors) == 0
}

func (msg *Message) ValidateImage() bool {
	msg.Errors = make(map[string]string)

	// Check if image is valid JPEG, SVG, PNG or GIF
	if !strings.HasSuffix(msg.ImageHeader.Filename, ".jpg") && !strings.HasSuffix(msg.ImageHeader.Filename, ".jpeg") && !strings.HasSuffix(msg.ImageHeader.Filename, ".svg") && !strings.HasSuffix(msg.ImageHeader.Filename, ".png") && !strings.HasSuffix(msg.ImageHeader.Filename, ".gif") {
		msg.Errors["Image"] = "Choose valid image type: JPEG, SVG, PNG or GIF"

	}
	// Check if the image is too big (max 20MB)
	if msg.ImageHeader.Size > 20000000 {
		msg.Errors["Image"] = "File size is too big. Choose image that is smaller than 20MB"
	}
	return len(msg.Errors) == 0

}
