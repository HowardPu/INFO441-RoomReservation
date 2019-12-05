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
		SELECT TOP 1 U.userID, U.email, U.userName, U.passHash, UT.userTypeName
		FROM tblUser U JOIN tblUserType UT ON U.userTypeID = UT.userTypeID
		WHERE U.userID = ?`

	// search user in db
	userInfo, err := s.db.Query(insq, id)
	if err != nil {
		return nil, err
	}

	// scan the user, if error occurs, return it
	user := User{}
	userInfo.Next()
	scanErr := userInfo.Scan(&(user.ID), &(user.Email), &(user.UserName), &(user.PassHash), &(user.Type))
	userInfo.Close()

	if scanErr != nil {
		return nil, scanErr
	}

	// return user
	return &user, nil
}

func (s *MsSqlStore) GetByEmail(email string) (*User, error) {
	insq := `
	SELECT TOP 1 U.userID, U.email, U.userName, U.passHash, UT.userTypeName
	FROM tblUser U JOIN tblUserType UT ON U.userTypeID = UT.userTypeID
	WHERE U.email = ?`

	// search user by email
	userInfo, err := s.db.Query(insq, email)
	if err != nil {
		return nil, err
	}

	user := User{}
	userInfo.Next()
	// return error if occurs
	scanErr := userInfo.Scan(&(user.ID), &(user.Email), &(user.UserName), &(user.PassHash), &(user.Type))
	userInfo.Close()
	if scanErr != nil {
		return nil, scanErr
	}

	// return suer
	return &user, nil
}

func (s *MsSqlStore) GetByUserName(username string) (*User, error) {
	insq := `
	SELECT U.userID, U.email, U.userName, U.passHash, UT.userTypeName
	FROM tblUser U JOIN tblUserType UT ON U.userTypeID = UT.userTypeID
	WHERE U.userName = ?`

	// search user in db by username
	userInfo, err := s.db.Query(insq, username)
	if err != nil {
		return nil, err
	}

	// return any scan error,
	user := User{}
	userInfo.Next()
	scanErr := userInfo.Scan(&(user.ID), &(user.Email), &(user.UserName), &(user.PassHash), &(user.Type))
	userInfo.Close()
	if scanErr != nil {
		return nil, scanErr
	}

	// return thr user info
	return &user, nil
}

func (s *MsSqlStore) Insert(user *User) (*User, error) {

	// parse user info
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
			@userName = ?,
			@email = ?,
			@passHash = ?,
			@userTypeName = ?`

	// insert into sql db
	_, err := s.db.Exec(transaction,
		result.UserName,
		result.Email,
		result.PassHash,
		result.Type,
	)

	// return any error that exists
	if err != nil {
		return nil, err
	}

	// get the latest inserted user id
	latestInsertedSQL := `SELECT IDENT_CURRENT('tblUser')`
	lastestID, lastestErr := s.db.Query(latestInsertedSQL)

	if lastestErr != nil {
		return nil, lastestErr
	}

	// return any error that occurs
	lastestID.Next()
	scanErr := lastestID.Scan(&(result.ID))
	lastestID.Close()

	// set assigned id into user struct, and return it
	if scanErr != nil {
		return nil, scanErr
	}

	return &result, nil
}
