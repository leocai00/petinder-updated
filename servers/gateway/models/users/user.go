package users

import (
	"time"
	"encoding/hex"
	"golang.org/x/crypto/bcrypt"
	"crypto/md5"
	"strings"
	"fmt"
	"net/mail"
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

// SignIn is a struct to collect sign in
type SignIn struct {
	ID        int64     `json:"id"`
	UserID    int64     `json:"userID"`
	DateTime  time.Time `json:"dateTime"`
	IPAddress string    `json:"ipAddress"`
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
	_, err := mail.ParseAddress(nu.Email)
	if err != nil {
		return err
	}
	
	if len(nu.Password) < 6 {		
		return fmt.Errorf("password is not at least 6 characters")
	}

	if nu.Password != nu.PasswordConf {
		return fmt.Errorf("passwords do not match")
	}

	if len(nu.UserName) == 0 {
		return fmt.Errorf("user name cannot be zero length")
	}

	if strings.Contains(nu.UserName, " ") {
		return fmt.Errorf("user name may not contain spaces")
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
	err := nu.Validate()
	if err != nil {
		return nil, err
	}

	newUser := &User{}
	newUser.Email = nu.Email

	err = newUser.SetPassword(nu.Password)
	if err != nil {
		return nil, err
	}

	newUser.UserName = nu.UserName
	newUser.FirstName = nu.FirstName
	newUser.LastName = nu.LastName

	b := md5.Sum([]byte(strings.ToLower(strings.TrimSpace(nu.Email))))
	newUser.PhotoURL = gravatarBasePhotoURL + hex.EncodeToString(b[:])
	
	return newUser, nil
}

//FullName returns the user's full name, in the form:
// "<FirstName> <LastName>"
//If either first or last name is an empty string, no
//space is put between the names. If both are missing,
//this returns an empty string
func (u *User) FullName() string {
	//TODO: implement according to comment above
	if u.FirstName == "" {
		return u.LastName
	}
	if u.LastName == "" {
		return u.FirstName
	}
	return u.FirstName + " " + u.LastName
}

//SetPassword hashes the password and stores it in the PassHash field
func (u *User) SetPassword(password string) error {
	//TODO: use the bcrypt package to generate a new hash of the password
	//https://godoc.org/golang.org/x/crypto/bcrypt
	ps, err := bcrypt.GenerateFromPassword([]byte(password), bcryptCost)
	if err != nil {
		return err
	}
	u.PassHash = ps
	return nil
}

//Authenticate compares the plaintext password against the stored hash
//and returns an error if they don't match, or nil if they do
func (u *User) Authenticate(password string) error {
	//TODO: use the bcrypt package to compare the supplied
	//password with the stored PassHash
	//https://godoc.org/golang.org/x/crypto/bcrypt
	return bcrypt.CompareHashAndPassword(u.PassHash, []byte(password))
}

//ApplyUpdates applies the updates to the user. An error
//is returned if the updates are invalid
func (u *User) ApplyUpdates(updates *Updates) error {
	//TODO: set the fields of `u` to the values of the related
	//field in the `updates` struct
	if updates.FirstName == "" {
		return fmt.Errorf("cannot update to an empty first name")
	}
	u.FirstName = updates.FirstName

	if updates.LastName == "" {
		return fmt.Errorf("cannot update to an empty last name")
	}
	u.LastName = updates.LastName
	return nil
}
