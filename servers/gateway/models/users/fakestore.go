package users

import (
	"github.com/final-project-petinder/servers/gateway/indexes"
	"fmt"
)

// FakeStore struct
type FakeStore struct {
}

// NewFakeStore constructs and returns a pointer to a FakeStore struct
func NewFakeStore() *FakeStore {
	return &FakeStore{}
}

//GetByID returns the User with the given ID
func (s *FakeStore) GetByID(id int64) (*User, error) {
	// We can trigger an error by passing in an id of 10
	if id == 10 {
		return nil, fmt.Errorf("Error getting user with ID: %d", id)
	}
	u := &User{
		Email:     "fake@example.com",
		UserName:  "fake",
		FirstName: "Fake",
		LastName:  "Example",
	}
	return u, nil
}

//GetByEmail returns the User with the given email
func (s *FakeStore) GetByEmail(email string) (*User, error) {
	if email == "fake@uw.edu" {
		return nil, fmt.Errorf("Error getting user with Email: %v", email)
	}
	nu := &NewUser{
		Email:        "fake@gmail.com",
		Password:     "fakepassword",
		PasswordConf: "fakepassword",
		UserName:     "newfake",
		FirstName:    "fakeFirstName",
		LastName:     "fakeLastName",
	}
	return nu.ToUser()
}

//GetByUserName returns the User with the given Username
func (s *FakeStore) GetByUserName(username string) (*User, error) {
	if username == "pizza" {
		return nil, fmt.Errorf("Error getting user with username: %v", username)
	}
	u := &User{
		Email:     "burger@example.com",
		UserName:  "burger",
		FirstName: "Ham",
		LastName:  "Burger",
	}
	return u, nil
}

//Insert inserts the user into the database, and returns
//the newly-inserted User, complete with the DBMS-assigned ID
func (s *FakeStore) Insert(user *User) (*User, error) {
	u := &User{
		Email:     "dog@example.com",
		UserName:  "dog",
		FirstName: "Dog",
		LastName:  "Big",
	}
	if user == u {
		return nil, fmt.Errorf("Error inserting user: %v", u)
	}
	nuser := &User{
		Email:     "cat@example.com",
		UserName:  "cat",
		FirstName: "Cat",
		LastName:  "Small",
	}
	return nuser, nil
}

//Update applies UserUpdates to the given user ID
//and returns the newly-updated user
func (s *FakeStore) Update(id int64, updates *Updates) (*User, error) {
	if updates.FirstName == "" {
		return nil, fmt.Errorf("Error updating user with empty updates")
	}
	u := &User{
		Email:     "ugh@example.com",
		UserName:  "ugh",
		FirstName: "Ugh",
		LastName:  "Hgu",
	}
	return u, nil
}

//Delete deletes the user with the given ID
func (s *FakeStore) Delete(id int64) error {
	if id == 12 {
		return fmt.Errorf("Error deleting user with ID: %d", id)
	}
	return nil
}

func (s *FakeStore) Load(t *indexes.Trie) error {
	return nil
}
