package main

import (
	"crypto/tls"
	"log/slog"
	"net/http"
	"os"
	"time"
)

const dbPath = "./database/data.db"

func main() {
	env := os.Getenv("ENV")
	if env == "" {
		env = "production"
	}

	allowedOrigin := "http://localhost:3000"
	if env == "production" {
		allowedOrigin = "https://dyyni.org"
	}
	slog.Info("env", "mode", env, "allowedOrigin", allowedOrigin)

	db, err := InitDB(dbPath)
	if err != nil {
		slog.Error("failed to init database", "error", err)
		os.Exit(1)
	}
	defer db.Close()

	limiter := NewLimiter(5, time.Minute)
	allowlist := LoadAllowlist()

	go StartPruneLoop(db)
	go StartLimiterCleanup(limiter)

	production := env == "production"

	mux := http.NewServeMux()
	mux.HandleFunc("/blackbox", BlackboxHandler(db, limiter, allowedOrigin, allowlist, production))

	if env == "production" {
		cert, err := generateSelfSignedCert()
		if err != nil {
			slog.Error("failed to generate cert", "error", err)
			os.Exit(1)
		}

		healthMux := http.NewServeMux()
		healthMux.HandleFunc("/health", HealthHandler)

		go func() {
			slog.Info("health server listening", "addr", "127.0.0.1:8081")
			if err := http.ListenAndServe("127.0.0.1:8081", healthMux); err != nil {
				slog.Error("health server failed", "error", err)
			}
		}()

		server := &http.Server{
			Addr:         ":443",
			Handler:      mux,
			ReadTimeout:  5 * time.Second,
			WriteTimeout: 10 * time.Second,
			IdleTimeout:  60 * time.Second,
			TLSConfig: &tls.Config{
				Certificates: []tls.Certificate{cert},
				MinVersion:   tls.VersionTLS12,
			},
		}

		slog.Info("listening", "addr", ":443", "tls", true)
		if err := server.ListenAndServeTLS("", ""); err != nil {
			slog.Error("server failed", "error", err)
			os.Exit(1)
		}
	} else {
		mux.HandleFunc("/health", HealthHandler)
		slog.Info("listening", "addr", ":8080", "tls", false)
		if err := http.ListenAndServe(":8080", mux); err != nil {
			slog.Error("server failed", "error", err)
			os.Exit(1)
		}
	}
}
