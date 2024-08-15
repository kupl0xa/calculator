package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/gorilla/mux"
	"github.com/kupl0xa/calculator/internal/models"
	"github.com/kupl0xa/calculator/internal/services"
)

func TestCreateTaskHandler(t *testing.T) {
	service := services.NewTaskService()
	handler := NewTaskHandler(service)

	body := []byte(`{"x": 10, "y": 5, "operator": "+"}`)
	req, err := http.NewRequest("POST", "/tasks", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler.CreateTask(rr, req)

	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("wrong status code: expected %v, got %v", http.StatusCreated, status)
	}

	var response map[string]int
	if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
		t.Fatal(err)
	}
	if response["id"] != 1 {
		t.Errorf("expected id 1, got %d", response["id"])
	}

	body = []byte(`{"x": 10, "y": 0, "operator": "/"}`)
	req, err = http.NewRequest("POST", "/tasks", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}
	rr = httptest.NewRecorder()
	handler.CreateTask(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("wrong status code: expected %v, got %v", http.StatusBadRequest, status)
	}

	body = []byte(`{"x": "a", "y": 0, "operator": "*"}`)
	req, err = http.NewRequest("POST", "/tasks", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}
	rr = httptest.NewRecorder()
	handler.CreateTask(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("wrong status code: expected %v, got %v", http.StatusBadRequest, status)
	}

	body = []byte(`{"x": 5, "y": 0, "operator": "f"}`)
	req, err = http.NewRequest("POST", "/tasks", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}
	rr = httptest.NewRecorder()
	handler.CreateTask(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("wrong status code: expected %v, got %v", http.StatusBadRequest, status)
	}
}

func TestGetTaskResultHandler(t *testing.T) {
	service := services.NewTaskService()
	handler := NewTaskHandler(service)
	id := service.CreateTask(10, 5, "+")
	service.ExecuteTask(id)
	router := mux.NewRouter()
	router.HandleFunc("/tasks/{id}", handler.GetTaskResult).Methods("GET")

	req, err := http.NewRequest("GET", "/tasks/"+strconv.Itoa(id), nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("wrong status code: expected %v, got %v", http.StatusOK, status)
	}

	var response models.TaskResultResponse
	if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
		t.Fatal(err)
	}
	if response.ID != id || response.Result != 15 {
		t.Errorf("expected ID:1, Result:15, got %+v", response)
	}

	req, err = http.NewRequest("GET", "/tasks/d", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr = httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("wrong status code: expected %v, got %v", http.StatusBadRequest, status)
	}

	req, err = http.NewRequest("GET", "/tasks/123", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr = httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("wrong status code: expected %v, got %v", http.StatusNotFound, status)
	}
}

func TestListTaskHandler(t *testing.T) {
	service := services.NewTaskService()
	handler := NewTaskHandler(service)
	service.CreateTask(10, 5, "+")
	service.CreateTask(10, 5, "*")

	req, err := http.NewRequest("GET", "/tasks", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler.ListTasks(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("wrong status code: expected %v, got %v", http.StatusOK, status)
	}

	var response []models.TaskStatusResponse
	if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
		t.Fatal(err)
	}
	if len(response) != 2 {
		t.Errorf("expected 2 tasks, got %d", len(response))
	}
}
