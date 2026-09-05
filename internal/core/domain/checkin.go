package domain

import "time"

type CheckIn struct {
	TripID    string    `json:"tripId" firestore:"tripId"`
	Lat       float64   `json:"lat" firestore:"lat"`
	Lng       float64   `json:"lng" firestore:"lng"`
	Timestamp time.Time `json:"timestamp" firestore:"timestamp"`
	Message   string    `json:"message,omitempty" firestore:"message,omitempty"`
	ImageURL  string    `json:"imageUrl,omitempty" firestore:"imageUrl,omitempty"`
}
