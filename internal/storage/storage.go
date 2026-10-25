package storage

import (
	"math/rand/v2"
	"sync"
)

type MemoryStorage struct {
	mu    sync.RWMutex
	links map[string]string
}

func NewMemory() *MemoryStorage {
	return &MemoryStorage{links: make(map[string]string)}
}

func (s *MemoryStorage) Get(code string) (string, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	value, ok := s.links[code]

	return value, ok
}

func (s *MemoryStorage) Save(link string) string {
	s.mu.Lock()
	defer s.mu.Unlock()
	codeUrl := generatorCode()

	_, ok := s.links[codeUrl]
	for ok {
		codeUrl = generatorCode()
		_, ok = s.links[codeUrl]
	}

	s.links[codeUrl] = link
	return codeUrl
}

const length = 8

func generatorCode() string {
	alphabet := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

	result := make([]rune, length)

	for i := 0; i < length; i++ {
		randomIndex := rand.IntN(len(alphabet))
		result[i] = alphabet[randomIndex]
	}
	return string(result)
}
