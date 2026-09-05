package ports

import (
	"context"
	"time"

	"trip-history-api/internal/core/domain"
)

type TripRepository interface {
	CreateTrip(ctx context.Context, trip domain.Trip) (string, error)
	UpdateTrip(ctx context.Context, tripID string, updateData map[string]interface{}) error
	GetTripsByUserID(ctx context.Context, userID string) ([]domain.Trip, error)
	GetTripsHistory(ctx context.Context, userID string, startTime time.Time, endTime time.Time) ([]domain.Trip, error)
	
	SaveLocationHistory(ctx context.Context, tripID string, loc domain.LocationUpdate) error
	GetLocationHistory(ctx context.Context, tripID string) ([]domain.LocationUpdate, error)
	
	SaveCheckIn(ctx context.Context, checkIn domain.CheckIn) error
	ReportIssue(ctx context.Context, issue domain.Issue) error
	ClearAllTrips(ctx context.Context) error
}

type JobPlanRepository interface {
	CreateJobPlan(ctx context.Context, plan domain.JobPlan) (string, error)
	GetJobPlans(ctx context.Context) ([]domain.JobPlan, error)
	GetJobPlansByDriver(ctx context.Context, driverID string) ([]domain.JobPlan, error)
	UpdateJobPlanStatus(ctx context.Context, planID string, status string) error
	ClearAllJobPlans(ctx context.Context) error
}

type LocationCache interface {
	UpdateLatestLocation(ctx context.Context, tripID string, loc domain.LocationUpdate) error
	ClearAll(ctx context.Context) error
}

type UserRepository interface {
	CreateUserInAuth(ctx context.Context, email, password, name string) (string, error)
	DeleteUserInAuth(ctx context.Context, uid string) error
	SaveUser(ctx context.Context, user domain.User) error
	GetUser(ctx context.Context, uid string) (domain.User, error)
	GetAllUsers(ctx context.Context) ([]domain.User, error)
	UpdateFCMToken(ctx context.Context, uid string, token string) error
	DeleteUserInFirestore(ctx context.Context, uid string) error
	ClearAllUsers(ctx context.Context) error
	GenerateEmailVerificationLink(ctx context.Context, email string) (string, error)
}

type LogRepository interface {
	SaveLog(ctx context.Context, log domain.SystemLog) error
	ClearAllLogs(ctx context.Context) error
}
