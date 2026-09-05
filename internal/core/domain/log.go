package domain

import "time"

type SystemLog struct {
	UserID    string    `json:"userId" firestore:"userId"`
	EventType string    `json:"eventType" firestore:"eventType"` // e.g. "GPS_LOST", "NETWORK_LOST"
	Message   string    `json:"message" firestore:"message"`
	Timestamp time.Time `json:"timestamp" firestore:"timestamp"`
}
