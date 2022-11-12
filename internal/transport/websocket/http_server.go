package websocket

import (
	"net/http"
	"time"

	"github.com/go-chi/chi"

	"github.com/XXena/chat/internal/config"
	"github.com/XXena/chat/pkg/logger"
)

func RunServer(cfg *config.Config, log *logger.Logger, router *chi.Mux, listenErr chan error) {
	log.Info("starting webSocket server on %s", cfg.WebSocket.Port)

	rt, err := time.ParseDuration(cfg.ReadTimeout)
	if err != nil {
		log.Error("failed to run ws server: unable to parse read timeout", err)
		listenErr <- err
		return
	}

	wt, err := time.ParseDuration(cfg.WriteTimeout)
	if err != nil {
		log.Error("failed to run ws server: unable to parse write timeout", err)
		listenErr <- err
		return
	}

	srv := http.Server{
		Addr:         cfg.WebSocket.Host + cfg.WebSocket.Port,
		Handler:      router,
		ReadTimeout:  rt,
		WriteTimeout: wt,
	}

	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		listenErr <- err
	}
}
