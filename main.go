package main

import (
	"chat_room/jwt"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

const (
	SUCCESS = 0
	FAIL    = 1
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins
	},
}

type client struct {
	sender   string
	receiver string
	connect  *websocket.Conn
}
type chatRoom struct {
	connections map[string]map[string]*websocket.Conn
	lock        *sync.Mutex
	register    chan client
	unregister  chan client
	message     chan Message
}

type Message struct {
	sender   string
	receiver string
	content  []byte
}

func (c *chatRoom) set(sender string, receiver string, conn *websocket.Conn) bool {
	c.lock.Lock()
	defer c.lock.Unlock()
	if _, ok := c.connections[sender]; ok {
		return false
	}
	c.connections[sender][receiver] = conn
	return true
}

func (c *chatRoom) delete(sender string, receiver string) bool {
	c.lock.Lock()
	defer c.lock.Unlock()
	if _, ok := c.connections[sender]; ok {
		delete(c.connections, sender)
		return true
	} else {
		return false
	}
}

func (c *chatRoom) addUser(client client) {
	c.register <- client
}
func (c *chatRoom) removedUser(client client) {
	c.unregister <- client
}

var chatRoomInstance *chatRoom

func NewChatRoomInstance(c context.Context) *chatRoom {
	ch := &chatRoom{
		connections: make(map[string]map[string]*websocket.Conn),
		lock:        &sync.Mutex{},
		register:    make(chan client),
		unregister:  make(chan client),
	}
	go ch.detectRegister(c)
	go ch.processMessage(c)
	return ch
}

func main() {
	r := gin.Default()
	// detect register and unregister
	ctx, cancel := context.WithCancel(context.Background())
	chatRoomInstance = NewChatRoomInstance(ctx)
	defer cancel()

	r.GET("/websocket", jwt.AuthenticateJWT(), websocketHandler)
	r.Run()
}
func (cs *chatRoom) processMessage(c context.Context) {
	for {
		select {
		case c, ok := <-cs.message:
			if !ok {
				log.Println("error on processMessage")
			}
			cs.lock.Lock()
			if sender, senderOk := cs.connections[c.sender]; senderOk {
				if receiver, receiverOk := sender[c.receiver]; receiverOk {
					receiver.WriteJSON(c.content)
				}
			}
			if receiver, receiverOk := cs.connections[c.receiver]; receiverOk {
				if sender, senderOK := receiver[c.sender]; senderOK {
					sender.WriteJSON(c.content)
				}
			}
			cs.lock.Unlock()
		case <-c.Done():
			return
		}
	}
}
func (cs *chatRoom) detectRegister(c context.Context) {
	for {
		select {
		case c, ok := <-cs.register:
			if !ok {
				return
			}
			cs.set(c.sender, c.receiver, c.connect)
		case c, ok := <-cs.unregister:
			if !ok {
				return
			}
			cs.delete(c.sender, c.receiver)
		case <-c.Done():
			return
		}
	}
}
func websocketHandler(ctx *gin.Context) {
	conn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		log.Println("WebSocket upgrade failed:", err)
		return
	}
	defer conn.Close()

	receiver := ctx.Request.URL.Query().Get("id")

	claim, isExist := ctx.Get("jwt")

	if !isExist {
		log.Println("WebSocket upgrade failed:", err)
		return
	}
	sender, ok := claim.(*jwt.UserClaims)
	senderClient := client{sender: sender.UserID, receiver: receiver, connect: conn}
	if ok {
		// Register the client connection with the assigned ID
		chatRoomInstance.addUser(senderClient)
	} else {
		log.Panicln("fail auth socket")
		return
	}
	defer func() { chatRoomInstance.removedUser(senderClient) }()

	for {
		// Read message from client
		_, msg, err := conn.ReadMessage()
		if err != nil {
			if _, ok := err.(*websocket.CloseError); ok {
				log.Println("close room socket")
			} else {
				log.Println("Failed to read message:", err)
			}
			break
		}
		var message Message
		err = json.Unmarshal(msg, &message)

		if err != nil {
			log.Println("illegal RoomSocketEvent Data type")
			break
		}
		chatRoomInstance.message <- message
	}
}
