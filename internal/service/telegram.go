package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

const telegramAPI = "https://api.telegram.org/bot%s/sendMessage"

type TelegramMessage struct {
	ChatID string `json:"chat_id"`
	Text   string `json:"text"`
}

func SendTelegramMessage(token, chatID, text string) error {
	url := fmt.Sprintf(telegramAPI, token)

	msg := TelegramMessage{
		ChatID: chatID,
		Text:   text,
	}

	msgBytes, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("error marshalling message: %w", err)
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(msgBytes))
	if err != nil {
		return fmt.Errorf("error sending request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return nil
}
