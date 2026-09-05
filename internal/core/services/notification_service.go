package services

import (
	"context"

	"firebase.google.com/go/v4/messaging"
	"trip-history-api/internal/core/ports"
)

type notificationService struct {
	client *messaging.Client
}

func NewNotificationService(client *messaging.Client) ports.NotificationService {
	return &notificationService{
		client: client,
	}
}

func (s *notificationService) SendJobPlanNotification(ctx context.Context, fcmToken string, title string, body string) error {
	if fcmToken == "" {
		return nil // Do nothing if the user has no token
	}

	msg := &messaging.Message{
		Notification: &messaging.Notification{
			Title: title,
			Body:  body,
		},
		Token: fcmToken,
	}

	_, err := s.client.Send(ctx, msg)
	return err
}
