package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// LINEMessagingServiceAdapter implements LINEMessagingService
type LINEMessagingServiceAdapter struct {
	channelAccessToken string
}

// NewLINEMessagingServiceAdapter creates a new LINE messaging service adapter
func NewLINEMessagingServiceAdapter(channelAccessToken string) *LINEMessagingServiceAdapter {
	return &LINEMessagingServiceAdapter{
		channelAccessToken: channelAccessToken,
	}
}

type textMessage struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

type pushMessageRequest struct {
	To       string      `json:"to"`
	Messages []textMessage `json:"messages"`
}

type broadcastMessageRequest struct {
	Messages []textMessage `json:"messages"`
}

// SendTextMessage sends a text message to a specific user
func (s *LINEMessagingServiceAdapter) SendTextMessage(userID, message string) error {
	url := "https://api.line.me/v2/bot/message/push"
	
	reqBody := pushMessageRequest{
		To: userID,
		Messages: []textMessage{
			{
				Type: "text",
				Text: message,
			},
		},
	}

	return s.sendRequest(url, reqBody)
}

// SendFlexMessage sends a flex message to a specific user
func (s *LINEMessagingServiceAdapter) SendFlexMessage(userID string, flexMessage interface{}) error {
	url := "https://api.line.me/v2/bot/message/push"
	
	reqBody := map[string]interface{}{
		"to": userID,
		"messages": []interface{}{flexMessage},
	}

	return s.sendRequest(url, reqBody)
}

// BroadcastMessage sends a message to all users
func (s *LINEMessagingServiceAdapter) BroadcastMessage(message string) error {
	url := "https://api.line.me/v2/bot/message/broadcast"
	
	reqBody := broadcastMessageRequest{
		Messages: []textMessage{
			{
				Type: "text",
				Text: message,
			},
		},
	}

	return s.sendRequest(url, reqBody)
}

// sendRequest sends HTTP request to LINE API
func (s *LINEMessagingServiceAdapter) sendRequest(url string, body interface{}) error {
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+s.channelAccessToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("LINE API error: %s", string(body))
	}

	return nil
}
