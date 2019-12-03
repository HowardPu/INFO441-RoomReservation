package models

import (
	"database/sql"
	"errors"
	"time"

	_ "github.com/denisenkom/go-mssqldb"
)

/* Function Index:
1: GetByUserName?

2: ReleaseReservation(userName string, resID int) (string, error)
	remove specific a user's reservation of a room at given time

3: ReserveRoom(userName string, roomName string, beginTime int, duration int, reserveDate string) (int, error)
	reserve a room for a user at specific time

4: CheckOverlap(usedTime *[49]int, startTime int, duration int) error
	return error if there is an overlap of time interval

5: GetUsedTime(roomName string, date string) (*[49]int, error)
	return the reserved timespot of a room at a given time

6: GetRoomID(roomName string) (int, error)
	return the roomID by the room name

7: GetUserID(userName string) (int, error)
	return the userID by the user name

8: GetCurrentTime() string
	return the current data in format YY-MM-DD

9: AddRoom(roomName string, floor int, capacity int, roomType string, userName string) (int, error)
	Add the room into the database and return its id

10: DeleteRoom(roomName string, userName string) (int, error)
	Delete a specific room

10: AddEquipment(eqipName string, userName string) (int, error)
	Add an equipment into the database and return its id

11: AddEquipmentToRoom(equipName string, roomName string, userName string) (int, error)
	Add an equipment to a room and return its id

12: RemoveEquipmentToRoom(roomEquipID int, userName string) (int error)

12: GetLatestInsertedID(tableName string) (int, error)
	Get the id of latest inserted row in a table

13: GetReservationLists(userName string)
	Get the reservation list of a user

14: GetRoomLists(roomName string, capacity *int, floor *int, roomType string)
	Get the list of room which follows the search requirement:
		Note: "*" means any, nil for capacity and floor means any

15: GetEquipList(roomName string) ([]*Equipment, error)
	return the equipments in a room

16: GetIssues(roomName string, searchType string) ([]*Issue, error)
	get the issues
		serch type:
			"All": all issues,
			"RoomAll": all issues of a room
			"Unconfimed": unconfirmed issues of a room
			"Unsolved": unsolved issues of a room
			"Confirmed": confirmed issues of a room
			"Solved": solved issues of a room

17: AddIssue(issueBody string, roomName string) (int, error)
	add issue to a room

18: UpdateIssue(issueID int, updateType string, userName string) error
	update issue of a room
		updateType: confirm, solve

19: GetAllEquipment() ([]string, error)
	get all equipments

20: DeleteEquipment(euipName string, userName string) (equipID int, error)

21: get all reservations
What we need
*/

type ReservationStore struct {
	db *sql.DB
}

// Initialize a ReservationStore
func NewReservationStore(db *sql.DB) *ReservationStore {
	result := ReservationStore{}
	result.db = db
	return &result
}

func (s *ReservationStore) ReleaseReservation(userName string, resID int) (string, error) {
	if resID <= 0 {
		return "", errors.New("ID must be positive")
	}

	query := `
		SELECT TOP 1 roomName 
		FROM tblReservation RE JOIN tblRoom R ON R.roomID = RE.roomID
		WHERE RE.reservationID = ?`

	row, err := s.db.Query(query, resID)
	if err != nil {
		return "", errors.New("Cannot Execute Query")
	}

	if !row.Next() {
		return "", errors.New("Room Not Found")
	}

	roomName := ""

	scanErr := row.Scan(&roomName)

	if scanErr != nil {
		return "", errors.New("Cannot Scan Room")
	}

	transaction := ` 
		EXEC usp_releaseReservation
		@ResID = ?,
		@UserName = ?`

	_, tranErr := s.db.Exec(transaction, resID, userName)

	return roomName, tranErr
}

func (s *ReservationStore) ReserveRoom(userName string, roomName string, beginTime int,
	duration int, reserveDate string) (int, error) {

	if duration <= 0 {
		return -1, errors.New("Duration must be positive")
	}

	usedTime, err := s.GetUsedTime(roomName, reserveDate)

	if err != nil {
		return -1, err
	}

	overLapErr := s.CheckOverlap(usedTime, beginTime, duration)

	if overLapErr != nil {
		return -1, overLapErr
	}

	transaction := ` 
			EXEC usp_makeRoomReservation
			@userName = ?,
			@roomName = ?,
			@tranDate = ?,
			@reserveDate = ?,
			@beginTime = ?,
			@duration = ?
	`
	curDate := s.GetCurrentTime()

	_, tranErr := s.db.Exec(transaction,
		userName,
		roomName,
		curDate,
		reserveDate,
		beginTime,
		duration,
	)

	if tranErr != nil {
		return -1, tranErr
	}

	latestID, findErr := s.GetLatestInsertedID(`'tblReservation'`)

	return latestID, findErr
}

func (s *ReservationStore) CheckOverlap(usedTime *[49]int, startTime int, duration int) error {
	if duration <= 0 {
		return errors.New("Duration must be positive")
	}
	for i := startTime + 1; i <= startTime+duration; i++ {
		if usedTime[i] != 0 {
			return errors.New("There is a overlap")
		}
	}
	return nil
}

// name : `MGH441`
// date : `2010-1-1`
// return: 17, 18......42

// 8 - 9 fine
// 9 - 10 fine

func (s *ReservationStore) GetUsedTime(roomName string, date string) (*[49]int, error) {
	roomID, getRoomErr := s.GetRoomID(roomName)
	if getRoomErr != nil {
		return nil, getRoomErr
	}

	timeQuery := `SELECT beginTime, endTime FROM tblReservation WHERE reserveDate = ? AND roomID = ?`

	timeInfo, err := s.db.Query(timeQuery, date, roomID)
	defer timeInfo.Close()

	if err != nil {
		return nil, err
	}

	result := [49]int{}

	for timeInfo.Next() {
		var beginTime int
		var endTime int
		scanErr := timeInfo.Scan(&beginTime, &endTime)
		if scanErr != nil {
			return nil, scanErr
		}
		for i := beginTime + 1; i <= endTime; i++ {
			result[i] = result[i] + 1
		}
	}

	return &result, nil
}

func (s *ReservationStore) GetRoomID(roomName string) (int, error) {
	roomQuery := `SELECT TOP 1 roomID FROM tblRoom WHERE roomName = ?`

	roomInfo, err := s.db.Query(roomQuery, roomName)

	if err != nil {
		return -1, err
	}

	var roomID int

	roomInfo.Next()
	scanErr := roomInfo.Scan(&roomID)
	roomInfo.Close()
	if scanErr != nil {
		return -1, scanErr
	}

	return roomID, nil
}

func (s *ReservationStore) GetUserID(userName string) (int, error) {
	userQuery := `SELECT TOP 1 userID FROM tblUser WHERE userName = ?`

	userInfo, err := s.db.Query(userQuery, userName)

	if err != nil {
		return -1, err
	}

	var userID int

	userInfo.Next()
	scanErr := userInfo.Scan(&userID)
	userInfo.Close()
	if scanErr != nil {
		return -1, scanErr
	}

	return userID, nil
}

func (s *ReservationStore) GetCurrentTime() string {
	currentTime := time.Now()
	return currentTime.Format("2006-01-02")
}

func (s *ReservationStore) AddRoom(roomName string, floor int, capacity int, roomType string, userName string) (int, error) {
	transaction := `
		EXEC usp_createRoom
		@roomName = ?,
		@floor = ?,
		@capcity = ?,
		@roomTypeName = ?,
		@userName = ?
	`
	_, tranErr := s.db.Exec(transaction, roomName, floor, capacity, roomType, userName)

	if tranErr != nil {
		return -1, tranErr
	}

	latestID, findErr := s.GetLatestInsertedID(`'tblRoom'`)

	return latestID, findErr
}

// delete room

func (s *ReservationStore) DeleteRoom(roomName string, userName string) (int, error) {
	roomID, getErr := s.GetRoomID(roomName)
	if getErr != nil {
		return -1, errors.New("Cannot find room")
	}
	transaction := `
		EXEC usp_deleteRoom
		@roomName = ?,
		@userName = ?
	`

	_, tranErr := s.db.Exec(transaction, roomName, userName)

	return roomID, tranErr
}

func (s *ReservationStore) AddEquipment(eqipName string, userName string) (int, error) {
	transaction := `
			EXEC usp_addEquipment
			@equipName = ?,
			@userName = ?`
	_, tranErr := s.db.Exec(transaction, eqipName, userName)

	if tranErr != nil {
		return -1, tranErr
	}

	latestID, findErr := s.GetLatestInsertedID(`'tblEquipment'`)

	return latestID, findErr
}

func (s *ReservationStore) AddEquipmentToRoom(equipName string, roomName string, userName string) (int, error) {
	transaction := `
		EXEC usp_addEquipmentToRoom
		@equipName = ?,
		@roomName = ?,
		@addDate = ?,
		@userName = ?`

	curDate := s.GetCurrentTime()
	_, tranErr := s.db.Exec(transaction, equipName, roomName, curDate, userName)
	if tranErr != nil {
		return -1, tranErr
	}

	latestID, findErr := s.GetLatestInsertedID(`'tblEquipInRoom'`)

	return latestID, findErr
}

func (s *ReservationStore) DeleteEquipmentInRoom(roomEquipID int, userName string) error {
	date := s.GetCurrentTime()
	transaction := `
		EXEC usp_removeEquipmentInRoom
		@roomEquipID = ?,
		@userName = ?,
		@removeDate = ?`
	_, tranErr := s.db.Exec(transaction, roomEquipID, userName, date)
	return tranErr
}

func (s *ReservationStore) GetLatestInsertedID(tableName string) (int, error) {
	latestInsertedSQL := `SELECT IDENT_CURRENT(?)`

	lastestID, lastestErr := s.db.Query(latestInsertedSQL, tableName)

	if lastestErr != nil {
		return -1, lastestErr
	}

	var result int

	lastestID.Next()
	scanErr := lastestID.Scan(result)
	lastestID.Close()

	if scanErr != nil {
		return -1, scanErr
	}
	return result, nil
}

func (s *ReservationStore) GetReservationLists(userName string) ([]*Reservation, error) {
	result := []*Reservation{}
	query := `
		SELECT R.reservationID, R.tranDate, R.reserveDate, R.beginTime, R.endTime, RM.roomName, RT.roomTypeName
		FROM tblReservation R 
		JOIN tblUser U ON U.userID = R.userID
		JOIN tblRoom RM ON RM.roomID = R.roomID
		JOIN tblRoomType RT ON RT.roomTypeID = R.reservationID
		WHERE U.userName = ?
	`
	reservationInfo, err := s.db.Query(query, userName)

	if err != nil {
		return result, err
	}

	defer reservationInfo.Close()
	for reservationInfo.Next() {
		var id int
		var tranDate string
		var resDate string
		var beginTime int
		var endTime int
		var roomName string
		var roomType string
		scanErr := reservationInfo.Scan(&id, &tranDate, &resDate, &beginTime, &endTime, &roomName, &roomType)

		if scanErr != nil {
			return result, scanErr
		}

		res := Reservation{
			ID:          id,
			TranDate:    tranDate,
			ReserveDate: resDate,
			RoomName:    roomName,
			BeginTime:   beginTime,
			EndTime:     endTime,
			RoomType:    roomType,
		}

		result = append(result, &res)
	}

	return result, nil
}

func (s *ReservationStore) GetRoomLists(roomName string, capacity *int, floor *int, roomType string) ([]*Room, error) {
	query := `
		SELECT R.roomID, R.roomName, R.capacity, R.roomFloor, RT.roomTypeName 
		FROM tblRoom R JOIN tblRoomType RT ON R.roomTypeID = RT.roomTypeID
	`
	reservationInfo, err := s.db.Query(query)

	result := []*Room{}

	if err != nil {
		return result, err
	}

	defer reservationInfo.Close()
	for reservationInfo.Next() {
		var id int
		var rName string
		var cap int
		var flr int
		var rType string

		scanErr := reservationInfo.Scan(&id, &rName, &cap, &flr, &rType)

		if scanErr != nil {
			return result, scanErr
		}

		if roomName == "*" || roomName == rName {
			if capacity == nil || *capacity == cap {
				if floor == nil || *floor == flr {
					if rType == "*" || roomType == rType {
						curRoom := Room{
							ID:       id,
							RoomName: rName,
							Capacity: cap,
							Floor:    flr,
							RoomType: rType,
						}
						result = append(result, &curRoom)
					}
				}
			}
		}

	}
	return result, nil
}

func (s *ReservationStore) GetEquipList(roomName string) ([]*Equipment, error) {
	query := `
		SELECT E.equipName, ER.equipRoomID
		FROM tblEquipment E JOIN tblEquipInRoom ER ON ER.equipID = E.equipID
		JOIN tblRoom R ON R.roomID = ER.roomID
		WHERE R.roomName = '' AND ER.removeDate IS NULL`

	equipInfo, err := s.db.Query(query, roomName)
	defer equipInfo.Close()
	if err != nil {
		return nil, err
	}
	result := []*Equipment{}

	for equipInfo.Next() {
		var equipRoomID int
		var equip string
		scanErr := equipInfo.Scan(&equip, &equipRoomID)
		if scanErr != nil {
			return result, scanErr
		}
		curEquip := Equipment{
			Name:        equip,
			RoomEquipID: equipRoomID,
		}
		result = append(result, &curEquip)

	}
	return result, nil
}

func (s *ReservationStore) GetIssues(roomName string, searchType string) ([]*Issue, error) {
	queries := map[string]string{
		"All": `SELECT RI.roomIssueID, R.roomName, RI.roomIssueBody, RI.createDate 
		FROM tblRoomIssue RI JOIN tblRoom R ON R.roomID = RI.roomID`,
		"RoomAll": `SELECT RI.roomIssueID, R.roomName, RI.roomIssueBody, RI.createDate 
		FROM tblRoomIssue RI JOIN tblRoom R ON R.roomID = RI.roomID
		WHERE R.roomName = ?`,
		"Unconfimed": `SELECT RI.roomIssueID, R.roomName, RI.roomIssueBody, RI.createDate 
		FROM tblRoomIssue RI JOIN tblRoom R ON R.roomID = RI.roomID
		WHERE R.roomName = ? AND RI.confirmDate IS NULL`,
		"Unsolved": `SELECT RI.roomIssueID, R.roomName, RI.roomIssueBody, RI.createDate 
		FROM tblRoomIssue RI JOIN tblRoom R ON R.roomID = RI.roomID
		WHERE R.roomName = ? AND RI.solveDate IS NULL`,
		"Confirmed": `SELECT RI.roomIssueID, R.roomName, RI.roomIssueBody, RI.createDate 
		FROM tblRoomIssue RI JOIN tblRoom R ON R.roomID = RI.roomID
		WHERE R.roomName = ? AND RI.confirmDate IS NOT NULL`,
		"Solved": `SELECT RI.roomIssueID, R.roomName, RI.roomIssueBody, RI.createDate 
		FROM tblRoomIssue RI JOIN tblRoom R ON R.roomID = RI.roomID
		WHERE R.roomName = ? AND RI.solveDate IS NOT NULL`,
	}

	result := []*Issue{}
	query, found := queries[searchType]
	if !found {
		return result, nil
	}

	issueInfo, err := s.db.Query(query, roomName)
	defer issueInfo.Close()
	if err != nil {
		return nil, err
	}

	for issueInfo.Next() {
		var issueID int
		var roomName string
		var body string
		var createDate string
		scanErr := issueInfo.Scan(&issueID, &roomName, &body, &createDate)
		if scanErr != nil {
			return result, scanErr
		}
		currentIssue := Issue{
			ID:         issueID,
			RoomName:   roomName,
			Body:       body,
			CreateDate: createDate,
		}
		result = append(result, &currentIssue)
	}

	return result, nil
}

func (s *ReservationStore) AddIssue(issueBody string, roomName string) (int, error) {
	if len(issueBody) == 0 {
		return -1, errors.New("Plz write something for the issue")
	}
	date := s.GetCurrentTime()

	transaction := `
		EXEC usp_addIssue
		@roomName = ?,
		@roomIssue = ?,
		@issueDate = ?
	`
	_, tranErr := s.db.Exec(transaction, roomName, issueBody, date)
	if tranErr != nil {
		return -1, tranErr
	}

	latestID, findErr := s.GetLatestInsertedID(`'tblEquipInRoom'`)

	return latestID, findErr
}

func (s *ReservationStore) UpdateIssue(issueID int, updateType string, userName string) error {
	if issueID <= 0 {
		return errors.New("ID must be positive")
	}
	transactions := map[string]string{
		"confirm": `EXEC usp_confirmIssue
					@issueID = ?,
					@confirmDate = ?,
					@userName = ?`,
		"solve": `EXEC usp_solveIssue
					@issueID = ?,
					@solveDate = ?,
					@userName = ?`,
	}

	transaction, found := transactions[updateType]

	if !found {
		return errors.New("Update Method not found")
	}
	date := s.GetCurrentTime()
	_, tranErr := s.db.Exec(transaction, issueID, date, userName)

	return tranErr
}

func (s *ReservationStore) UpdateEquipName(oldname string, newname string, username string) error {
	transaction := `
		EXEC usp_updateEquipmentName
		@oldName = ?,
		@newName = ?,
		@userName = ?
	`

	_, tranErr := s.db.Exec(transaction, oldname, newname, username)

	return tranErr
}

func (s *ReservationStore) GetAllEquipment() ([]string, error) {
	query := `SELECT equipName FROM tblEquipment`
	equipInfo, err := s.db.Query(query)
	defer equipInfo.Close()
	if err != nil {
		return nil, err
	}
	result := []string{}
	for equipInfo.Next() {
		var equipName string
		scanErr := equipInfo.Scan(&equipName)
		if scanErr != nil {
			return result, scanErr
		}
		result = append(result, equipName)
	}
	return result, nil
}

func (s *ReservationStore) DeleteEquipment(equipName string, username string) error {
	transaction := `
		EXEC usp_deleteEquipment
		@equipName = ?,
		@userName = ?`
	_, tranErr := s.db.Exec(transaction, equipName, username)

	return tranErr
}
