package domain

import "time"

type Issue struct {
	IssueID     string    `json:"issueId" firestore:"-"`
	TripID      string    `json:"tripId" firestore:"tripId"`
	Title       string    `json:"title" firestore:"title"`
	Description string    `json:"description" firestore:"description"`
	ReportedAt  time.Time `json:"reportedAt" firestore:"reportedAt"`
}
