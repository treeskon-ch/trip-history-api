package domain

import "time"

type Trip struct {
	ID        string     `json:"id" firestore:"-"`
	PlanID    string     `json:"planId,omitempty" firestore:"planId,omitempty"`
	UserID    string     `json:"userId" firestore:"userId"`
	Status    string     `json:"status" firestore:"status"`
	StartTime time.Time  `json:"startTime" firestore:"startTime"`
	EndTime   *time.Time `json:"endTime,omitempty" firestore:"endTime,omitempty"`
	Distance  float64    `json:"distance,omitempty" firestore:"distance,omitempty"`
	ImageURL  string     `json:"imageUrl,omitempty" firestore:"imageUrl,omitempty"`
}
