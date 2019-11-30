package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"sync"

	"github.com/gorilla/websocket"
)

type ConnectionHoder struct {
	Connection *websocket.Conn
	UserID     int64
}

type SocketStore struct {
	Connections map[string]*ConnectionHoder
	userAuthMap map[int64][]string
	Lock        *sync.RWMutex
	AuthChan    chan string
}

// constructs a new SocketStore
func NewSocketStore() *SocketStore {
	return &SocketStore{
		Connections: make(map[string]*ConnectionHoder),
		userAuthMap: make(map[int64][]string),
		Lock:        &sync.RWMutex{},
		AuthChan:    make(chan string),
	}
}

func (s *SocketStore) AddNewConnection(auth string, userID int64, conn *websocket.Conn) error {
	s.Lock.Lock()
	defer s.Lock.Unlock()
	_, found := s.Connections[auth]
	if found {
		return errors.New("This token has already establish the websocket connection")
	}
	holder := ConnectionHoder{
		Connection: conn,
		UserID:     userID,
	}
	s.Connections[auth] = &holder
	_, found = s.userAuthMap[userID]
	if !found {
		s.userAuthMap[userID] = []string{}
	}

	authTokens, _ := s.userAuthMap[userID]

	for _, token := range authTokens {
		if token == auth {
			return errors.New("Token and User ID relationship has been established, maybe forget to delete")
		}
	}
	s.AuthChan <- auth
	s.userAuthMap[userID] = append(authTokens, auth)

	return nil
}

func (s *SocketStore) RemoveConnection(auth string) error {
	s.Lock.Lock()
	defer s.Lock.Unlock()
	_, found := s.Connections[auth]
	if !found {
		return errors.New("No Websocket Connection For this Auth Token")
	}
	userID := s.Connections[auth].UserID
	s.Connections[auth].Connection.Close()
	delete(s.Connections, auth)

	authTokens, _ := s.userAuthMap[userID]

	if authTokens == nil {
		return fmt.Errorf("No Relationship between user ID %d and Auth Token %v", userID, auth)
	}

	for i, token := range authTokens {
		if token == auth {
			authTokens = append(authTokens[:i], authTokens[i+1:]...)
		}
	}
	s.userAuthMap[userID] = authTokens
	return nil
}

type Message struct {
	Type    string      `json:"type,omitempty"`
	Action  interface{} `json:"action,omitempty"`
	UserIDs []int64     `json:"userIDs,omitempty"`
}

func (s *SocketStore) ProcessMessage(message []byte) error {
	msgJSON := Message{}
	marshalErr := json.Unmarshal(message, &msgJSON)
	if marshalErr != nil {
		log.Printf("Unmarshal Error %v \n", marshalErr)
		return marshalErr
	}

	userIDs := msgJSON.UserIDs
	if userIDs == nil {
		return errors.New("No User Id Found")
	}

	if len(userIDs) == 0 {
		for user := range s.userAuthMap {
			userIDs = append(userIDs, user)
		}
	}
	for i := 0; i < len(userIDs); i++ {
		curID := userIDs[i]
		authTokens, foundTokens := s.userAuthMap[curID]
		if !foundTokens || len(authTokens) == 0 {
			log.Printf("No auth token for this user %d", curID)
			return errors.New("No auth token for the user")
		} else {
			log.Printf("current auth token: %v", authTokens)
			for j := 0; j < len(authTokens); j++ {
				authToken := authTokens[j]
				connection, connFound := s.Connections[authToken]
				if !connFound {
					log.Printf("No auth token for this connection %v", authToken)
					return errors.New("No auth token for the connection")
				} else {
					conn := connection.Connection
					writeErr := conn.WriteMessage(websocket.TextMessage, message)
					if writeErr != nil {
						log.Println("Websocket Error")
						s.RemoveConnection(authToken)
					}
				}
			}
		}
	}
	return nil
}
