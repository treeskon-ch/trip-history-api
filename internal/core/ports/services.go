package ports

import (
	"context"
	"time"

	"trip-history-api/internal/core/domain"
)

type NotificationService interface {
	SendJobPlanNotification(ctx context.Context, fcmToken string, title string, body string) error
}

type TripService interface {
	StartTrip(ctx context.Context, userID string, planID string) (string, error)
	EndTrip(ctx context.Context, tripID string, distance float64, imageURL string) error
	GetTrips(ctx context.Context, userID string) ([]domain.Trip, error)
	GetTripsHistory(ctx context.Context, userID string, startTime time.Time, endTime time.Time) ([]domain.Trip, error)
	
	UpdateLocation(ctx context.Context, tripID string, loc domain.LocationUpdate) error
	GetTripRoute(ctx context.Context, tripID string) ([]domain.LocationUpdate, error)
	
	CheckIn(ctx context.Context, checkIn domain.CheckIn) error
	ReportIssue(ctx context.Context, issue domain.Issue) error
}

type JobPlanService interface {
	CreateJobPlan(ctx context.Context, plan domain.JobPlan) (string, error)
	GetAllJobPlans(ctx context.Context) ([]domain.JobPlan, error)
	GetJobPlansByDriver(ctx context.Context, driverID string) ([]domain.JobPlan, error)
	UpdateJobPlanStatus(ctx context.Context, planID string, status string) error
}

type UserService interface {
	RegisterUser(ctx context.Context, email, password, name, role string) (domain.User, string, error)
	DeleteUser(ctx context.Context, uid string) error
	GetAllUsers(ctx context.Context) ([]domain.User, error)
	UpdateFCMToken(ctx context.Context, uid string, token string) error
}

type LogService interface {
	SaveLog(ctx context.Context, log domain.SystemLog) error
}

type SystemService interface {
	ClearAllData(ctx context.Context) error
}
