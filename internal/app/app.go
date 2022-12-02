package app

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	ws "github.com/XXena/chat/internal/service/websocket"

	"github.com/XXena/chat/pkg/logger"

	"github.com/XXena/chat/internal/config"
	httpTransport "github.com/XXena/chat/internal/transport/http"
)

func Run(cfg *config.Config) {
	l := logger.New(cfg.Log.Level)
	errChan := make(chan error, 1)
	router := ws.InitRoutes()

	go func() {
		httpTransport.RunServer(cfg, l, router, errChan)
	}()

	// Waiting signals
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		l.Info(fmt.Sprintf("app - Run - signal: %s", s))
	case err := <-errChan:
		l.Error(fmt.Errorf("app - Run - error notify: %w", err))
		//case <-wsClient.Done:
		//	log.Println("Receiver Channel Closed! Exiting....")
	}
	return

}
