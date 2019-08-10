package users

import (
	"testing"
)

//TODO: add tests for the various functions in user.go, as described in the assignment.
//use `go test -cover` to ensure that you are covering all or nearly all of your code paths.
func TestValidate(t *testing.T) {
	cases := []struct {
		name          string
		hint          string
		nu            *NewUser
		expectedError bool
	}{
		{
			"Valid New User",
			"This is a valid new user, so this should work",
			&NewUser{"test@example.com", "password", "password", "test", "firstName", "lastName"},
			false,
		},
		{
			"Invalid Email",
			"Remember to validate email",
			&NewUser{"invalidemail", "password", "password", "test", "firstName", "lastName"},
			true,
		},
		{
			"Password Too Short",
			"Remember to check the length of the password",
			&NewUser{"test@example.com", "blah", "blah", "test", "firstName", "lastName"},
			true,
		},
		{
			"Password Not Match",
			"Remember to check the matching of the passwords",
			&NewUser{"test@example.com", "password", "notmatch", "test", "firstName", "lastName"},
			true,
		},
		{
			"Empty User Name",
			"Remember to check the length of the user name",
			&NewUser{"test@example.com", "password", "password", "", "firstName", "lastName"},
			true,
		},
		{
			"User Name Contains Spaces",
			"Remember to check for spaces in the user name",
			&NewUser{"test@example.com", "password", "password", "te st", "firstName", "lastName"},
			true,
		},
	}

	for _, c := range cases {
		err := c.nu.Validate()
		if err != nil && !c.expectedError {
			t.Errorf("case %s: unexpected error %v\nHINT: %s", c.name, err, c.hint)
		}
		if c.expectedError && err == nil {
			t.Errorf("case %s: expected error but didn't get one\nHINT: %s", c.name, c.hint)
		}
	}
}

func TestToUser(t *testing.T) {
	cases := []struct {
		name          string
		hint          string
		nu            *NewUser
		expectedUser  *User
		expectedError bool
	}{
		{
			"Valid Conversion",
			"This is a valid conversion, so this should work",
			&NewUser{"test@example.com", "password", "password", "test", "firstName", "lastName"},
			&User{0, "test@example.com", nil, "test", "firstName", "lastName", "https://www.gravatar.com/avatar/55502F40DC8B7C769880B10874ABC9D0"},
			false,
		},
		{
			"Email Contains Upper Case Letters and Spaces",
			"Remember to format the email",
			&NewUser{"Test@example.com ", "password", "password", "test", "firstName", "lastName"},
			&User{0, "Test@example.com ", nil, "test", "firstName", "lastName", "https://www.gravatar.com/avatar/55502F40DC8B7C769880B10874ABC9D0"},
			false,
		},
	}

	for _, c := range cases {
		c.expectedUser.SetPassword(c.nu.Password)
		_, err := c.nu.ToUser()
		if err != nil && !c.expectedError {
			t.Errorf("case %s: unexpected error %v\nHINT: %s", c.name, err, c.hint)
		}
		if c.expectedError && err == nil {
			t.Errorf("case %s: expected error but didn't get one\nHINT: %s", c.name, c.hint)
		}
	}
}

func TestFullName(t *testing.T) {
	cases := []struct {
		name             string
		hint             string
		FirstName        string
		LastName	     string
		expectedFullName string
	}{
		{
			"Complete Full Name",
			"This is a complete full name, so this should work",
			"firstName",
			"lastName",
			"firstName lastName",
		},
		{
			"Both Names Empty",
			"Remember to check both first and last names",
			"",
			"",
			"",
		},
		{
			"First Name Empty",
			"Remember to check first name",
			"",
			"lastName",
			"lastName",
		},
		{
			"Last Name Empty",
			"Remember to check last name",
			"firstName",
			"",
			"firstName",
		},
	}

	for _, c := range cases {		
		u := &User{}
		u.FirstName = c.FirstName
		u.LastName = c.LastName
		if u.FullName() != c.expectedFullName {
			t.Errorf("case %s: incorrect full name\nHINT: %s", c.name, c.hint)
		}
	}
}

func TestAuthenticate(t *testing.T) {
	cases := []struct {
		name          string
		hint          string
		passsword     string
		expectedError bool
	}{
		{
			"Correct Password",
			"The password is correct, so this should work",
			"password",
			false,
		},
		{
			"Incorrect Password",
			"Remember to return the error if the passwords don't match",
			"notpassword",
			true,
		},
		{
			"Empty Password",
			"Remember to return the error if the passwords don't match",
			"",
			true,
		},
	}

	for _, c := range cases {		
		u := &User{}
		u.SetPassword("password")
		
		err := u.Authenticate(c.passsword)
		if err != nil && !c.expectedError {
			t.Errorf("case %s: unexpected error %v\nHINT: %s", c.name, err, c.hint)
		}
		if c.expectedError && err == nil {
			t.Errorf("case %s: expected error but didn't get one\nHINT: %s", c.name, c.hint)
		}
	}
}

func TestApplyUpdates(t *testing.T) {
	cases := []struct {
		name          string
		hint          string
		updates       *Updates
		expectedError bool
	}{
		{
			"Valid Updates",
			"This is a valid update, so this should work",
			&Updates{"newFirstName", "newLastName"},
			false,
		},
		{
			"Invalid First Name",
			"Remember to check if the new first name is valid",
			&Updates{"", "newLastName"},
			true,
		},
		{
			"Invalid Last Name",
			"Remember to check if the new last name is valid",
			&Updates{"newFirstName", ""},
			true,
		},
	}

	for _, c := range cases {		
		u := &User{}
		u.FirstName = "firstName"
		u.LastName = "lastName"

		err := u.ApplyUpdates(c.updates)
		if err != nil && !c.expectedError {
			t.Errorf("case %s: unexpected error %v\nHINT: %s", c.name, err, c.hint)
		}
		if c.expectedError && err == nil {
			t.Errorf("case %s: expected error but didn't get one\nHINT: %s", c.name, c.hint)
		}
	}
}