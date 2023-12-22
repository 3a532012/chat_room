package main

import (
	"encoding/json"
	"log"

	"github.com/gorilla/websocket"
)

type Client struct {
	ID         string
	Name       string
	Connection *websocket.Conn
	Message    chan []byte
	Room       *Room
	Exist      chan struct{}
}

func NewClient(id string, name string, conn *websocket.Conn, room *Room) *Client {
	client := &Client{
		ID:         id,
		Name:       name,
		Connection: conn,
		Message:    make(chan []byte),
		Room:       room,
		Exist:      make(chan struct{}),
	}
	go client.run()
	return client
}
func (client *Client) readPump() {
	for {
		_, jsonMessage, err := client.Connection.ReadMessage()
		if err != nil {
			log.Println("read message error")
			return
		}

		client.handleNewMessage(jsonMessage)
	}
}

func (client *Client) run() {
	defer func() {
		close(client.Message)
		close(client.Exist)
		client.Connection.Close()
		log.Printf("close client %s \n", client.ID)
	}()
	for {
		select {
		case m, ok := <-client.Message:
			if !ok {
				log.Println("close client message")
				return
			}
			client.sendMessage(m)
		case <-client.Exist:
			return
		}
	}
}
func (client *Client) handleNewMessage(message []byte) {
	client.Room.Broadcast <- message
}
func (client *Client) sendMessage(message []byte) {
	var m Message
	if err := json.Unmarshal(message, &m); err != nil {
		log.Println("wrong format in message")
		return
	}
	client.Connection.WriteJSON(m)
}
