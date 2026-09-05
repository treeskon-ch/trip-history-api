package domain

import "time"

type JobLocation struct {
	Lat     float64 `json:"lat" firestore:"lat"`
	Lng     float64 `json:"lng" firestore:"lng"`
	Address string  `json:"address" firestore:"address"`
}

type JobPlan struct {
	PlanID        string        `json:"planId" firestore:"-"`
	DriverID      string        `json:"driverId" firestore:"driverId"`
	Description   string        `json:"description" firestore:"description"`
	StartTime     time.Time     `json:"startTime" firestore:"startTime"`
	PickupPoint   JobLocation   `json:"pickupPoint" firestore:"pickupPoint"`
	DropOffPoints []JobLocation `json:"dropOffPoints" firestore:"dropOffPoints"`
	Status        string        `json:"status" firestore:"status"` // e.g., "pending", "in_progress", "completed", "cancelled"
	CreatedAt     time.Time     `json:"createdAt" firestore:"createdAt"`
}
