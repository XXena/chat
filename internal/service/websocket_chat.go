package service

import (
	"log"
	"time"

	"github.com/gorilla/websocket"
)

type WebsocketChat struct {
	Hub    IHub
	chatID string

	// The websocket connection
	wsConn *websocket.Conn

	// Buffered channel for outbound messages
	sendChan chan []byte
}

func NewWebsocketChat(hub IHub, conn *websocket.Conn, bufferSize int, chatID string) *WebsocketChat {
	return &WebsocketChat{
		Hub:      hub,
		wsConn:   conn,
		sendChan: make(chan []byte, bufferSize),
		chatID:   chatID,
	}
}

func (c WebsocketChat) GetID() string {
	return c.chatID
}

func (c WebsocketChat) GetSendChan() chan []byte {
	return c.sendChan
}

// ReadPump pumps messages from the websocket connection to the Hub
func (c WebsocketChat) ReadPump() {
	defer func() {
		c.Hub.Unregister(c)
		err := c.wsConn.Close()
		if err != nil {
			return
		}
	}()
	c.wsConn.SetReadLimit(maxMessageSize)
	err := c.wsConn.SetReadDeadline(time.Now().Add(pongWait))
	if err != nil {
		return
	}
	c.wsConn.SetPongHandler(func(string) error {
		err := c.wsConn.SetReadDeadline(time.Now().Add(pongWait))
		if err != nil {
			return err
		}
		return nil
	})
	for {
		_, msg, err := c.wsConn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
				log.Printf("error: %v", err)
			}
			break
		}
		m := Message{
			UserID: "", // todo
			ChatID: c.GetID(),
			Data:   msg,
		}

		c.Hub.Broadcast(m)
	}
}

// WritePump pumps messages from the Hub to the websocket connection
func (c WebsocketChat) WritePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		err := c.wsConn.Close()
		if err != nil {
			return
		}
	}()
	for {
		select {
		case message, ok := <-c.sendChan:
			if !ok {
				err := c.write(websocket.CloseMessage, []byte{})
				if err != nil {
					return
				}
				return
			}
			if err := c.write(websocket.TextMessage, message); err != nil {
				return
			}
		case <-ticker.C:
			if err := c.write(websocket.PingMessage, []byte{}); err != nil {
				return
			}
		}
	}
}

// write writes a message with the given message type and payload
func (c WebsocketChat) write(mt int, payload []byte) error {
	err := c.wsConn.SetWriteDeadline(time.Now().Add(writeWait))
	if err != nil {
		return err
	}
	return c.wsConn.WriteMessage(mt, payload)
}
