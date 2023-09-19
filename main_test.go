package main

import (
	"counter/services"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCounterHandler(t *testing.T) {
	testFilePath := "requests.log"
	err := startup(testFilePath)
	if err != nil {
		t.Fatal(err)
	}

	req := httptest.NewRequest(http.MethodGet, "/counter", nil)
	w := httptest.NewRecorder()
	Counter(w, req)
	res := w.Result()
	defer res.Body.Close()
	data, err := io.ReadAll(res.Body)
	if err != nil {
		t.Errorf("expected error to be nil got %v", err)
	}
	if string(data) != `{"counter":0}` {
		t.Errorf(`expected {"counter":0} got %v`, string(data))
	}

	err = services.DeleteLogFile(testFilePath)
	if err != nil {
		t.Fatalf("Error removing the file: %v", err)
	}
}

func TestCounterHandlerMultiple(t *testing.T) {
	testFilePath := "requests.log"
	err := startup(testFilePath)
	if err != nil {
		t.Fatal(err)
	}

	for i := 0; i < 10; i++ {
		req := httptest.NewRequest(http.MethodGet, "/counter", nil)
		w := httptest.NewRecorder()
		Counter(w, req)
		res := w.Result()
		defer res.Body.Close()
		data, err := io.ReadAll(res.Body)
		if err != nil {
			t.Errorf("expected error to be nil got %v", err)
		}
		expected := fmt.Sprintf(`{"counter":%d}`, i)
		if string(data) != expected {
			t.Errorf(`expected %s got %v`, expected, string(data))
		}
	}

	err = services.DeleteLogFile(testFilePath)
	if err != nil {
		t.Fatalf("Error removing the file: %v", err)
	}
}
