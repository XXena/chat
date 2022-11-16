package websocket

import (
	"github.com/XXena/chat/internal/config"
	"github.com/XXena/chat/pkg/logger"
	"github.com/gorilla/websocket"
)

func NewClient(cfg *config.Config, log *logger.Logger) *websocket.Conn {
	socketUrl := "ws://" + cfg.WebSocket.Client.Host + cfg.WebSocket.Client.Port + "/socket"
	conn, _, err := websocket.DefaultDialer.Dial(socketUrl, nil)
	if err != nil {
		log.Fatal("Error connecting to Websocket Server:", err)
	}
	defer conn.Close()

	return conn
}
