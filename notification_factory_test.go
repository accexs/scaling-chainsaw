package modak_rate_limiter

import (
	"golang.org/x/time/rate"
	"testing"
)

func TestNotificationFactory_CreateNotification(t *testing.T) {
	configs := map[NotificationType]RateLimitConfig{
		Status:    {Value: 2, Timing: "minute"},
		News:      {Value: 1, Timing: "day"},
		Marketing: {Value: 3, Timing: "hour"},
	}

	factory := NewNotificationFactory(configs)

	recipient := "test@example.com"

	statusNotification, err := factory.CreateNotification(Status, recipient)
	if err != nil {
		t.Errorf("Error creating status notification: %v", err)
	}
	if statusNotification.Type != Status {
		t.Errorf("Expected notification type to be Status, got %v", statusNotification.Type)
	}

	newsNotification, err := factory.CreateNotification(News, recipient)
	if err != nil {
		t.Errorf("Error creating news notification: %v", err)
	}
	if newsNotification.Type != News {
		t.Errorf("Expected notification type to be News, got %v", newsNotification.Type)
	}

	marketingNotification, err := factory.CreateNotification(Marketing, recipient)
	if err != nil {
		t.Errorf("Error creating marketing notification: %v", err)
	}
	if marketingNotification.Type != Marketing {
		t.Errorf("Expected notification type to be Marketing, got %v", marketingNotification.Type)
	}
}

func TestNotification_Allow(t *testing.T) {
	limiter := rate.NewLimiter(rate.Limit(2), 1)
	notification := &Notification{
		Type:      Status,
		Recipient: "test@example.com",
		Limiter:   limiter,
	}

	if !notification.Allow() {
		t.Error("Expected notification to be allowed")
	}
}

func TestNotificationFactory_CreateNotification_InvalidType(t *testing.T) {
	configs := map[NotificationType]RateLimitConfig{
		Status: {Value: 2, Timing: "minute"},
		News:   {Value: 1, Timing: "day"},
	}

	factory := NewNotificationFactory(configs)

	recipient := "test@example.com"

	// Try to create a notification with an unsupported type
	_, err := factory.CreateNotification(Marketing, recipient)
	if err == nil {
		t.Error("Expected error for creating notification with unsupported type")
	}
}

func TestNotification_Allow_RateLimitExceeded(t *testing.T) {
	limiter := rate.NewLimiter(rate.Limit(1), 1) // Rate limit: 1 per minute
	notification := &Notification{
		Type:      News,
		Recipient: "test@example.com",
		Limiter:   limiter,
	}

	// Allow the first notification
	if !notification.Allow() {
		t.Error("Expected notification to be allowed")
	}

	// Try to allow another notification within the same minute
	if notification.Allow() {
		t.Error("Expected notification to be rate limited")
	}
}
