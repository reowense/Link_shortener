package storage

import (
	"fmt"
	"sync"
	"testing"
)

func TestStorageConcurrentSave(t *testing.T) {

	stor := NewMemory()
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
