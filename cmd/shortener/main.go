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
	"sync"
)

var (
	ErrInvalidURL        = errors.New("не валидный URL")
	ErrUnsupportedScheme = errors.New("неподдерживаемая схема, разрешены только http и https")
	ErrInvalidHost       = errors.New("хост пуст")
)

const (
	length   = 8
	maxBytes = 4 * 1024
)

type URL struct {
	Url string `json:"url"`
}

type Storage struct {
	mu    sync.RWMutex
	links map[string]string
}

func NewStorage() *Storage {
	return &Storage{links: make(map[string]string)}
}

func (s *Storage) Get(code string) (string, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	value, ok := s.links[code]
	if !ok {
		return "", false
	}

	return value, true
}

func (s *Storage) Save(url string) string {
	s.mu.Lock()
	defer s.mu.Unlock()
	codeUrl := generatorCode()

	_, ok := s.links[codeUrl]
	for ok {
		codeUrl = generatorCode()
		_, ok = s.links[codeUrl]
	}

	s.links[codeUrl] = url
	return codeUrl
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

func (s *Storage) shortenHandler(w http.ResponseWriter, r *http.Request) {
	var req URL

	r.Body = http.MaxBytesReader(w, r.Body, maxBytes)

	err := json.NewDecoder(r.Body).Decode(&req)

	if err != nil {
		var maxBytesErr *http.MaxBytesError
		if errors.As(err, &maxBytesErr) {
			http.Error(w, "размер файла превышает лимит", http.StatusRequestEntityTooLarge)
			return
		}
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

	codeUrl := s.Save(req.Url)
	w.WriteHeader(http.StatusCreated)

	fmt.Fprint(w, codeUrl)
}

func (s *Storage) redirectHandler(w http.ResponseWriter, r *http.Request) {
	code := r.PathValue("code")

	value, ok := s.Get(code)
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
	storage := NewStorage()
	http.HandleFunc("GET /{code}", storage.redirectHandler)
	http.HandleFunc("GET /health", healthHandler)
	http.HandleFunc("POST /shorten", storage.shortenHandler)

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
