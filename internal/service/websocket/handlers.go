package websocket

import (
	"log"
	"net/http"

	"github.com/XXena/chat/pkg/logger"
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
	r.Get("/ws/{chat}", socketHandler)
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
			log.Print("Error closing websocket connection:", err)
		}
	}(conn)

	chatID := chi.URLParam(r, "chat")
	// эхо сервер
	for {
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("Error during message reading:", err)
			break
		}
		log.Printf("ID: %s Received: %s", chatID, message)
		err = conn.WriteMessage(messageType, message)
		if err != nil {
			log.Println("Error during message writing:", err)
			break
		}
		log.Printf("ID: %s Sent: %s", chatID, message)
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
