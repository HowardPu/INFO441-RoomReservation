package users

import (
	"errors"
	"testing"

	M "github.com/DATA-DOG/go-sqlmock"
	"github.com/google/go-cmp/cmp"
)

func TestNewMsSQLStore(t *testing.T) {
	db, _, _ := M.New(M.QueryMatcherOption(M.QueryMatcherEqual))
	store := NewMsSqlStore(db)

	if store.db == nil {
		t.Errorf("Wrong Initialization, db should be initialized")
	}
	defer db.Close()
}

func TestGetByEmail(t *testing.T) {
	cases := []struct {
		name        string
		email       string
		expectError bool
		expectCol   []string
		expectRow   string
		expectUser  *User
	}{
		{
			"Regular",
			"aaa@a.com",
			false,
			[]string{"userID", "email", "userName", "passHash", "photoURL", "firstName", "lastName"},
			"1, aaa@a.com, Hedonist, xbujjjhb, dnwin, John, Smith",
			&User{1, "aaa@a.com", []byte("xbujjjhb"), "Hedonist", "John", "Smith", "dnwin"},
		}, {
			"No such EMail",
			"bbb@b.com",
			true,
			[]string{},
			"",
			nil,
		},
	}

	db, mock, _ := M.New(M.QueryMatcherOption(M.QueryMatcherEqual))
	defer db.Close()
	query := `
		SELECT TOP 1 U.userID, U.email, U.userName, U.passHash, U.photoURL, UN.firstName, UN.lastName 
		FROM tblUser U JOIN tblUserName UN ON U.userID = UN.userID 
		WHERE UN.endDate IS NULL AND U.email = ?`
	store := NewMsSqlStore(db)
	for _, c := range cases {
		mock.ExpectQuery(query).
			WithArgs(c.email).
			WillReturnRows(M.NewRows(c.expectCol).FromCSVString(c.expectRow))
		user, err := store.GetByEmail(c.email)

		if err != nil && !c.expectError {
			t.Errorf("No Error Expected but get one: %v", err)
		}

		if err == nil {
			if c.expectError {
				t.Errorf("Expect Error but Somehow Passed")
			} else {
				if !cmp.Equal(*user, *(c.expectUser)) {
					t.Errorf("user dat is wrong")
				}
			}
		}
	}
}

func TestGetByUserName(t *testing.T) {
	cases := []struct {
		name        string
		username    string
		expectError bool
		expectCol   []string
		expectRow   string
		expectUser  *User
	}{
		{
			"Regular",
			"Hedonist",
			false,
			[]string{"userID", "email", "userName", "passHash", "photoURL", "firstName", "lastName"},
			"1, aaa@a.com, Hedonist, xbujjjhb, dnwin, John, Smith",
			&User{1, "aaa@a.com", []byte("xbujjjhb"), "Hedonist", "John", "Smith", "dnwin"},
		}, {
			"No such EMail",
			"Justice",
			true,
			[]string{},
			"",
			nil,
		},
	}

	db, mock, _ := M.New(M.QueryMatcherOption(M.QueryMatcherEqual))
	defer db.Close()
	query := `
		SELECT TOP 1 U.userID, U.email, U.userName, U.passHash, U.photoURL, UN.firstName, UN.lastName 
		FROM tblUser U JOIN tblUserName UN ON U.userID = UN.userID 
		WHERE UN.endDate IS NULL AND U.userName = ?`
	store := NewMsSqlStore(db)
	for _, c := range cases {
		mock.ExpectQuery(query).
			WithArgs(c.username).
			WillReturnRows(M.NewRows(c.expectCol).FromCSVString(c.expectRow))
		user, err := store.GetByUserName(c.username)

		if err != nil && !c.expectError {
			t.Errorf("No Error Expected but get one: %v", err)
		}

		if err == nil {
			if c.expectError {
				t.Errorf("Expect Error but Somehow Passed")
			} else {
				if !cmp.Equal(*user, *(c.expectUser)) {
					t.Errorf("user dat is wrong")
				}
			}
		}
	}

}

func TestGetByID(t *testing.T) {

	cases := []struct {
		name        string
		id          int64
		expectError bool
		expectCol   []string
		expectRow   string
		expectUser  *User
	}{
		{
			"Regular",
			1,
			false,
			[]string{"userID", "email", "userName", "passHash", "photoURL", "firstName", "lastName"},
			"1, aaa@a.com, Hedonist, xbujjjhb, dnwin, John, Smith",
			&User{1, "aaa@a.com", []byte("xbujjjhb"), "Hedonist", "John", "Smith", "dnwin"},
		}, {
			"No such ID",
			2,
			true,
			[]string{},
			"",
			nil,
		},
	}

	db, mock, _ := M.New(M.QueryMatcherOption(M.QueryMatcherEqual))
	defer db.Close()
	query := `
		SELECT TOP 1 U.userID, U.email, U.userName, U.passHash, U.photoURL, UN.firstName, UN.lastName 
		FROM tblUser U JOIN tblUserName UN ON U.userID = UN.userID 
		WHERE UN.endDate IS NULL AND U.userID = ?`
	store := NewMsSqlStore(db)
	for _, c := range cases {
		mock.ExpectQuery(query).
			WithArgs(c.id).
			WillReturnRows(M.NewRows(c.expectCol).FromCSVString(c.expectRow))
		user, err := store.GetById(c.id)

		if err != nil && !c.expectError {
			t.Errorf("No Error Expected but get one: %v", err)
		}

		if err == nil {
			if c.expectError {
				t.Errorf("Expect Error but Somehow Passed")
			} else {
				if !cmp.Equal(*user, *(c.expectUser)) {
					t.Errorf("user dat is wrong")
				}
			}
		}
	}
}

func TestDelete(t *testing.T) {
	db, mock, _ := M.New(M.QueryMatcherOption(M.QueryMatcherEqual))
	userStore := NewMsSqlStore(db)
	query := `
		SELECT TOP 1 U.userID, U.email, U.userName, U.passHash, U.photoURL, UN.firstName, UN.lastName 
		FROM tblUser U JOIN tblUserName UN ON U.userID = UN.userID 
		WHERE UN.endDate IS NULL AND U.userID = ?`
	mock.ExpectQuery(query).
		WithArgs(1).
		WillReturnRows(M.NewRows([]string{"userID", "email", "userName", "passHash", "photoURL", "firstName", "lastName"}).
			FromCSVString("1, aaa@a.com, Hedonist, xbujjjhb, dnwin, John, Smith"))
	mock.ExpectExec(`EXEC usp_removeUser @U_Name = ?`).
		WithArgs("Hedonist").
		WillReturnResult(M.NewResult(0, 0))
	err := userStore.Delete(1)

	if err != nil {
		t.Errorf("No Error Expected but get 1 for delete: %v", err)
	}

	mock.ExpectQuery(query).
		WithArgs(2).
		WillReturnRows(M.NewRows([]string{""}).
			FromCSVString(""))

	err2 := userStore.Delete(2)

	if err2 == nil {
		t.Errorf("Expect error (no ID found) but no error")
	}

	mock.ExpectQuery(query).
		WithArgs(3).
		WillReturnRows(M.NewRows([]string{"userID", "email", "userName", "passHash", "photoURL", "firstName", "lastName"}).
			FromCSVString("3, aaa@a.com, Justice, xbujjjhb, dnwin, John, Smith"))
	mock.ExpectExec(`EXEC usp_removeUser @U_Name = ?`).
		WithArgs("Justice").
		WillReturnError(errors.New("Can't delete"))
	err3 := userStore.Delete(3)

	if err3 == nil {
		t.Errorf("Expect error (can't delete) but no error")
	}
}

// miss errors
func TestInsert(t *testing.T) {
	user1 := &User{-1, "aaa@a.com", []byte("xbujjjhb"), "Hedonist", "John", "Smith", "jhbjbhj"}
	db, mock, _ := M.New(M.QueryMatcherOption(M.QueryMatcherEqual))
	userStore := NewMsSqlStore(db)
	query := `EXEC usp_addNewUser 
				@U_Name = ?, 
				@E_Mail = ?, 
				@P_Hash = ?, 
				@P_URL = ?, 
				@F_Name = ?, 
				@L_Name = ?`

	mock.ExpectExec(query).
		WithArgs("Hedonist", "aaa@a.com", []byte("xbujjjhb"), "jhbjbhj", "John", "Smith").
		WillReturnResult(M.NewResult(1, 1))
	mock.ExpectQuery("SELECT IDENT_CURRENT('tblUser')").
		WillReturnRows(M.NewRows([]string{"userID"}).FromCSVString("1"))

	newUser1, user1Err := userStore.Insert(user1)

	if user1Err != nil {
		t.Errorf("No Error Expected but get 1 for first user: %v", user1Err)
	}

	if (*newUser1).ID != 1 {
		t.Errorf("Wrong ID for the first user: %d", (*newUser1).ID)
	}

	user2 := &User{-1, "bbb@b.com", []byte("b"), "b", "b", "b", "b"}

	mock.ExpectExec(query).
		WithArgs("b", "bbb@b.com", []byte("b"), "b", "b", "b").
		WillReturnError(errors.New("DB is down"))
	_, user1Err2 := userStore.Insert(user2)

	if user1Err2 == nil {
		t.Errorf("Expect Error (DB is down) but somehow passed")
	}

	db2, mock2, _ := M.New(M.QueryMatcherOption(M.QueryMatcherEqual))
	defer db2.Close()
	userstore2 := NewMsSqlStore(db2)

	mock2.ExpectExec(query).
		WithArgs("Hedonist", "aaa@a.com", []byte("xbujjjhb"), "jhbjbhj", "John", "Smith").
		WillReturnResult(M.NewResult(1, 1))

	mock2.ExpectQuery("SELECT IDENT_CURRENT('tblUser')").
		WillReturnError(errors.New("DB is down"))

	_, user1ErrBusinessRule := userstore2.Insert(user1)

	if user1ErrBusinessRule == nil {
		t.Errorf("Expect Error (Business Rule Violation) but somehow passed")
	}

	db3, mock3, _ := M.New(M.QueryMatcherOption(M.QueryMatcherEqual))
	defer db3.Close()
	userstore3 := NewMsSqlStore(db3)

	mock3.ExpectExec(query).
		WithArgs("Hedonist", "aaa@a.com", []byte("xbujjjhb"), "jhbjbhj", "John", "Smith").
		WillReturnResult(M.NewResult(1, 1))

	mock3.ExpectQuery("SELECT IDENT_CURRENT('tblUser')").WillReturnRows(M.NewRows([]string{}).FromCSVString(""))

	_, user1Rollback := userstore3.Insert(user1)

	if user1Rollback == nil {
		t.Errorf("Expect Error (Rollback) but somehow passed")
	}

	defer db.Close()
}

func TestUpdate(t *testing.T) {
	db, mock, _ := M.New(M.QueryMatcherOption(M.QueryMatcherEqual))
	userStore := NewMsSqlStore(db)
	update1 := &Updates{"a", "b"}
	update2 := &Updates{"", ""}

	query := `
		SELECT TOP 1 U.userID, U.email, U.userName, U.passHash, U.photoURL, UN.firstName, UN.lastName 
		FROM tblUser U JOIN tblUserName UN ON U.userID = UN.userID 
		WHERE UN.endDate IS NULL AND U.userID = ?`

	mock.ExpectQuery(query).
		WithArgs(1).
		WillReturnRows(M.NewRows([]string{}).FromCSVString(""))

	_, err1 := userStore.Update(1, update1)

	if err1 == nil {
		t.Errorf("Expect error (no ID found) but no error")
	}

	_, err2 := userStore.Update(1, update2)
	if err2 == nil {
		t.Errorf("Expect error (update violation) but no error")
	}

	mock.ExpectQuery(query).
		WithArgs(3).
		WillReturnRows(M.NewRows([]string{"userID", "email", "userName", "passHash", "photoURL", "firstName", "lastName"}).
			FromCSVString("3, aaa@a.com, Justice, xbujjjhb, dnwin, John, Smith"))
	mock.ExpectExec(`EXEC usp_updateUserName @u_Name = ?, @updateFName = ?, @updateLName = ?`).
		WithArgs("Justice", "a", "b").
		WillReturnError(errors.New("Can't Update"))

	_, err3 := userStore.Update(3, update1)

	if err3 == nil {
		t.Errorf("Expect error (can't update) but no error")
	}

	mock.ExpectQuery(query).
		WithArgs(4).
		WillReturnRows(M.NewRows([]string{"userID", "email", "userName", "passHash", "photoURL", "firstName", "lastName"}).
			FromCSVString("4, aaa@a.com, Howard, xbujjjhb, dnwin, John, Smith"))
	mock.ExpectExec(`EXEC usp_updateUserName @u_Name = ?, @updateFName = ?, @updateLName = ?`).
		WithArgs("Howard", "a", "b").
		WillReturnResult(M.NewResult(0, 0))

	user, err4 := userStore.Update(4, update1)

	if err4 != nil {
		t.Errorf("Expect no error but get one: %v", err4)
	}

	expectUser := User{4, "aaa@a.com", []byte("xbujjjhb"), "Howard", "a", "b", "dnwin"}

	if !cmp.Equal(*user, expectUser) {
		t.Errorf("Update Failed")
	}
	defer db.Close()
}
