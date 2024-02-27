package modak_rate_limiter

import (
	"testing"
)

func TestNotificationService_Send(t *testing.T) {
	configs := map[NotificationType]RateLimitConfig{
		Status:    {Value: 2, Timing: "minute"},
		News:      {Value: 1, Timing: "day"},
		Marketing: {Value: 3, Timing: "hour"},
	}

	factory := NewNotificationFactory(configs)

	service := NewNotificationService(factory)

	err := service.Send(News, "user@example.com", "news 1")
	if err != nil {
		t.Errorf("Error sending news notification: %v", err)
	}

	err = service.Send(Status, "user@example.com", "status 1")
	if err != nil {
		t.Errorf("Error sending status notification: %v", err)
	}

	err = service.Send(Marketing, "user@example.com", "marketing 1")
	if err != nil {
		t.Errorf("Error sending marketing notification: %v", err)
	}
}

func TestNotificationService_Send_RateLimitExceeded(t *testing.T) {
	configs := map[NotificationType]RateLimitConfig{
		Status: {Value: 1, Timing: "minute"}, // Rate limit: 1 per minute
	}

	factory := NewNotificationFactory(configs)

	service := NewNotificationService(factory)

	err := service.Send(Status, "user@example.com", "status 1")
	if err != nil {
		t.Errorf("Error sending status notification: %v", err)
	}

	// Try sending another status notification within the same minute
	err = service.Send(Status, "user@example.com", "status 2")
	if err == nil {
		t.Error("Expected error for rate limit exceeded")
	}
}
