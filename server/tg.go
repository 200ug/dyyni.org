package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

const tgBaseURL = "https://api.telegram.org"

type TGSender struct {
	BotToken string `env:"TELEGRAM_BOT_TOKEN"`
	ChatID   string `env:"TELEGRAM_CHAT_ID"`
}

func NewTGSender() (*TGSender, error) {
	if err := godotenv.Load(); err != nil {
		return nil, err
	}
	var tgSender TGSender
	if err := env.Parse(&tgSender); err != nil {
		return nil, err
	}
	if !isTokenValid(tgSender.BotToken) {
		return nil, fmt.Errorf("invalid telegram bot token")
	}
	if tgSender.ChatID == "" {
		slog.Info("missing telegram chat id, visit https://api.telegram.org/bot{BOT_ID}/getUpdates to get it")
		return nil, fmt.Errorf("missing chat id, send a message to the bot and visit ")
	}
	return &tgSender, nil
}

func isTokenValid(botToken string) bool {
	getMeURL := fmt.Sprintf("%s/bot%s/getMe", tgBaseURL, botToken)
	resp, err := http.Get(getMeURL)
	if err != nil {
		return false
	}
	if resp.StatusCode != 200 {
		return false
	}
	return true
}

func (tgs *TGSender) SendMessage(text string) error {
	url := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", tgs.BotToken)
	body, _ := json.Marshal(map[string]string{
		"chat_id":    tgs.ChatID,
		"text":       text,
		"parse_mode": "MarkdownV2",
	})
	resp, err := http.Post(url, "application/json", bytes.NewReader(body))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("telegram api returned %d", resp.StatusCode)
	}
	return nil
}
