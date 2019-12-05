package users

import (
	"crypto/md5"
	"fmt"
	"io"
	"net/mail"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

//bcryptCost is the default bcrypt cost to use when hashing passwords
var bcryptCost = 13
var RegUserType = "Normal"

//User represents a user account in the database
type User struct {
	ID       int64  `json:"userID"`
	Email    string `json:""`  //never JSON encoded/decoded
	PassHash []byte `json:"-"` //never JSON encoded/decoded
	UserName string `json:"userName"`
	Type     string `json:"userType"`
}

//Credentials represents user sign-in credentials
type Credentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

//NewUser represents a new user signing up for an account
type NewUser struct {
	Email        string `json:"email"`
	Password     string `json:"password"`
	PasswordConf string `json:"passwordConf"`
	UserName     string `json:"userName"`
}

//Validate validates the new user and returns an error if
//any of the validation rules fail, or nil if its valid
func (nu *NewUser) Validate() error {
	//TODO: validate the new user according to these rules:
	//- Email field must be a valid email address (hint: see mail.ParseAddress)
	//- Password must be at least 6 characters
	//- Password and PasswordConf must match
	//- UserName must be non-zero length and may not contain spaces
	//use fmt.Errorf() to generate appropriate error messages if
	//the new user doesn't pass one of the validation rules

	userDat := *nu
	_, err := mail.ParseAddress(userDat.Email)

	if err != nil {
		return fmt.Errorf("Invalid Email Address")
	}

	if len(userDat.Password) < 6 {
		return fmt.Errorf("Insecured Password (at least 6 characters)")
	}

	if userDat.Password != userDat.PasswordConf {
		return fmt.Errorf("Password does not match Password Confirm")
	}

	if len(userDat.UserName) == 0 || strings.IndexAny(userDat.UserName, " ") != -1 {
		return fmt.Errorf("Please type something and have no space for username")
	}

	return nil
}

//ToUser converts the NewUser to a User, setting the
//User type and PassHash fields appropriately
func (nu *NewUser) ToUser() (*User, error) {
	err := (*nu).Validate()

	if err != nil {
		return nil, err
	}

	userRef := User{}

	newUserDat := *nu

	errPass := userRef.SetPassword(newUserDat.Password)

	if errPass != nil {
		return nil, errPass
	}

	h := md5.New()
	emailTrim := strings.TrimSpace(newUserDat.Email)
	userRef.Email = emailTrim
	emailLower := strings.ToLower(emailTrim)
	io.WriteString(h, emailLower)

	userRef.ID = 0
	userRef.UserName = strings.TrimSpace(newUserDat.UserName)
	userRef.Type = RegUserType
	return &userRef, nil
}

//SetPassword hashes the password and stores it in the PassHash field
func (u *User) SetPassword(password string) error {
	//TODO: use the bcrypt package to generate a new hash of the password
	//https://godoc.org/golang.org/x/crypto/bcrypt

	hashByte, err := bcrypt.GenerateFromPassword([]byte(password), bcryptCost)

	if err != nil {
		return err
	}

	(*u).PassHash = hashByte
	return nil
}

//Authenticate compares the plaintext password against the stored hash
//and returns an error if they don't match, or nil if they do
func (u *User) Authenticate(password string) error {
	//TODO: use the bcrypt package to compare the supplied
	//password with the stored PassHash
	//https://godoc.org/golang.org/x/crypto/bcrypt
	time.Sleep(1 * time.Second)
	err := bcrypt.CompareHashAndPassword((*u).PassHash, []byte(password))

	return err
}
