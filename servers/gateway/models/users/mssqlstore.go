package users

import (
	"database/sql"

	_ "github.com/denisenkom/go-mssqldb"
)

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
		SELECT U.userID, U.email, U.userName, U.passHash, UT.userTypeName
		FROM tblUser U JOIN tblUserType UT ON U.userTypeID = UT.userTypeID
		WHERE U.userID = ?`

	userInfo, err := s.db.Query(insq, id)
	if err != nil {
		return nil, err
	}

	user := User{}
	userInfo.Next()
	scanErr := userInfo.Scan(&(user.ID), &(user.Email), &(user.UserName), &(user.PassHash), &(user.Type))
	userInfo.Close()

	if scanErr != nil {
		return nil, scanErr
	}

	return &user, nil
}

func (s *MsSqlStore) GetByEmail(email string) (*User, error) {
	insq := `
	SELECT U.userID, U.email, U.userName, U.passHash, UT.userTypeName
	FROM tblUser U JOIN tblUserType UT ON U.userTypeID = UT.userTypeID
	WHERE U.email = ?`
	userInfo, err := s.db.Query(insq, email)
	if err != nil {
		return nil, err
	}

	user := User{}
	userInfo.Next()
	scanErr := userInfo.Scan(&(user.ID), &(user.Email), &(user.UserName), &(user.PassHash), &(user.Type))
	userInfo.Close()
	if scanErr != nil {
		return nil, scanErr
	}

	return &user, nil
}

func (s *MsSqlStore) GetByUserName(username string) (*User, error) {
	insq := `
	SELECT U.userID, U.email, U.userName, U.passHash, UT.userTypeName
	FROM tblUser U JOIN tblUserType UT ON U.userTypeID = UT.userTypeID
	WHERE U.userName = ?`
	userInfo, err := s.db.Query(insq, username)
	if err != nil {
		return nil, err
	}

	user := User{}
	userInfo.Next()
	scanErr := userInfo.Scan(&(user.ID), &(user.Email), &(user.UserName), &(user.PassHash), &(user.Type))
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
		userDat.Type,
	}

	transaction :=
		`EXEC usp_addNewUser 
		@userName VARCHAR(64),
		@email VARCHAR(64),
		@passHash BINARY(60),
		@userTypeName VARCHAR(32)`

	_, err := s.db.Exec(transaction,
		result.UserName,
		result.Email,
		result.PassHash,
		result.Type,
	)

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
