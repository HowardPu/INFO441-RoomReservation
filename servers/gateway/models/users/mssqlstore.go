package users

import (
	"database/sql"
	"sync"

	_ "github.com/denisenkom/go-mssqldb"
)

var locker = &sync.RWMutex{}

type MsSqlStore struct {
	db *sql.DB
}

// Initialize a MSSqlStore
func NewMsSqlStore(db *sql.DB) *MsSqlStore {
	result := MsSqlStore{}
	result.db = db
	return &result
}

func (s *MsSqlStore) GetById(id int64) (*User, error) {
	insq := `
		SELECT TOP 1 U.userID, U.email, U.userName, U.passHash, U.photoURL, UN.firstName, UN.lastName
		FROM tblUser U
		JOIN tblUserName UN ON U.userID = UN.userID
		WHERE UN.endDate IS NULL AND U.userID  = ?`

	userInfo, err := s.db.Query(insq, id)
	if err != nil {
		return nil, err
	}

	user := User{}
	userInfo.Next()
	scanErr := userInfo.Scan(&(user.ID), &(user.Email), &(user.UserName), &(user.PassHash), &(user.PhotoURL), &(user.FirstName), &(user.LastName))
	userInfo.Close()

	if scanErr != nil {
		return nil, scanErr
	}

	return &user, nil
}

func (s *MsSqlStore) GetByEmail(email string) (*User, error) {
	insq := `
		SELECT TOP 1 U.userID, U.email, U.userName, U.passHash, U.photoURL, UN.firstName, UN.lastName
		FROM tblUser U
		JOIN tblUserName UN ON U.userID = UN.userID
		WHERE UN.endDate IS NULL AND U.email  = ?`
	userInfo, err := s.db.Query(insq, email)
	if err != nil {
		return nil, err
	}

	user := User{}
	userInfo.Next()
	scanErr := userInfo.Scan(&(user.ID), &(user.Email), &(user.UserName), &(user.PassHash), &(user.PhotoURL), &(user.FirstName), &(user.LastName))
	userInfo.Close()
	if scanErr != nil {
		return nil, scanErr
	}

	return &user, nil
}

func (s *MsSqlStore) GetByUserName(username string) (*User, error) {
	insq := `
		SELECT TOP 1 U.userID, U.email, U.userName, U.passHash, U.photoURL, UN.firstName, UN.lastName
		FROM tblUser U
		JOIN tblUserName UN ON U.userID = UN.userID
		WHERE UN.endDate IS NULL AND U.userName  = ?`
	userInfo, err := s.db.Query(insq, username)
	if err != nil {
		return nil, err
	}

	user := User{}
	userInfo.Next()
	scanErr := userInfo.Scan(&(user.ID), &(user.Email), &(user.UserName), &(user.PassHash), &(user.PhotoURL), &(user.FirstName), &(user.LastName))
	userInfo.Close()
	if scanErr != nil {
		return nil, scanErr
	}

	return &user, nil
}

func (s *MsSqlStore) Insert(user *User) (*User, error) {

	userDat := *user

	result := User{
		userDat.ID,
		userDat.Email,
		userDat.PassHash,
		userDat.UserName,
		userDat.FirstName,
		userDat.LastName,
		userDat.PhotoURL,
	}

	transaction :=
		`EXEC usp_addNewUser 
		@U_Name = ?, 
		@E_Mail = ?, 
		@P_Hash = ?, 
		@P_URL = ?, 
		@F_Name = ?, 
		@L_Name = ?`

	_, err := s.db.Exec(transaction,
		result.UserName,
		result.Email,
		result.PassHash,
		result.PhotoURL,
		result.FirstName,
		result.LastName)

	if err != nil {
		return nil, err
	}

	latestInsertedSQL := `SELECT IDENT_CURRENT('tblUser')`
	lastestID, lastestErr := s.db.Query(latestInsertedSQL)

	if lastestErr != nil {
		return nil, lastestErr
	}

	lastestID.Next()
	scanErr := lastestID.Scan(&(result.ID))
	lastestID.Close()

	if scanErr != nil {
		return nil, scanErr
	}

	return &result, nil
}

func (s *MsSqlStore) Update(id int64, updates *Updates) (*User, error) {

	curUser, err := s.GetById(id)
	if err != nil {
		return nil, err
	}

	updateErr := curUser.ApplyUpdates(updates)
	if updateErr != nil {
		return nil, updateErr
	}

	transaction :=
		`EXEC usp_updateUserName
		@u_Name = ?,
		@updateFName = ?,
		@updateLName = ?`

	_, tranErr := s.db.Exec(transaction, (*curUser).UserName, (*curUser).FirstName, (*curUser).LastName)

	if tranErr != nil {
		return nil, tranErr
	}

	return curUser, nil
}

func (s *MsSqlStore) Delete(id int64) error {

	user, getErr := s.GetById(id)

	if getErr != nil {
		return getErr
	}

	transaction := `
		EXEC usp_removeUser
		@U_Name = ?
	`
	_, tranErr := s.db.Exec(transaction, (*user).UserName)

	if tranErr != nil {
		return tranErr
	}

	return nil
}

func (s *MsSqlStore) AddSignInInfo(userName string, date string, ip string) error {

	transaction := `
		EXEC dbo.usp_addSignIn
		@U_Name = ?,
		@Date = ?,
		@IP = ?
	`
	_, tranErr := s.db.Exec(transaction, userName, date, ip)

	if tranErr != nil {
		return tranErr
	}

	return nil
}
