package websocket

import (
	"log"

	"github.com/gorilla/websocket"
)

var Done chan interface{}

func ReceiveHandler(connection *websocket.Conn) {
	defer close(Done)
	for {
		_, msg, err := connection.ReadMessage()
		if err != nil {
			log.Println("Error in receive:", err)
			return
		}
		log.Printf("Received: %s\n", msg)
	}
}
