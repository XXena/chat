package app

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/websocket"

	wsClient "github.com/XXena/chat/internal/client/websocket"
	"github.com/XXena/chat/internal/handlers"

	"github.com/XXena/chat/pkg/logger"

	"github.com/XXena/chat/internal/config"
	wsTransport "github.com/XXena/chat/internal/transport/websocket"
)

func Run(cfg *config.Config) {
	l := logger.New(cfg.Log.Level)
	errChan := make(chan error, 1)

	// todo services
	// todo handlers
	router := handlers.InitRoutes()

	go func() {
		wsTransport.RunServer(cfg, l, router, errChan)
	}()

	conn := wsTransport.NewClient(cfg, l)
	go wsClient.ReceiveHandler(conn)

	// todo grpc server
	// todo graceful shutdown

	// Waiting signals
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	// Имитирует поведение клиента
	for {
		select {
		case <-time.After(time.Duration(1) * time.Millisecond * 1000):
			// Send an echo packet every second
			err := conn.WriteMessage(websocket.TextMessage, []byte("Hello from GolangDocs!"))
			if err != nil {
				log.Println("Error during writing to websocket:", err)
				return
			}

		case <-interrupt:
			// We received a SIGINT (Ctrl + C). Terminate gracefully...
			log.Println("Received SIGINT interrupt signal. Closing all pending connections")

			// Close our websocket connection
			err := conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("Error during closing websocket:", err)
				return
			}

			select {
			case s := <-interrupt:
				l.Info(fmt.Sprintf("app - Run - signal: %s", s))
			case err := <-errChan:
				l.Error(fmt.Errorf("app - Run - error notify: %w", err))
			case <-wsClient.Done:
				log.Println("Receiver Channel Closed! Exiting....")
			case <-time.After(time.Duration(1) * time.Second):
				log.Println("Timeout in closing receiving channel. Exiting....")
			}
			return
		}

	}

	//select {
	//case s := <-interrupt:
	//	l.Info(fmt.Sprintf("app - Run - signal: %s", s))
	//case err := <-errChan:
	//	l.Error(fmt.Errorf("app - Run - error notify: %w", err))
	//}

}
