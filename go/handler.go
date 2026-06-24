package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net"
	"net/http"
	"strings"
)

const (
	// NOTE: ensure this matches with src/blackbox.js value
	MaxMessageLength = 160
	MaxRequestBody   = 128
)

type payload struct {
	Message string `json:"message"`
}

func getClientIP(r *http.Request, allowlist *IPAllowlist, production bool) string {
	if production && allowlist != nil {
		remoteIP, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			remoteIP = r.RemoteAddr
		}
		if !allowlist.Contains(net.ParseIP(remoteIP)) {
			return ""
		}
		return r.Header.Get("CF-Connecting-IP")
	}

	// NOTE: only dev environment stuff below

	if ip := r.Header.Get("CF-Connecting-IP"); ip != "" {
		return ip
	}

	if ip := r.Header.Get("X-Real-IP"); ip != "" {
		return ip
	}

	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return ip
}

func HealthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func BlackboxHandler(tgs *TGSender, limiter *Limiter, allowedOrigin string, allowlist *IPAllowlist, production bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "DENY")
		w.Header().Set("Referrer-Policy", "no-referrer")
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

		ip := getClientIP(r, allowlist, production)
		if ip == "" {
			w.WriteHeader(http.StatusForbidden)
			return
		}

		if !limiter.Allow(ip) {
			slog.Warn("rate limited", "ip", ip)
			w.WriteHeader(http.StatusTooManyRequests)
			return
		}

		body, err := io.ReadAll(io.LimitReader(r.Body, MaxRequestBody))
		if err != nil {
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

		// NOTE: ip geolocation must be enabled for the domain in
		//		 the CF dashboard (network -> ip geolocation)

		country := r.Header.Get("CF-IPCountry")
		flag := CCToFlag(country)
		text := fmt.Sprintf("%s `%s`: %q", flag, ip, msg)

		if err := tgs.SendMessage(text); err != nil {
			slog.Error("telegram send message failed", "ip", ip, "error", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		slog.Info("message sent", "ip", ip, "msg", msg)
		w.WriteHeader(http.StatusOK)
	}
}
