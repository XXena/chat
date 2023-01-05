package websocket

import (
	"log"
	"net/http"

	"github.com/XXena/chat/internal/config"

	"github.com/XXena/chat/internal/service"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/gorilla/websocket"
)

var (
	upgrader = websocket.Upgrader{} // use default options
	Done     chan interface{}       // todo
)

type Handler struct {
	Hub service.IHub
	Cfg *config.Config
}

func (h Handler) InitRoutes() *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/test", func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte("welcome"))
		if err != nil {
			return
		}
	})
	r.Get("/ws/{chat}", h.socketHandler)
	return r
}

func (h Handler) socketHandler(w http.ResponseWriter, r *http.Request) {
	// Upgrade our raw HTTP connection to a websocket based one
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("Error during connection upgradation:", err)
		return
	}

	chatID := chi.URLParam(r, "chat")
	chat := service.NewWebsocketChat(
		h.Hub,
		conn,
		h.Cfg.SendBufferSize, // todo проверить, что читается buffer size 256
		chatID)

	chat.Hub.Register(chat)
	go chat.WritePump()
	go chat.ReadPump()

}
