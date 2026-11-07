package main

import (
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/reowense/Link_shortener/internal/config"
	"github.com/reowense/Link_shortener/internal/handler"
	"github.com/reowense/Link_shortener/internal/storage"
)

func main() {
	jsonLog := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(jsonLog)

	cfg := config.NewConfig()
	store := storage.NewMemory()
	h := handler.New(store)

	mux := http.NewServeMux()
	h.Register(mux)

	newServer := http.Server{
		Addr:              cfg.Port,
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       60 * time.Second,
	}
	slog.Info("Успешный старт",
		"PORT",
		cfg.Port)

	err := newServer.ListenAndServe()

	if err != nil {
		slog.Error("Сервер не поднят", "error", err)
		os.Exit(1)
	}
}
