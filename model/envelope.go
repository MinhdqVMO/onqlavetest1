package model

import "time"

type MessageEnvelope struct {
	ID         string            `json:"ID"`
	Data       []byte            `json:"Data"`
	Timestamp  time.Time         `json:"Timestamp"`
	Attributes map[string]string `json:"Attributes"`
}
