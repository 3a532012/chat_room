package main

import (
	"log"
	"sort"
	"strings"
)

type WebsocketMainInstance struct {
	Clients          map[string]*Client
	Rooms            map[string]*Room
	RegisterClient   chan *Client
	UnregisterClient chan *Client
	RegisterRoom     chan *Room
	UnregisterRoom   chan *Room
	Exist            chan struct{}
}

func NewWebsocketMainInstance() *WebsocketMainInstance {
	w := &WebsocketMainInstance{
		Clients:          make(map[string]*Client),
		Rooms:            make(map[string]*Room),
		RegisterClient:   make(chan *Client),
		UnregisterClient: make(chan *Client),
		RegisterRoom:     make(chan *Room),
		UnregisterRoom:   make(chan *Room),
		Exist:            make(chan struct{}),
	}
	go w.run()
	return w
}
func (ws *WebsocketMainInstance) findPrivateRoom(senderID string, reciverID string) (*Room, bool) {
	name := []string{senderID, reciverID}
	sort.Strings(name)
	combined := strings.Join(name, "_")
	for _, room := range ws.Rooms {
		if room.Name == combined && room.IsPrivate {
			return room, true
		}
	}
	return nil, false
}
func (ws *WebsocketMainInstance) disconnection() {
	ws.Exist <- struct{}{}
	for _, room := range ws.Rooms {
		room.Exist <- struct{}{}
	}
	for _, client := range ws.Clients {
		client.Exist <- struct{}{}
	}
}

func (ws *WebsocketMainInstance) run() {
	defer func() {
		close(ws.RegisterClient)
		close(ws.UnregisterClient)
		close(ws.RegisterRoom)
		close(ws.UnregisterRoom)
		close(ws.Exist)
		log.Println("close websocket instance")
	}()
	for {
		select {
		case client, ok := <-ws.RegisterClient:
			if !ok {
				log.Println("close websocket register client channel")
				return
			}
			ws.addClient(client)
		case client, ok := <-ws.UnregisterClient:
			if !ok {
				log.Println("close websocket unregister client channel")
				return
			}
			ws.removeClient(client)
		case room, ok := <-ws.RegisterRoom:
			if !ok {
				log.Println("close websocket register room channel")
				return
			}
			ws.addRoom(room)
		case room, ok := <-ws.UnregisterRoom:
			if !ok {
				log.Println("close websocket unregister room channel")
				return
			}
			ws.removeRoom(room)
		case <-ws.Exist:
			return
		}
	}
}
func (ws *WebsocketMainInstance) addClient(client *Client) {
	if _, exist := ws.Clients[client.ID]; exist {
		log.Println("duplicate client register")
	} else {
		log.Printf("add client %s from websocketinstance \n", client.ID)
		ws.Clients[client.ID] = client
	}
}
func (ws *WebsocketMainInstance) removeClient(client *Client) {
	if _, exist := ws.Clients[client.ID]; exist {
		delete(ws.Clients, client.ID)
		log.Printf("remove client %s from websocketinstance \n", client.ID)
	} else {
		log.Println("not exist client to delete")
	}
}
func (ws *WebsocketMainInstance) addRoom(room *Room) {
	if _, exist := ws.Rooms[room.ID]; exist {
		log.Println("duplicate room register")
	} else {
		log.Printf("add room %s from websocketinstance \n", room.Name)
		ws.Rooms[room.ID] = room
	}
}
func (ws *WebsocketMainInstance) removeRoom(room *Room) {
	if _, exist := ws.Rooms[room.ID]; exist {
		delete(ws.Rooms, room.ID)
		log.Printf("remove room %s from websocketinstance \n", room.ID)
	} else {
		log.Println("not exist room to delete")
	}
}
