package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand/v2"
	"net/http"
	"strings"
)

var store = map[string]string{}

const length = 8

type URL struct {
	Url string `json:"url"`
}

func generatorCode() string {
	alphabet := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

	result := make([]rune, length)

	for i := 0; i < length; i++ {
		randomIndex := rand.IntN(len(alphabet))
		result[i] = alphabet[randomIndex]
	}
	return string(result)
}

func shortenHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	var req URL

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Некорректный JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	req.Url = strings.TrimSpace(req.Url)

	if req.Url == "" {
		http.Error(w, "URL пуст", http.StatusBadRequest)
		return
	}

	codeUrl := generatorCode()
	_, ok := store[codeUrl]

	for ok {
		codeUrl = generatorCode()
		_, ok = store[codeUrl]
	}

	store[codeUrl] = req.Url
	w.WriteHeader(http.StatusCreated)

	fmt.Fprint(w, codeUrl)
}

func redirectHandler(w http.ResponseWriter, r *http.Request) {
	code := r.PathValue("code")
	value, ok := store[code]
	if !ok {
		http.Error(w, "Не найдено", http.StatusNotFound)
		return
	}

	http.Redirect(w, r, value, http.StatusFound)

}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "ok")
}

func main() {

	http.HandleFunc("GET /{code}", redirectHandler)
	http.HandleFunc("GET /health", healthHandler)
	http.HandleFunc("POST /shorten", shortenHandler)

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
