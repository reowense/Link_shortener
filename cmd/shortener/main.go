package main

import (
 "github.com/reowense/Link_shortener/internal/config"
 "github.com/reowense/Link_shortener/internal/handler"
 "github.com/reowense/Link_shortener/internal/storage"
 "log"
 "net/http"
 "time"
)

func main() {
 cfg := config.NewConfig()
 store := storage.NewMemory()
 h := handler.New(store)

 mux := http.NewServeMux()
 h.Register(mux)

 newServer := http.Server{
  Addr:              ":" + cfg.Port,
  Handler:           mux,
  ReadHeaderTimeout: 5 * time.Second,
  ReadTimeout:       10 * time.Second,
  WriteTimeout:      10 * time.Second,
  IdleTimeout:       60 * time.Second,
 }

 err := newServer.ListenAndServe()

 if err != nil {
  log.Fatal("ListenAndServe: ", err)
 }
}