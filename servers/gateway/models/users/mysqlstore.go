package users

import (
	"github.com/final-project-petinder/servers/gateway/indexes"
	"strings"
	"fmt"
	"database/sql"
)

const (
	selectStatement = "select * from users"
	insertStatement = "insert into users (email, pass_hash, user_name, first_name, last_name, photo_url) values (?, ?, ?, ?, ?, ?)"
	updateStatement = "update users set first_name=?, last_name=? where id=?"
	deleteStatement = "delete from users where id=?"
	insertSignIn = "insert into sign_ins (userID, dateTime, ipAddress) values (?, ?, ?)"
)

// MySQLStore interface allowing you to abstract the sql client
type MySQLStore struct {
	Client *sql.DB
}

// NewMySQLStore constructs and returns a pointer to a MySQLStore struct
func NewMySQLStore(db *sql.DB) *MySQLStore {
	if db != nil {
		return &MySQLStore{
			Client: db,
		}
	}
	return nil
}

func (s *MySQLStore) baseGet(where string, by interface{}) (*User, error) {
	query := selectStatement + " where " + where + "=?"
	u := &User{}
	err := s.Client.QueryRow(query, by).Scan(&u.ID, &u.Email, &u.PassHash, &u.UserName, &u.FirstName, &u.LastName, &u.PhotoURL)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrUserNotFound
		}
		return nil, fmt.Errorf("Error scanning: %v", err)
	}
	return u, nil
}

//GetByID returns the User with the given ID
func (s *MySQLStore) GetByID(id int64) (*User, error) {
	// Process request and return correct output
	return s.baseGet("id", id)
}

//GetByEmail returns the User with the given email
func (s *MySQLStore) GetByEmail(email string) (*User, error) {
	// Process request and return correct output
	return s.baseGet("email", email)
}

//GetByUserName returns the User with the given Username
func (s *MySQLStore) GetByUserName(username string) (*User, error) {
	// Process request and return correct output
	return s.baseGet("username", username)
}

//Insert inserts the user into the database, and returns
//the newly-inserted User, complete with the DBMS-assigned ID
func (s *MySQLStore) Insert(user *User) (*User, error) {
	// Process request and return correct output
	r, err := s.Client.Exec(insertStatement, user.Email, user.PassHash, user.UserName, user.FirstName, user.LastName, user.PhotoURL)
	if err != nil {
		return nil, fmt.Errorf("Error executing: %v", err)
	}

	id, err := r.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("Error getting newly generated ID: %v", err)
	}

	user.ID = id
	return user, nil
}

//Update applies UserUpdates to the given user ID
//and returns the newly-updated user
func (s *MySQLStore) Update(id int64, updates *Updates) (*User, error) {
	// Process request and return correct output
	r, err := s.Client.Exec(updateStatement, updates.FirstName, updates.LastName, id)
	if err != nil {
		return nil, fmt.Errorf("Error updating: %v", err)
	}

	ra, err := r.RowsAffected()
	if err != nil {
		return nil, fmt.Errorf("Error getting the number of rows affected: %v", err)
	}

	if ra == 0 {
		return nil, ErrUserNotFound
	}
	
	return s.GetByID(id)
}

//Delete deletes the user with the given ID
func (s *MySQLStore) Delete(id int64) error {
	// Process request and return correct output
	_, err := s.Client.Exec(deleteStatement, id)
	if err != nil {
		return fmt.Errorf("Error deleting: %v", err)
	}
	return nil
}

//Load loads
func (s *MySQLStore) Load(t *indexes.Trie) error {
	r, err := s.Client.Query(selectStatement)
	if err != nil {
		return fmt.Errorf("Error selecting: %v", err)
	}
	defer r.Close()

	for r.Next() {
		user := &User{}
		err := r.Scan(&user.ID, &user.Email, &user.PassHash, &user.UserName, &user.FirstName, &user.LastName, &user.PhotoURL)
		if err != nil {
			if err == sql.ErrNoRows {
				return ErrUserNotFound
			}
			return fmt.Errorf("Error scanning: %v", err)
		}

		arr := strings.Split(strings.ToLower(user.UserName), " ")
		for i := 0; i < len(arr); i++ {
			t.Add(arr[i], user.ID)
		}

		arr = strings.Split(strings.ToLower(user.LastName), " ")
		for i := 0; i < len(arr); i++ {
			t.Add(arr[i], user.ID)
		}

		arr = strings.Split(strings.ToLower(user.FirstName), " ")
		for i := 0; i < len(arr); i++ {
			t.Add(arr[i], user.ID)
		}
	}
	return nil
}

// InsertSignIn inserts
func (s *MySQLStore) InsertSignIn(signIn *SignIn) (*SignIn, error) {
	res, err := s.Client.Exec(insertSignIn, signIn.UserID, signIn.DateTime, signIn.IPAddress)
	if err != nil {
		return nil, fmt.Errorf("Error inserting: %v", err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("Error getting last id: %v", err)
	}

	signIn.ID = id
	return signIn, nil
}