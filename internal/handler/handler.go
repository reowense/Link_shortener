package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

type URL struct {
	Url string `json:"url"`
}

const (
	maxBytes = 4 * 1024
)

var (
	ErrInvalidURL        = errors.New("не валидный URL")
	ErrUnsupportedScheme = errors.New("неподдерживаемая схема, разрешены только http и https")
	ErrInvalidHost       = errors.New("хост пуст")
)

type StorURL interface {
	Get(code string) (string, bool)
	Save(link string) string
}

type Handler struct {
	store StorURL
}

func New(s StorURL) *Handler {
	return &Handler{
		store: s,
	}
}

func (s Handler) Register(mux *http.ServeMux) {
	mux.HandleFunc("GET /{code}", s.RedirectHandler)
	mux.HandleFunc("GET /health", HealthHandler)
	mux.HandleFunc("POST /shorten", s.ShortenHandler)
}

func (s Handler) ShortenHandler(w http.ResponseWriter, r *http.Request) {
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

	codeUrl := s.store.Save(req.Url)
	w.WriteHeader(http.StatusCreated)

	fmt.Fprint(w, codeUrl)
}

func (s Handler) RedirectHandler(w http.ResponseWriter, r *http.Request) {
	code := r.PathValue("code")

	value, ok := s.store.Get(code)
	if !ok {
		http.Error(w, "не найдено", http.StatusNotFound)
		return
	}

	http.Redirect(w, r, value, http.StatusFound)
}

func HealthHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "ok")
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
