package websocket

import (
	"log"
	"net/http"

	logger "github.com/XXena/chat/pkg/logger"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{} // use default options
var Done chan interface{}

func InitRoutes() *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/test", func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte("welcome"))
		if err != nil {
			return
		}
	})

	r.Get("/ws", socketHandler)
	return r
}

func socketHandler(w http.ResponseWriter, r *http.Request) {
	// Upgrade our raw HTTP connection to a websocket based one
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("Error during connection upgradation:", err)
		return
	}
	defer func(conn *websocket.Conn) {
		err := conn.Close()
		if err != nil {

		}
	}(conn)

	// эхо сервер
	for {
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("Error during message reading:", err)
			break
		}
		log.Printf("Received: %s", message)
		err = conn.WriteMessage(messageType, message)
		if err != nil {
			log.Println("Error during message writing:", err)
			break
		}
		log.Printf("Sent: %s", message)
	}
}

func ReceiveHandler(connection *websocket.Conn, log *logger.Logger) {
	defer close(Done)
	for {
		_, msg, err := connection.ReadMessage()
		if err != nil {
			log.Error("Error in receive:", err)
			return
		}
		log.Debug("Received: %s\n", msg)
	}
}

func SendHandler(messsageType int, connection *websocket.Conn, log *logger.Logger, message []byte) {
	defer close(Done)
	for {
		err := connection.WriteMessage(messsageType, message)
		if err != nil {
			log.Error("Error in send:", err)
			return
		}
		log.Debug("Sent: %s\n", message)
	}
}
