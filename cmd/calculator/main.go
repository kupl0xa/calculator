package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/kupl0xa/calculator/internal/router"
	"github.com/kupl0xa/calculator/internal/tracer"
)

func main() {
	tp, err := tracer.InitTracer()
	if err != nil {
		slog.Error("Failed to iintialize tracer", "error", err.Error())
		os.Exit(1)
	}
	defer tracer.ShutdownTracer(tp)

	r := router.NewRouter()

	srv := &http.Server{
		Addr:         "0.0.0.0:8080",
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      r,
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	go func() {
		slog.Info("Starting server on port 8080")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("Can't start server", "error", err.Error())
			os.Exit(1)
		}

	}()

	<-stop
	ctx, cancel := context.WithTimeout(context.Background(), 5)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		slog.Error(err.Error())
	}
	slog.Info("Shutting down")
	os.Exit(0)
}
