package adapters

import (
	"fmt"
	"time"
)

const telegramAPI = "https://api.telegram.org/bot%s/sendMessage"

type Notifier struct {
	ChatID string
	Text   string
	URL    string
}

func NewNotifier(chatID, token string) *Notifier {
	return &Notifier{
		ChatID: chatID,
		URL:    fmt.Sprintf(telegramAPI, token),
	}
}

func (n *Notifier) Send(text string) error {
	formattedTime := time.Now().Format("2006-01-02 15:04:05")
	fmt.Printf("sending message: %s, time: %s\n", text, formattedTime)
	//msg := Notifier{
	//	ChatID: n.ChatID,
	//	Text:   text,
	//}
	//
	//msgBytes, err := json.Marshal(msg)
	//if err != nil {
	//	return fmt.Errorf("error marshalling message: %w", err)
	//}
	//
	//resp, err := http.Post(n.URL, "application/json", bytes.NewBuffer(msgBytes))
	//if err != nil {
	//	return fmt.Errorf("error sending request: %w", err)
	//}
	//defer resp.Body.Close()
	//
	//if resp.StatusCode != http.StatusOK {
	//	return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	//}

	return nil
}
