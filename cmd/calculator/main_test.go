package main

import (
	"net/http"
	"testing"
	"time"
)

func TestMain(t *testing.T) {
	go main()
	time.Sleep(1 * time.Second)

	resp, err := http.Get("http://localhost:8080/tasks")
	if err != nil {
		t.Fatalf("Failed to send GET request: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected status code 200, got %v", resp.StatusCode)
	}
}
