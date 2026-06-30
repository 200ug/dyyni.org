package main

import (
	"crypto/tls"
	"io"
	"log"
	"log/slog"
	"net/http"
	"os"
	"time"
)

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

	tg, err := NewTGSender()
	if err != nil {
		slog.Error("telegram sender creation failed", "error", err)
		os.Exit(1)
	}

	limiter := NewLimiter(5, time.Minute)
	allowlist := LoadAllowlist()

	go StartLimiterCleanup(limiter)

	production := env == "production"

	mux := http.NewServeMux()
	mux.HandleFunc("/blackbox", BlackboxHandler(tg, limiter, allowedOrigin, allowlist, production))

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
			ErrorLog:     log.New(io.Discard, "", 0), // prevents log spam from handshake errors etc.
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
