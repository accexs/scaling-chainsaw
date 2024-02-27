package modak_rate_limiter

import (
	"errors"
	"golang.org/x/time/rate"
)

// NotificationType represents the type of notification.
type NotificationType int

const (
	Status NotificationType = iota
	News
	Marketing
)

// notificationTypeNames a string representation of the notification type
var notificationTypeNames = map[NotificationType]string{
	Status:    "Status",
	News:      "News",
	Marketing: "Marketing",
}

// RateLimitConfig represents the rate limit configuration for a notification type.
type RateLimitConfig struct {
	Value  int
	Timing string // e.g., "minute", "hour", "day"
}

// NotificationFactory creates notifications with configurable rate limits.
type NotificationFactory struct {
	configs map[NotificationType]RateLimitConfig
}

// NewNotificationFactory creates a new NotificationFactory with given configurations.
func NewNotificationFactory(configs map[NotificationType]RateLimitConfig) *NotificationFactory {
	return &NotificationFactory{
		configs: configs,
	}
}

// CreateNotification creates a notification of the specified type for the given recipient.
func (nf *NotificationFactory) CreateNotification(notificationType NotificationType, recipient string) (*Notification, error) {
	config, ok := nf.configs[notificationType]
	if !ok {
		return nil, errors.New("unsupported notification type")
	}

	var limiter *rate.Limiter
	switch config.Timing {
	case "minute":
		limiter = rate.NewLimiter(rate.Limit(config.Value), config.Value)
	case "hour":
		limiter = rate.NewLimiter(rate.Limit(config.Value), config.Value*60)
	case "day":
		limiter = rate.NewLimiter(rate.Limit(config.Value), config.Value*1440)
	default:
		limiter = rate.NewLimiter(rate.Limit(config.Value), 1)
	}

	name, ok := notificationTypeNames[notificationType]
	if !ok {
		return nil, errors.New("unsupported notification type")
	}

	return &Notification{
		Type:      notificationType,
		Name:      name,
		Recipient: recipient,
		Limiter:   limiter,
	}, nil
}

// Notification represents a notification.
type Notification struct {
	Type      NotificationType
	Name      string
	Recipient string
	Limiter   *rate.Limiter
}

// Allow checks whether the notification is allowed to be sent.
func (n *Notification) Allow() bool {
	return n.Limiter.Allow()
}
