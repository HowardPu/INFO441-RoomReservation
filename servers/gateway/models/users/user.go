package users

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"net/mail"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

//gravatarBasePhotoURL is the base URL for Gravatar image requests.
//See https://id.gravatar.com/site/implement/images/ for details
const gravatarBasePhotoURL = "https://www.gravatar.com/avatar/"

//bcryptCost is the default bcrypt cost to use when hashing passwords
var bcryptCost = 13

//User represents a user account in the database
type User struct {
	ID        int64  `json:"id"`
	Email     string `json:"-"` //never JSON encoded/decoded
	PassHash  []byte `json:"-"` //never JSON encoded/decoded
	UserName  string `json:"userName"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	PhotoURL  string `json:"photoURL"`
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
	FirstName    string `json:"firstName"`
	LastName     string `json:"lastName"`
}

//Updates represents allowed updates to a user profile
type Updates struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
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
//PhotoURL and PassHash fields appropriately
func (nu *NewUser) ToUser() (*User, error) {
	//TODO: call Validate() to validate the NewUser and
	//return any validation errors that may occur.
	//if valid, create a new *User and set the fields
	//based on the field values in `nu`.
	//Leave the ID field as the zero-value; your Store
	//implementation will set that field to the DBMS-assigned
	//primary key value.
	//Set the PhotoURL field to the Gravatar PhotoURL
	//for the user's email address.
	//see https://en.gravatar.com/site/implement/hash/
	//and https://en.gravatar.com/site/implement/images/

	//TODO: also call .SetPassword() to set the PassHash
	//field of the User to a hash of the NewUser.Password

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
	imageURL := gravatarBasePhotoURL + hex.EncodeToString(h.Sum(nil))

	userRef.ID = 0
	userRef.PhotoURL = imageURL
	userRef.UserName = strings.TrimSpace(newUserDat.UserName)
	userRef.FirstName = strings.TrimSpace(newUserDat.FirstName)
	userRef.LastName = strings.TrimSpace(newUserDat.LastName)

	return &userRef, nil
}

//FullName returns the user's full name, in the form:
// "<FirstName> <LastName>"
//If either first or last name is an empty string, no
//space is put between the names. If both are missing,
//this returns an empty string
func (u *User) FullName() string {
	//TODO: implement according to comment above

	result := ""

	fName := (*u).FirstName
	lName := (*u).LastName

	if len(fName) != 0 {
		result += fName
	}

	if len(lName) != 0 {
		if len(result) != 0 {
			result += " "
		}
		result += lName
	}

	return result
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

//ApplyUpdates applies the updates to the user. An error
//is returned if the updates are invalid
func (u *User) ApplyUpdates(updates *Updates) error {
	//TODO: set the fields of `u` to the values of the related
	//field in the `updates` struct

	if updates == nil {
		return errors.New("Cannot update user")
	}

	updateDat := *updates
	fName := strings.TrimSpace(updateDat.FirstName)
	lName := strings.TrimSpace(updateDat.LastName)

	(*u).FirstName = fName
	(*u).LastName = lName

	return nil
}
