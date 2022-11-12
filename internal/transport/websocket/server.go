package websocket

import (
	"net/http"
	"time"

	"github.com/XXena/chat/internal/config"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"

	"github.com/XXena/chat/pkg/logger"
)

func initRoutes() *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/test", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("welcome"))
	})

	return r
}

func RunServer(cfg *config.Config, log *logger.Logger, listenErr chan error) {
	log.Info("starting webSocket server on %s", cfg.WebSocket.Port)
	router := initRoutes()

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
