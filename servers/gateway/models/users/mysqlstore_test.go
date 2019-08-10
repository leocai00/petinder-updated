package users

import (
	"crypto/md5"
	"database/sql"
	"encoding/hex"
	"fmt"
	"log"
	"reflect"
	"strings"
	"testing"

	"gopkg.in/DATA-DOG/go-sqlmock.v1"
)

const getByID = `select \* from users where id=\?`
const getByEmail = `select \* from users where email=\?`
const getByUserName = `select \* from users where username=\?`
const insert = `insert into users \(email, pass_hash, user_name, first_name, last_name, photo_url\) values \(\?, \?, \?, \?, \?, \?\)`
const updateStatement = `update users set first_name=\?, last_name=\? where id=\?`
const delete = `delete from users where id=\?`

func Hash(email string) string {
	h := md5.New()
	h.Write([]byte(strings.TrimSpace(strings.ToLower(email))))
	return string(hex.EncodeToString(h.Sum(nil)))
}
func TestGetByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("error creating sql mock: %v", err)
	}

	defer db.Close()

	expectedUser := &User{1, "test@example.com", nil, "test", "firstName", "lastName", gravatarBasePhotoURL + Hash("test@example.com")}
	expectedUser.SetPassword("password")

	rows := sqlmock.NewRows([]string{"id", "email", "pass_hash", "user_name", "first_name", "last_name", "photo_url"})
	rows.AddRow(expectedUser.ID, expectedUser.Email, expectedUser.PassHash, expectedUser.UserName, expectedUser.FirstName, expectedUser.LastName, expectedUser.PhotoURL)

	mock.ExpectQuery(getByID).WithArgs(expectedUser.ID).WillReturnRows(rows)

	store := NewMySQLStore(db)

	u, err := store.GetByID(expectedUser.ID)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if err == nil && !reflect.DeepEqual(u, expectedUser) {
		t.Errorf("the user queried does not match with the expected user")
	}

	mock.ExpectQuery(getByID).WithArgs(-1).WillReturnError(sql.ErrNoRows)
	_, err = store.GetByID(-1)
	if err == nil {
		t.Errorf("expected error: %v, but recieved nil", sql.ErrNoRows)
	}

	mock.ExpectQuery(getByID).WithArgs(expectedUser.ID).WillReturnError(fmt.Errorf("error querying"))
	
	_, err = store.GetByID(expectedUser.ID)
	if err == nil {
		t.Errorf("expected error: %v, but recieved nil", fmt.Errorf("error querying"))
	}

	err = mock.ExpectationsWereMet()
	if err != nil {
		t.Errorf("unmet sqlmock expectations: %v", err)
	}
}

func TestGetByEmail(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("error creating sql mock: %v", err)
	}

	defer db.Close()

	expectedUser := &User{1, "test@example.com", nil, "test", "firstName", "lastName", gravatarBasePhotoURL + Hash("test@example.com")}
	expectedUser.SetPassword("password")

	rows := sqlmock.NewRows([]string{"id", "email", "pass_hash", "user_name", "first_name", "last_name", "photo_url"})
	rows.AddRow(expectedUser.ID, expectedUser.Email, expectedUser.PassHash, expectedUser.UserName, expectedUser.FirstName, expectedUser.LastName, expectedUser.PhotoURL)

	mock.ExpectQuery(getByEmail).WithArgs(expectedUser.Email).WillReturnRows(rows)

	store := NewMySQLStore(db)

	u, err := store.GetByEmail(expectedUser.Email)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if err == nil && !reflect.DeepEqual(u, expectedUser) {
		t.Errorf("the user queried does not match with the expected user")
	}

	mock.ExpectQuery(getByEmail).WithArgs("invalidemail").WillReturnError(sql.ErrNoRows)
	_, err = store.GetByEmail("invalidemail")
	if err == nil {
		t.Errorf("expected error: %v, but recieved nil", sql.ErrNoRows)
	}

	mock.ExpectQuery(getByEmail).WithArgs(expectedUser.Email).WillReturnError(fmt.Errorf("error querying"))

	_, err = store.GetByEmail(expectedUser.Email)
	if err == nil {
		t.Errorf("expected error: %v, but recieved nil", fmt.Errorf("error querying"))
	}

	err = mock.ExpectationsWereMet()
	if err != nil {
		t.Errorf("unmet sqlmock expectations: %v", err)
	}
}
func TestGetByUserName(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("error creating sql mock: %v", err)
	}

	defer db.Close()

	expectedUser := &User{1, "test@example.com", nil, "test", "firstName", "lastName", gravatarBasePhotoURL + Hash("test@example.com")}
	expectedUser.SetPassword("password")

	rows := sqlmock.NewRows([]string{"id", "email", "pass_hash", "user_name", "first_name", "last_name", "photo_url"})
	rows.AddRow(expectedUser.ID, expectedUser.Email, expectedUser.PassHash, expectedUser.UserName, expectedUser.FirstName, expectedUser.LastName, expectedUser.PhotoURL)

	mock.ExpectQuery(getByUserName).WithArgs(expectedUser.UserName).WillReturnRows(rows)

	store := NewMySQLStore(db)

	u, err := store.GetByUserName(expectedUser.UserName)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if err == nil && !reflect.DeepEqual(u, expectedUser) {
		t.Errorf("the user queried does not match with the expected user")
	}

	mock.ExpectQuery(getByUserName).WithArgs("nottest").WillReturnError(sql.ErrNoRows)
	_, err = store.GetByUserName("nottest")
	if err == nil {
		t.Errorf("expected error: %v, but recieved nil", sql.ErrNoRows)
	}

	mock.ExpectQuery(getByUserName).WithArgs(expectedUser.UserName).WillReturnError(fmt.Errorf("error querying"))
	_, err = store.GetByUserName(expectedUser.UserName)
	if err == nil {
		t.Errorf("expected error: %v, but recieved nil", fmt.Errorf("error querying"))
	}

	err = mock.ExpectationsWereMet()
	if err != nil {
		t.Errorf("unmet sqlmock expectations: %v", err)
	}
}
func TestInsert(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("error creating sql mock: %v", err)
	}

	defer db.Close()

	u := &User{
		Email:     "test@example.com",
		UserName:  "test",
		FirstName: "firstName",
		LastName:  "lastName",
		PhotoURL:  gravatarBasePhotoURL + Hash("test@example.com"),
	}
	u.SetPassword("password")

	store := NewMySQLStore(db)

	mock.ExpectExec(insert).WithArgs(u.Email, u.PassHash, u.UserName, u.FirstName, u.LastName, u.PhotoURL).WillReturnResult(sqlmock.NewResult(2, 1))

	user, err := store.Insert(u)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if err == nil && !reflect.DeepEqual(user, u) {
		t.Errorf("the user returned does not match with the input user")
	}

	invalidUser := &User{
		Email:     "invalidemail",
		UserName:  "test2",
		FirstName: "firstName",
		LastName:  "lastName",
		PhotoURL:  gravatarBasePhotoURL + Hash("invalidemail"),
	}

	mock.ExpectExec(insert).WithArgs(invalidUser.Email, invalidUser.PassHash, invalidUser.UserName, invalidUser.FirstName, invalidUser.LastName, invalidUser.PhotoURL).WillReturnError(fmt.Errorf("error inserting"))

	_, err = store.Insert(invalidUser)
	if err == nil {
		t.Errorf("expected error: %v but recieved nil", fmt.Errorf("error inserting"))
	}

	err = mock.ExpectationsWereMet()
	if err != nil {
		t.Errorf("unmet sqlmock expectations: %v", err)
	}
}

func TestUpdate(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatalf("error creating sql mock: %v", err)
	}

	defer db.Close()

	u := &User{1, "test@example.com", nil, "test", "firstName", "lastName", gravatarBasePhotoURL + Hash("test@example.com")}
	u.SetPassword("password")

	store := NewMySQLStore(db)
	rows := sqlmock.NewRows([]string{"id", "email", "pass_hash", "user_name", "first_name", "last_name", "photo_url"})
	update := &Updates{"firstName2", "lastName2"}
	rows.AddRow(u.ID, u.Email, u.PassHash, u.UserName, update.FirstName, update.LastName, u.PhotoURL)

	mock.ExpectExec(updateStatement).WithArgs(update.FirstName, update.LastName, u.ID).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectQuery(getByID).WithArgs(u.ID).WillReturnRows(rows)
	
	user, err := store.Update(u.ID, update)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if err == nil && reflect.DeepEqual(u, user) {
		t.Errorf("user did not update")
	}

	invalidUser := &User{
		ID:        -1,
		Email:     "invalidemail",
		UserName:  "test2",
		FirstName: "firstName",
		LastName:  "lastName",
		PhotoURL:  gravatarBasePhotoURL + Hash("invalidemail"),
	}

	mock.ExpectExec(updateStatement).WithArgs(update.FirstName, update.LastName, invalidUser.ID).WillReturnResult(sqlmock.NewResult(-1, 0))

	_, err = store.Update(invalidUser.ID, update)
	if err != ErrUserNotFound {
		t.Errorf("expected error: %v but recieved nil", ErrUserNotFound)
	}

	err = mock.ExpectationsWereMet()
	if err != nil {
		t.Errorf("unmet sqlmock expectations: %v", err)
	}

	noUser := &User{
		ID:        2,
		Email:     "nouser@example.com",
		UserName:  "nouser",
		FirstName: "firstName",
		LastName:  "lastName",
		PhotoURL:  gravatarBasePhotoURL + Hash("nouser@example.com"),
	}

	mock.ExpectExec(updateStatement).WithArgs(update.FirstName, update.LastName, noUser.ID).WillReturnResult(sqlmock.NewResult(2, 0))

	_, err = store.Update(noUser.ID, update)
	if err != ErrUserNotFound {
		t.Errorf("expected error: %v but recieved nil", ErrUserNotFound)
	}

	err = mock.ExpectationsWereMet()
	if err != nil {
		t.Errorf("unmet sqlmock expectations: %v", err)
	}
}

func TestDelete(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("error creating sql mock: %v", err)
	}

	defer db.Close()

	u := &User{1, "test@example.com", nil, "test", "firstName", "lastName", gravatarBasePhotoURL + Hash("test@example.com")}
	u.SetPassword("password")

	store := NewMySQLStore(db)

	mock.ExpectExec(delete).WithArgs(u.ID).WillReturnResult(sqlmock.NewResult(1, 1))

	err = store.Delete(u.ID)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	invalidUser := &User{
		ID:        -1,
		Email:     "invalidemail",
		UserName:  "test2",
		FirstName: "firstName",
		LastName:  "lastName",
		PhotoURL:  gravatarBasePhotoURL + Hash("invalidemail"),
	}

	mock.ExpectExec(delete).WithArgs(invalidUser.ID).WillReturnError(fmt.Errorf("error deleting"))

	err = store.Delete(invalidUser.ID)
	if err == nil {
		t.Errorf("expected error: %v but recieved nil", fmt.Errorf("error deleting"))
	}

	err = mock.ExpectationsWereMet()
	if err != nil {
		t.Errorf("unmet sqlmock expectations: %v", err)
	}
}