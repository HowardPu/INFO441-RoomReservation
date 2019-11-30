package users

import (
	"reflect"
	"testing"

	"golang.org/x/crypto/bcrypt"
)

//TODO: add tests for the various functions in user.go, as described in the assignment.
//use `go test -cover` to ensure that you are covering all or nearly all of your code paths.

func TestToUser(t *testing.T) {
	emailHash := "93ff63f8dcc70d71cc7e116ddbc84967"
	profileURL := gravatarBasePhotoURL + emailHash
	pwd := "111111"
	userName := "abc"
	fName := "a"
	lName := "bc"
	hashByte, _ := bcrypt.GenerateFromPassword([]byte(pwd), bcryptCost)

	cases := []struct {
		name        string
		newUser     NewUser
		expectError bool
		expectUser  *User
	}{
		{
			"Invalid",
			NewUser{"qqq", pwd, pwd, userName, fName, lName},
			true,
			nil,
		}, {
			"Regular",
			NewUser{"qqq@gmail.com", pwd, pwd, userName, fName, lName},
			false,
			&User{0, "qqq@gmail.com", hashByte, userName, fName, lName, profileURL},
		}, {
			"Profile: Email Test Case",
			NewUser{"QqQ@gmail.com", pwd, pwd, userName, fName, lName},
			false,
			&User{0, "QqQ@gmail.com", hashByte, userName, fName, lName, profileURL},
		}, {
			"Profile: Email Test Trim",
			NewUser{"   qqq@gmail.com   ", pwd, pwd, userName, fName, lName},
			false,
			&User{0, "qqq@gmail.com", hashByte, userName, fName, lName, profileURL},
		},
	}

	for _, c := range cases {
		user, err := (&(c.newUser)).ToUser()
		if !c.expectError && err != nil {
			t.Errorf("case %s: do not expect error but we get -%s-", c.name, err)
		}

		if c.expectError && err == nil {
			t.Errorf("case %s: do expect error but we do not get get", c.name)
		}

		if user == nil && c.expectUser != nil {
			t.Errorf("Does not return user")
		}

		if user != nil && c.expectUser == nil {
			t.Errorf("return user while nil is expected")
		}

		if user != nil && c.expectUser != nil {

			v := reflect.ValueOf(*user)

			m := reflect.ValueOf(*(c.expectUser))

			for i := 0; i < v.NumField(); i++ {
				key := string(v.Type().Field(i).Name)

				if key == "PassHash" {
					err := user.Authenticate(pwd)
					if err != nil {
						t.Errorf("false password setting")
					}
				} else {
					curVal := v.Field(i).String()
					expectVal := m.Field(i).String()
					if curVal != expectVal {
						t.Errorf("case %s: %s is wrong", c.name, key)
					}
				}
			}
		}
	}
}

func TestValidate(t *testing.T) {
	cases := []struct {
		name        string
		newUser     NewUser
		expectError bool
	}{
		{
			"Regular",
			NewUser{"qqq@hmail.com", "111111", "111111", "abc", "a", "bc"},
			false,
		}, {
			"Invalid Email",
			NewUser{"qqq", "111111", "111111", "abc", "a", "bc"},
			true,
		}, {
			"Invalid Password",
			NewUser{"qqq@hmail.com", "11111", "11111", "abc", "a", "bc"},
			true,
		}, {
			"Invalid Password Confirm",
			NewUser{"qqq@hmail.com", "111111", "222222", "abc", "a", "bc"},
			true,
		}, {
			"Empty User Name",
			NewUser{"qqq@hmail.com", "111111", "111111", "", "a", "bc"},
			true,
		}, {
			"User Name with Space",
			NewUser{"qqq@hmail.com", "111111", "111111", "a bc", "a", "bc"},
			true,
		},
	}

	for _, c := range cases {
		err := (&(c.newUser)).Validate()

		if !c.expectError && err != nil {
			t.Errorf("case %s: do not expect error but we get -%s-", c.name, err)
		}

		if c.expectError && err == nil {
			t.Errorf("case %s: do expect error but we do not get get", c.name)
		}
	}
}

func TestFullName(t *testing.T) {
	cases := []struct {
		name       string
		user       *User
		expectName string
	}{
		{
			"Regular",
			GetUser("a", "b"),
			"a b",
		}, {
			"No First Name",
			GetUser("", "b"),
			"b",
		}, {
			"No Last Name",
			GetUser("a", ""),
			"a",
		}, {
			"No Fiest No Last",
			GetUser("", ""),
			"",
		},
	}

	for _, c := range cases {
		fullName := (*(c.user)).FullName()

		if fullName != c.expectName {
			t.Errorf("case %s: expect -%s- but we get -%s-", c.name, c.expectName, fullName)
		}
	}
}

func TestAuthenticate(t *testing.T) {
	pwd := "hihciuajc"
	pwd2 := "dxwbjhcwb"

	user := User{}

	hashByte, _ := bcrypt.GenerateFromPassword([]byte(pwd), bcryptCost)

	user.PassHash = hashByte

	err := (&user).Authenticate(pwd)

	if err != nil {
		t.Errorf("case expect no error for authenticate but get one: %s", err)
	}

	err = (&user).Authenticate(pwd2)

	if err == nil {
		t.Errorf("case expect fail authentication but success: user passord: %s, given password: %s", pwd, pwd2)
	}
}

func TestSetPassword(t *testing.T) {
	pwd := "hihciuajc"

	user := User{}
	(&user).SetPassword(pwd)

	err := bcrypt.CompareHashAndPassword(user.PassHash, []byte(pwd))

	if err != nil {
		t.Errorf("case wrong password hashing: %s", err)
	}
}

func TestApplyUpdates(t *testing.T) {
	cases := []struct {
		name        string
		update 		*Updates
		user		*User
		expectError bool
		expectUser  *User
	}{
		{
			"Regular Case",
			&Updates{"a", "b"},
			GetTestUser1(),
			false,
			GetUser("a", "b"),
		}, {
			"Regular Case with Space",
			&Updates{"     a     ", "     b     "},
			GetTestUser1(),
			false,
			GetUser("a", "b"),
		}, {
			"Diff First Name, No Last Name Case",
			&Updates{"a", ""},
			GetTestUser1(),
			false,
			GetUser("a", ""),
		}, {
			"Diff Last Name, No First Name Case",
			&Updates{"", "b"},
			GetTestUser1(),
			false,
			GetUser("", "b"),
		}, {
			"Diff First Name with space, No Last Name with space",
			&Updates{"     a     ", "            "},
			GetTestUser1(),
			false,
			GetUser("a", ""),
		}, {
			"Diff Last Name with space, No First Name Case with space",
			&Updates{"           ", "    b      "},
			GetTestUser1(),
			false,
			GetUser("", "b"),
		}, {
			"No Firstname, No Lastname",
			&Updates{"", ""},
			GetTestUser1(),
			false,
			GetUser("", ""),
		}, {
			"Empty Firstname with Space, empty Lastname with Space",
			&Updates{" ", " "},
			GetTestUser1(),
			false,
			GetUser("", ""),
		}, {
			"Same Firstname, no lastname",
			&Updates{"c", ""},
			GetTestUser1(),
			false,
			GetUser("c", ""),
		}, {
			"Same Firstname with space, no lastname with space",
			&Updates{"    c    ", "     "},
			GetTestUser1(),
			false,
			GetUser("c", ""),
		}, {
			"No Firstname, Same Lastname",
			&Updates{"", "d"},
			GetTestUser1(),
			false,
			GetUser("", "d"),
		}, {
			"No Firstname with space, Same Lastname with space",
			&Updates{"     ", "     d      "},
			GetTestUser1(),
			false,
			GetUser("", "d"),
		}, {
			"Same Firstname, Same Lastname",
			&Updates{"c", "d"},
			GetTestUser1(),
			false,
			GetUser("c", "d"),
		}, {
			"Same Firstname with space, Same Lastname with space",
			&Updates{"     c     ", "     d     "},
			GetTestUser1(),
			false,
			GetUser("c", "d"),
		}, {
			"No Firstname, New LastName",
			nil,
			GetTestUser1(),
			true,
			GetTestUser1(),
		},
	}

	for _, c := range cases {
		err := c.user.ApplyUpdates(c.update)

		if err != nil && !c.expectError {
			t.Errorf("case %s: unexpected error but we get %s", c.name, err)
		}

		if err == nil && c.expectError {
			t.Errorf("case %s: expected error but we don't get", c.name)
		}

		expectFName := (*(c.expectUser)).FirstName
		expectLName := (*(c.expectUser)).LastName

		currentFName := (*(c.user)).FirstName
		currentLName := (*(c.user)).LastName

		if expectFName != currentFName {
			t.Errorf("case %s: expected first name %s but we update %s", c.name, expectFName, currentFName)
		}

		if expectLName != currentLName {
			t.Errorf("case %s: expected last name %s but we update %s", c.name, expectLName, currentLName)
		}
	}
}

func GetTestUser1() *User {
	return GetUser("c", "d")
}

func GetUser(fName string, lName string) *User {
	user := User{}
	user.FirstName = fName
	user.LastName = lName
	return &user
}
