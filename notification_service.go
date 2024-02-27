package modak_rate_limiter

import (
	"errors"
	"fmt"
)

// NotificationService sends notifications using pre-configured notifications from the NotificationFactory.
type NotificationService struct {
	notifications map[NotificationType]*Notification
}

// NewNotificationService creates a new NotificationService.
func NewNotificationService(factory *NotificationFactory) *NotificationService {
	notifications := make(map[NotificationType]*Notification)

	for notificationType := range factory.configs {
		notification, _ := factory.CreateNotification(notificationType, "")
		notifications[notificationType] = notification
	}

	return &NotificationService{
		notifications: notifications,
	}
}

// Send sends a notification of the specified type to the target recipient with the given message.
func (ns *NotificationService) Send(notificationType NotificationType, recipient, message string) error {
	notification, ok := ns.notifications[notificationType]
	if !ok {
		return errors.New("unsupported notification type")
	}

	// Set the recipient for the notification
	notification.Recipient = recipient

	// Check if the notification is allowed
	if !notification.Allow() {
		fmt.Printf("Cant send %v notification to %s: %s. Rate limit excedeed.\n", notification.Name, recipient, message)
		return fmt.Errorf("rate limit exceeded for %v notifications to recipient %s", notification.Name, recipient)
	}

	// Send the notification
	fmt.Printf("Sending %v notification to %s: %s\n", notification.Name, recipient, message)
	return nil
}
