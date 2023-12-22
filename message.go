package main

import (
	"encoding/json"
	"log"
)

type Message struct {
	Sender  string `json:"sender"`
	RoomID  string `json:"room_id"`
	Content string `json:"content"`
}

func (message *Message) encode() []byte {
	json, err := json.Marshal(message)
	if err != nil {
		log.Println(err)
	}

	return json
}
