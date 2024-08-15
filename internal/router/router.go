package router

import (
	"github.com/gorilla/mux"
	"github.com/kupl0xa/calculator/internal/handlers"
	"github.com/kupl0xa/calculator/internal/middleware"
	"github.com/kupl0xa/calculator/internal/services"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gorilla/mux/otelmux"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func NewRouter() *mux.Router {
	r := mux.NewRouter()
	taskService := services.NewTaskService()
	taskHandler := handlers.NewTaskHandler(taskService)

	r.Use(otelmux.Middleware("calculator-app"))
	r.Use(middleware.LoggingMiddleware)
	r.Handle("/metrics", promhttp.Handler())
	r.HandleFunc("/tasks", taskHandler.CreateTask).Methods("POST")
	r.HandleFunc("/tasks", taskHandler.ListTasks).Methods("GET")
	r.HandleFunc("/tasks/{id}", taskHandler.GetTaskResult).Methods("GET")

	return r
}
