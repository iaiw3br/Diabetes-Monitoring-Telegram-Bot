package adapters

import (
	"fmt"
	"os"
)

const telegramAPI = "https://api.telegram.org/bot%s/sendMessage"

type Notifier struct {
	ChatID string
	Text   string
	URL    string
}

type TelegramMessage struct {
	ChatID string `json:"chat_id"`
	Text   string `json:"text"`
}

func NewNotifier() *Notifier {
	token := os.Getenv("TELEGRAM_BOT_TOKEN")
	chatID := os.Getenv("TELEGRAM_CHAT_ID")
	return &Notifier{
		ChatID: chatID,
		URL:    fmt.Sprintf(telegramAPI, token),
	}
}

func (n *Notifier) Send(text string) error {
	fmt.Println("Sending message", text)
	//msg := TelegramMessage{
	//	ChatID: n.ChatID,
	//	Text:   text,
	//}
	//
	//body, err := json.Marshal(msg)
	//if err != nil {
	//	return err
	//}
	//
	//resp, err := http.Post(n.URL, "application/json", bytes.NewBuffer(body))
	//if err != nil {
	//	return err
	//}
	//defer resp.Body.Close()
	//
	//if resp.StatusCode != http.StatusOK {
	//	return fmt.Errorf("unexpected status: %s", resp.Status)
	//}

	return nil
}
