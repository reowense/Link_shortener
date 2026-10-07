package main

import (
	"fmt"
	"log"
	"net/http"
)

func healthHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "ok")
}

func main() {
	http.HandleFunc("/health", healthHandler)

	err:=http.ListenAndServe(":8080",nil)
	if err!=nil{
		log.Fatal("ListenAndServe: ",err)
	}
}