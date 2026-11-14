package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/reowense/Link_shortener/internal/config"
	"github.com/reowense/Link_shortener/internal/handler"
	"github.com/reowense/Link_shortener/internal/middleware"
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
		Handler:           middleware.Logging(mux),
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       60 * time.Second,
	}
	slog.Info("Успешный старт",
		"PORT",
		cfg.Port)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	go func() {

		err := newServer.ListenAndServe()
		if !errors.Is(err, http.ErrServerClosed) {
			slog.Error("Сервер не поднят", "error", err)
			os.Exit(1)
		}

	}()
	<-ctx.Done()
	slog.Info("Получен сигнал остановки сервера")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err := newServer.Shutdown(shutdownCtx)
	if err != nil {
		slog.Error("Ошибка в остановке сервера", "error", err)
	}

	slog.Info("Сервер остановлен")
}
