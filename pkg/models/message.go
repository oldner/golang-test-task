package models

import "time"

type Message struct {
	Sender    string    `json:"sender"`
	Receiver  string    `json:"receiver"`
	Content   string    `json:"message"`
	Timestamp time.Time `json:"timestamp"`
}
