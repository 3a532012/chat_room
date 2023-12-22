package main

import (
	"log"

	"github.com/google/uuid"
)

type Room struct {
	ID         string
	Name       string
	Register   chan *Client
	Unregister chan *Client
	Clients    map[*Client]bool
	IsPrivate  bool
	Broadcast  chan []byte
	Exist      chan struct{}
}

func NewRoom(name string, isPrivate bool) *Room {
	r := &Room{
		ID:         uuid.New().String(),
		Name:       name,
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Clients:    make(map[*Client]bool),
		IsPrivate:  isPrivate,
		Broadcast:  make(chan []byte),
		Exist:      make(chan struct{}),
	}
	go r.run()
	return r
}

func (r *Room) run() {
	defer func() {
		close(r.Register)
		close(r.Unregister)
		close(r.Broadcast)
		close(r.Exist)
		log.Printf("close room %s \n", r.Name)
	}()
	for {
		select {
		case client, ok := <-r.Register:
			if !ok {
				log.Println("close room register")
				return
			}
			r.addClient(client)
		case client, ok := <-r.Unregister:
			if !ok {
				log.Println("close room unregister")
				return
			}
			r.removeClient(client)
		case broadcast, ok := <-r.Broadcast:
			if !ok {
				log.Println("close room broadcast")
				return
			}
			r.broadcastMessage(broadcast)
		case <-r.Exist:
			return
		}
	}
}
func (r *Room) broadcastMessage(message []byte) {
	for client := range r.Clients {
		client.Message <- message
	}
}
func (r *Room) addClient(client *Client) {
	if _, ok := r.Clients[client]; ok {
		log.Printf("duplicate register clien in room :%s \n", client.ID)
	} else {
		log.Printf("add client %s in room %s \n", client.ID, r.Name)
		r.Clients[client] = true
	}
}
func (r *Room) removeClient(client *Client) {
	if _, ok := r.Clients[client]; ok {
		delete(r.Clients, client)
		log.Printf("remove client %s from room %s \n", client.ID, r.Name)
	} else {
		log.Printf("unable unregister client in room :%s \n", client.ID)
	}
}
