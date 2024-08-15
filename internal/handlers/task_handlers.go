package handlers

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/kupl0xa/calculator/internal/models"
	"github.com/kupl0xa/calculator/internal/services"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	taskCounter = prometheus.NewCounterVec(prometheus.CounterOpts{Name: "task_requests_total", Help: "Total number of task requests"}, []string{"endpoint"})
)

func init() {
	prometheus.MustRegister(taskCounter)
}

type TaskHandler struct {
	Service *services.TaskService
}

func NewTaskHandler(service *services.TaskService) *TaskHandler {
	return &TaskHandler{Service: service}
}

func (h *TaskHandler) CreateTask(w http.ResponseWriter, r *http.Request) {
	taskCounter.WithLabelValues("create").Inc()

	var input struct {
		X        int    `json:"x"`
		Y        int    `json:"y"`
		Operator string `json:"operator"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		slog.Error("CreateTask", "error", err.Error())
		http.Error(w, "Incorrect data", http.StatusBadRequest)
		return
	}
	switch input.Operator {
	case "/":
		if input.Y == 0 {
			slog.Error("CreateTask", "error", "Division by zero")
			http.Error(w, "Division by zero is not supported", http.StatusBadRequest)
			return
		}
	case "+", "-", "*":
	default:
		slog.Error("CreateTask", "Wrong operator", input.Operator)
		http.Error(w, "Operator should be +,-,*,/", http.StatusBadRequest)
		return
	}
	id := h.Service.CreateTask(input.X, input.Y, input.Operator)
	go h.Service.ExecuteTask(id)
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(map[string]int{"id": id}); err != nil {
		slog.Error("CreateTask", "error", err.Error())
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}
}

func (h *TaskHandler) ListTasks(w http.ResponseWriter, r *http.Request) {
	taskCounter.WithLabelValues("list").Inc()

	tasks := h.Service.ListTasks()
	if err := json.NewEncoder(w).Encode(tasks); err != nil {
		slog.Error("ListTask", "error", err.Error())
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}
}

func (h *TaskHandler) GetTaskResult(w http.ResponseWriter, r *http.Request) {
	taskCounter.WithLabelValues("get_result").Inc()

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		slog.Error(err.Error())
		http.Error(w, "Invalid task ID", http.StatusBadRequest)
		return
	}
	result, err := h.Service.GetTaskResult(id)
	if err != nil {
		slog.Error(err.Error(), "id", id)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	if err := json.NewEncoder(w).Encode(models.TaskResultResponse{ID: id, Result: result}); err != nil {
		slog.Error("GetTaskResult", "error", err.Error())
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}
}
