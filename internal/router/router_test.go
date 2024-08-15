package router

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestRouter(t *testing.T) {
	// service := services.NewTaskService()
	// handler := handlers.NewTaskHandler(service)
	r := NewRouter()

	tests := []struct {
		method       string
		url          string
		expectedCode int
	}{
		{"POST", "/tasks", http.StatusCreated},
		{"GET", "/tasks/1", http.StatusOK},
		{"GET", "/tasks", http.StatusOK},
	}

	for _, tt := range tests {
		req, err := http.NewRequest(tt.method, tt.url, strings.NewReader(`{"x": 10, "y": 5, "operator": "+"}`))
		if err != nil {
			t.Fatal(err)
		}
		rr := httptest.NewRecorder()
		r.ServeHTTP(rr, req)

		if status := rr.Code; status != tt.expectedCode {
			t.Errorf("expected %v, got %v", tt.expectedCode, status)
		}
	}
}
