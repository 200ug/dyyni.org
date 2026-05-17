package main

import (
	"database/sql"
	"encoding/json"
	"io"
	"log/slog"
	"net"
	"net/http"
	"strings"
)

const (
	// NOTE: ensure this matches with src/blackbox.js value
	MaxMessageLength = 80
	MaxRequestBody   = 128
)

type payload struct {
	Message string `json:"message"`
}

func getClientIP(r *http.Request) string {
	if ip := r.Header.Get("X-Real-IP"); ip != "" {
		return ip
	}
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return ip
}

func BlackboxHandler(db *sql.DB, limiter *Limiter, allowedOrigin string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", allowedOrigin)
		w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		ip := getClientIP(r)

		if !limiter.Allow(ip) {
			slog.Warn("rate limited", "ip", ip)
			w.WriteHeader(http.StatusTooManyRequests)
			return
		}

		body, err := io.ReadAll(io.LimitReader(r.Body, MaxRequestBody))
		if err != nil {
			slog.Error("read body failed", "ip", ip, "error", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		var p payload
		if err := json.Unmarshal(body, &p); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		msg := strings.TrimSpace(p.Message)
		if msg == "" || len(msg) > MaxMessageLength {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if err := InsertMessage(db, msg); err != nil {
			slog.Error("insert failed", "ip", ip, "error", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		slog.Info("message stored", "ip", ip, "msg", msg)
		w.WriteHeader(http.StatusOK)
	}
}
