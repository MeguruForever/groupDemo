package models

import "time"

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
