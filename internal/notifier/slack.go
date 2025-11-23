package notifier

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/antigravity/dbbackup/internal/config"
)

type Notifier interface {
	Notify(message string) error
}

type SlackNotifier struct {
	WebhookURL string
}

func NewSlackNotifier(cfg config.NotifyConfig) *SlackNotifier {
	return &SlackNotifier{WebhookURL: cfg.SlackWebhookURL}
}

func (s *SlackNotifier) Notify(message string) error {
	if s.WebhookURL == "" {
		return nil // No-op if no webhook URL
	}

	payload := map[string]string{"text": message}
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	resp, err := http.Post(s.WebhookURL, "application/json", bytes.NewBuffer(jsonPayload))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("slack notification failed with status: %d", resp.StatusCode)
	}

	return nil
}
