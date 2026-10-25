package main

import (
	"github.com/reowense/Link_shortener/internal/handler"
	"github.com/reowense/Link_shortener/internal/storage"
	"log"
	"net/http"
)

func main() {
	store := storage.NewMemory()
	h := handler.New(store)

	mux := http.NewServeMux()
	h.Register(mux)

	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
