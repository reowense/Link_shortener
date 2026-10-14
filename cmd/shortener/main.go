package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math/rand/v2"
	"net/http"
	"net/url"
	"strings"
)

var (
	ErrInvalidURL        = errors.New("не валидный URL")
	ErrUnsupportedScheme = errors.New("неподдерживаемая схема, разрешены только http и https")
	ErrInvalidHost       = errors.New("хост пуст")
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

func validateInputURL(urlRaw string) error {
	parsedURL, err := url.ParseRequestURI(urlRaw)

	if err != nil {
		return ErrInvalidURL
	}

	if parsedURL.Scheme != "http" && parsedURL.Scheme != "https" {
		return ErrUnsupportedScheme
	}

	if parsedURL.Host == "" {
		return ErrInvalidHost
	}

	return nil
}

func shortenHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	var req URL

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "некорректный JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	req.Url = strings.TrimSpace(req.Url)

	if req.Url == "" {
		http.Error(w, "URL пуст", http.StatusBadRequest)
		return
	}

	err = validateInputURL(req.Url)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
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
		http.Error(w, "не найдено", http.StatusNotFound)
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
