package reservations

import (
	"database/sql"
	"fmt"
	"log"
	"sync"

	M "INFO441-RoomReservation/servers/reservation/models"

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

func (s *MsSqlStore) GetByUserName(username string) (*[]M.Room, error) {
	insq := `
		SELECT TOP 1 U.userID, U.email, U.userName, U.passHash, U.photoURL, UN.firstName, UN.lastName
		FROM tblUser U
		JOIN tblUserName UN ON U.userID = UN.userID
		WHERE UN.endDate IS NULL AND U.userName  = ?`

	rows, err := s.db.Query(insq, username)
	if err != nil {
		return nil, err
	}

	roomList := []M.Room{}

	for rows.Next() {
		var name, roomType, statusType string
		var id, floor, capacity int
		// Get values from row.
		scanErr := rows.Scan(&id, &name, &floor, &capacity, &roomType, &statusType)
		if scanErr != nil {
			log.Printf("Can't scan query result with error: %v \n", scanErr)
			return nil, scanErr
		}

		var room = M.Room{id, name, floor, capacity, roomType, statusType}
		roomList = append(roomList, room)
		fmt.Printf("ID: %d, Name: %s, Floor: %d, Capacity: %d, roomType: %s, statusType: %s, \n", id, name, floor, capacity, roomType, statusType)
	}

	return &roomList, nil
}
