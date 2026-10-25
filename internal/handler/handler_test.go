package handler

import (
	"errors"
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
