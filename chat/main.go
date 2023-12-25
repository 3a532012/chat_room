package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"sort"
	"strings"
	"syscall"
	"time"
	"utils/jwt"

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

var instance *WebsocketMainInstance

func main() {
	r := gin.Default()
	instance = NewWebsocketMainInstance()
	r.GET("/auth/:id", generateJWT) //just generate jwt toke for test
	r.GET("/websocket/:id", jwt.AuthenticateJWT(), websocketHandler)
	srv := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	go func() {
		// service connections
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutdown Server ...")
	instance.disconnection()
	time.Sleep(1 * time.Second)
	log.Println("Shutdown Finished")
}

func generateJWT(ctx *gin.Context) {
	id := ctx.Param("id")
	newJwt, err := jwt.CreateJWT(id)
	if err != nil {
		Success(ctx, FAIL, err.Error(), nil)
	}
	Success(ctx, SUCCESS, "Get jwt", newJwt)

}
func websocketHandler(ctx *gin.Context) {
	conn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		log.Println("WebSocket upgrade failed:", err)
		return
	}
	defer conn.Close()

	receiver := ctx.Param("id")
	log.Printf("reciver %s", receiver)
	claim, isExist := ctx.Get("jwt")

	if !isExist {
		log.Println("WebSocket upgrade failed:", err)
		return
	}
	sender, ok := claim.(*jwt.UserClaims)
	if !ok {
		log.Println("jwt format fail")
		return
	}
	room, exist := instance.findPrivateRoom(sender.UserID, receiver)
	if !exist {
		//private room name gernate rule
		name := []string{sender.UserID, receiver}
		sort.Strings(name)
		combined := strings.Join(name, "_")
		room = NewRoom(combined, true)
		instance.addRoom(room)
	}
	client := NewClient(sender.UserID, sender.Username, conn, room)

	instance.addClient(client)
	defer instance.removeClient(client)

	room.addClient(client)
	defer room.removeClient(client)
	client.readPump()
}
func Error(ctx *gin.Context, statusCode int, code int, message *string) {
	ctx.JSON(statusCode, gin.H{
		"errorCode":    code,
		"errorMessage": message,
	})
}

func Success(ctx *gin.Context, code int, message string, data interface{}) {
	ctx.JSON(http.StatusOK, gin.H{
		"code":    code,
		"message": message,
		"data":    data,
	})
}
