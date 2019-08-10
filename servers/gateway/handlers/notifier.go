package handlers

import (
	"log"
	"encoding/json"
	"github.com/streadway/amqp"
	"sync"
	"github.com/gorilla/websocket"
)

// SocketStore stores all the connections
type SocketStore struct {
	Connections map[int64]*websocket.Conn
	Lock        sync.Mutex
}

// NewSocketStore constructs a new SocketStore
func NewSocketStore() *SocketStore {
	return &SocketStore {
		Connections: make(map[int64]*websocket.Conn),
	}
}

// Msg is a struct to hold JSON objects received from the message queue
type Msg struct {
	Type    string      `json:"type,omitempty"`
	Action  interface{} `json:"action,omitempty"`
	UserIDs []int64     `json:"userIDs,omitempty"`
}

// InsertConnection is a thread-safe method for inserting a connection to the SocketStore
func (s *SocketStore) InsertConnection(conn *websocket.Conn, userID int64) {
	s.Lock.Lock()
	defer s.Lock.Unlock()
	// insert socket connection
	s.Connections[userID] = conn
	go s.ReadConnection(conn, userID)
}

// ReadConnection reads a connection
func (s *SocketStore) ReadConnection(conn *websocket.Conn, userID int64) {
	for {
		if _, _, err := conn.NextReader(); err != nil {
			conn.Close()
			delete(s.Connections, userID)
			break
		}
	}
}

// RemoveConnection is a thread-safe method for removing a connection
func (s *SocketStore) RemoveConnection(conn *websocket.Conn, userID int64) {
	s.Lock.Lock()
	// remove socket connection
	conn.Close()
	delete(s.Connections, userID)
	s.Lock.Unlock()
}

// ProcessMessages writes messages receievd at gateway
func (s *SocketStore) ProcessMessages (messages <-chan amqp.Delivery) {
	for message := range messages {
		s.Lock.Lock()
		Msg := &Msg{}
		err := json.Unmarshal(message.Body, Msg)
		if err != nil {
			log.Printf("Error decoding: %v", err)
			return
		}

		if len(Msg.UserIDs) == 0 {
			s.WriteToPublicConnections(message)
		} else {
			s.WriteToPrivateConnections(Msg.UserIDs, message)
		}

		message.Ack(false)
		s.Lock.Unlock()
	}
}

// WriteToPrivateConnections writes a message to a subset of connections
// (if the message is intended for a private channel)
func (s *SocketStore) WriteToPrivateConnections(users []int64, message amqp.Delivery) {
	for _, user := range users {
		if _, ok := s.Connections[user]; ok {
			err := s.Connections[user].WriteMessage(websocket.TextMessage, message.Body)
			if err != nil {
				s.Connections[user].Close()
				delete(s.Connections, user)
				return
			}
		}
	}
}

// WriteToPublicConnections writes a message to all live connections
// (if the message is posted on a public channel)
func (s *SocketStore) WriteToPublicConnections(message amqp.Delivery) {
	for user, connection := range s.Connections {
		err := connection.WriteMessage(websocket.TextMessage, message.Body)
		if err != nil {
			connection.Close()
			delete(s.Connections, user)
			return
		}
	}
}