package handlers

import (
	"errors"
	"log"
	"sync"

	"github.com/gorilla/websocket"
)

// connections: each auth token maps to its assigned websocket connection
// lock: lock for thread safe
// AuthChan: a channel listening a connection is created
//				for listening the client message
type SocketStore struct {
	Connections map[string]*websocket.Conn
	Lock        *sync.RWMutex
}

// constructs a new SocketStore
func NewSocketStore() *SocketStore {
	return &SocketStore{
		Connections: make(map[string]*websocket.Conn),
		Lock:        &sync.RWMutex{},
	}
}

func (s *SocketStore) AddNewConnection(auth string, conn *websocket.Conn) error {
	// begin lock for thread safe
	s.Lock.Lock()
	defer s.Lock.Unlock()

	// check the auth token not exist
	// if exist, throw error
	_, found := s.Connections[auth]
	if found {
		return errors.New("This token has already establish the websocket connection")
	}

	// set the auth token to this connection
	s.Connections[auth] = conn

	return nil
}

func (s *SocketStore) RemoveConnection(auth string) error {
	// lock the store
	s.Lock.Lock()
	defer s.Lock.Unlock()

	// found the connection
	// if the connection DNE, throw error
	_, found := s.Connections[auth]
	if !found {
		return errors.New("No Websocket Connection For this Auth Token")
	}

	// close and delete the connection from the map
	s.Connections[auth].Close()
	delete(s.Connections, auth)

	return nil
}

func (s *SocketStore) ProcessMessage(message []byte) error {
	// write the message for all connections
	// if there is err for writing the message
	// end the connection for that auth token
	for auth := range s.Connections {
		conn := s.Connections[auth]
		writeErr := conn.WriteMessage(websocket.TextMessage, message)
		if writeErr != nil {
			log.Println("Websocket Error")
			s.RemoveConnection(auth)
		}
	}

	return nil
}
