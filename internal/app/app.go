package app

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/XXena/chat/internal/handlers"

	"github.com/XXena/chat/pkg/logger"

	"github.com/XXena/chat/internal/config"
	"github.com/XXena/chat/internal/transport/websocket"
)

func Run(cfg *config.Config) {
	l := logger.New(cfg.Log.Level)
	errChan := make(chan error, 1)

	// todo services
	// todo handlers
	router := handlers.InitRoutes()

	go func() {
		websocket.RunServer(cfg, l, router, errChan)
	}()

	// todo grpc server
	// todo graceful shutdown

	// Waiting signals
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		l.Info(fmt.Sprintf("app - Run - signal: %s", s))
	case err := <-errChan:
		l.Error(fmt.Errorf("app - Run - error notify: %w", err))
	}

}
