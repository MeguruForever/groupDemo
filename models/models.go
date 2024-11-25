package models

import (
	"fmt"
	"github.com/gorilla/websocket"
	"sync"
	"time"
)

type Group struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type Message struct {
	GroupID string    `json:"group_id"`
	UserID  string    `json:"user_id"`
	Content string    `json:"content"`
	Time    time.Time `json:"time"`
}

type GroupPools struct {
	Mu    sync.Mutex
	Pools map[string]*SingleGroupPool
}
type SingleGroupPool struct {
	Connections    map[string]*websocket.Conn
	Broadcast      chan Message
	messageStore   []Message
	messageStoreMu sync.Mutex
}

func (pool *SingleGroupPool) BroadcastMessages() {
	for {
		message := <-pool.Broadcast
		for userId := range pool.Connections {
			conn := pool.Connections[userId]
			err := conn.WriteJSON(message)
			if err != nil {
				err := conn.Close()
				if err != nil {
					return
				}
				delete(pool.Connections, userId)
			}
		}
	}
}

func (pool *SingleGroupPool) SaveMessage(msg Message) {
	pool.messageStoreMu.Lock()
	defer pool.messageStoreMu.Unlock()
	pool.messageStore = append(pool.messageStore, msg)
	fmt.Println("backend save message", msg)
}
