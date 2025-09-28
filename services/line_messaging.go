package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"line-oa-backend/config"
	"line-oa-backend/models"
)

type LINEMessagingService struct {
	config *config.Config
}

type PushMessageRequest struct {
	To       string    `json:"to"`
	Messages []Message `json:"messages"`
}

type Message struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

type LINEMessagingErrorResponse struct {
	Message string `json:"message"`
	Details []struct {
		Message  string `json:"message"`
		Property string `json:"property"`
	} `json:"details"`
}

func NewLINEMessagingService(cfg *config.Config) *LINEMessagingService {
	return &LINEMessagingService{
		config: cfg,
	}
}

// SendBookingConfirmation sends a booking confirmation message to the user
func (s *LINEMessagingService) SendBookingConfirmation(lineUserID string, booking *models.Booking) error {
	message := s.formatBookingMessage(booking)
	return s.sendPushMessage(lineUserID, message)
}

// formatBookingMessage creates a formatted message for booking confirmation
func (s *LINEMessagingService) formatBookingMessage(booking *models.Booking) string {
	return fmt.Sprintf(
		"🎉 การจองของคุณได้รับการยืนยันแล้ว!\n\n"+
			"📋 รายละเอียดการจอง:\n"+
			"• บริการ: %s\n"+
			"• วันที่: %s\n"+
			"• เวลา: %s\n"+
			"• หมายเหตุ: %s\n\n"+
			"ขอบคุณที่ใช้บริการของเรา 🙏",
		booking.ServiceName,
		booking.BookingDate.Format("2 January 2006"),
		booking.BookingDate.Format("15:04"),
		getNotesOrDefault(booking.Notes),
	)
}

// sendPushMessage sends a push message to a LINE user
func (s *LINEMessagingService) sendPushMessage(lineUserID, messageText string) error {
	pushURL := "https://api.line.me/v2/bot/message/push"

	pushReq := PushMessageRequest{
		To: lineUserID,
		Messages: []Message{
			{
				Type: "text",
				Text: messageText,
			},
		},
	}

	jsonData, err := json.Marshal(pushReq)
	if err != nil {
		return fmt.Errorf("failed to marshal push message request: %w", err)
	}

	req, err := http.NewRequest("POST", pushURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+s.config.LINEChannelAccessToken)

	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		var errorResp LINEMessagingErrorResponse
		if err := json.Unmarshal(body, &errorResp); err != nil {
			return fmt.Errorf("failed to parse error response: status %d, body: %s", resp.StatusCode, string(body))
		}
		return fmt.Errorf("LINE messaging API error: %s", errorResp.Message)
	}

	return nil
}

// SendCustomMessage sends a custom message to a LINE user
func (s *LINEMessagingService) SendCustomMessage(lineUserID, messageText string) error {
	return s.sendPushMessage(lineUserID, messageText)
}

func getNotesOrDefault(notes string) string {
	if notes == "" {
		return "ไม่มี"
	}
	return notes
}
