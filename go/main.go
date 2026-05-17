package main

import (
	"log/slog"
	"net/http"
	"os"
	"time"
)

const dbPath = "./database/data.db"

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	env := os.Getenv("ENV")
	if env == "" {
		env = "production"
	}

	allowedOrigin := "http://localhost:3000"
	if env == "production" {
		allowedOrigin = "https://dyyni.org"
	}
	slog.Info("allowedOrigin", "origin", allowedOrigin)

	db, err := InitDB(dbPath)
	if err != nil {
		slog.Error("failed to init database", "error", err)
		os.Exit(1)
	}
	defer db.Close()

	limiter := NewLimiter(5, time.Minute)

	go StartPruneLoop(db)
	go StartLimiterCleanup(limiter)

	http.HandleFunc("/blackbox", BlackboxHandler(db, limiter, allowedOrigin))

	slog.Info("listening", "port", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		slog.Error("server failed", "error", err)
		os.Exit(1)
	}
}
