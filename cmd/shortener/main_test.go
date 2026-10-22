package main

import (
	"errors"
	"fmt"
	"sync"
	"testing"
)

func TestValidateInputURL(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  error
	}{
		{
			name:  "Нормальная ссылка",
			input: "https://go.dev/path?x=1",
			want:  nil,
		},
		{
			name:  "Пустая строка",
			input: "",
			want:  ErrInvalidURL,
		},
		{
			name:  "Нет схемы",
			input: "vk.com",
			want:  ErrInvalidURL,
		},
		{
			name:  "Неверная схема",
			input: "ftp://site.com",
			want:  ErrUnsupportedScheme,
		},
		{
			name:  "Неверный URL",
			input: "javascript:alert(1)",
			want:  ErrUnsupportedScheme,
		},
		{
			name:  "Нет хоста",
			input: "https://",
			want:  ErrInvalidHost,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := validateInputURL(tc.input)
			if !errors.Is(got, tc.want) {
				t.Errorf("Результат %q : %v    |    Ожидался %v", tc.input, got, tc.want)
			}
		})
	}
}

func TestStorageConcurrentSave(t *testing.T) {

	stor := NewStorage()
	const goCount = 100

	var wg sync.WaitGroup

	resCode := make([]string, goCount)
	resLink := make([]string, goCount)
	for i := 0; i < goCount; i++ {
		wg.Add(1)
		link := fmt.Sprintf("https://example.com/%d", i)
		go func() {
			defer wg.Done()
			resCode[i] = stor.Save(link)
			resLink[i] = link
		}()

	}
	wg.Wait()

	resCodeChek := make(map[string]struct{})

	for i := range resCode {
		if resCode[i] == "" {
			t.Errorf("Код пуст | Номер %d", i)
			continue
		}
		_, ok := resCodeChek[resCode[i]]
		if ok {
			t.Errorf("Найден повторный код  %d! %q", i, resCode[i])
		}

		value, ok := stor.Get(resCode[i])
		if value != resLink[i] {
			t.Errorf("Ожидали %q  | Получили %q", resLink[i], value)
		}
		if !ok {
			t.Errorf("Ссылка не найдена %q номер %d", resCode[i], i)
		}

		resCodeChek[resCode[i]] = struct{}{}
	}

	if goCount != len(resCodeChek) {
		t.Errorf("Не все прошли")
	}

}
