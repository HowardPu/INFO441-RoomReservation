package handlers

import (
	"errors"
	"log"
	"sync"

	"github.com/gorilla/websocket"
)

type SocketStore struct {
	Connections map[string]*websocket.Conn
	Lock        *sync.RWMutex
	AuthChan    chan string
}

// constructs a new SocketStore
func NewSocketStore() *SocketStore {
	return &SocketStore{
		Connections: make(map[string]*websocket.Conn),
		Lock:        &sync.RWMutex{},
		AuthChan:    make(chan string),
	}
}

func (s *SocketStore) AddNewConnection(auth string, conn *websocket.Conn) error {
	s.Lock.Lock()
	defer s.Lock.Unlock()
	_, found := s.Connections[auth]
	if found {
		return errors.New("This token has already establish the websocket connection")
	}

	s.Connections[auth] = conn

	s.AuthChan <- auth

	return nil
}

func (s *SocketStore) RemoveConnection(auth string) error {
	s.Lock.Lock()
	defer s.Lock.Unlock()
	_, found := s.Connections[auth]
	if !found {
		return errors.New("No Websocket Connection For this Auth Token")
	}
	s.Connections[auth].Close()
	delete(s.Connections, auth)

	return nil
}

func (s *SocketStore) ProcessMessage(message []byte) error {
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
