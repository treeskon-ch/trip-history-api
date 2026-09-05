package services

import (
	"context"
	"fmt"
	"time"

	"trip-history-api/internal/core/domain"
	"trip-history-api/internal/core/ports"
)

type tripService struct {
	tripRepo      ports.TripRepository
	locationCache ports.LocationCache
	userRepo      ports.UserRepository
}

func NewTripService(tripRepo ports.TripRepository, locationCache ports.LocationCache, userRepo ports.UserRepository) ports.TripService {
	return &tripService{
		tripRepo:      tripRepo,
		locationCache: locationCache,
		userRepo:      userRepo,
	}
}

func (s *tripService) StartTrip(ctx context.Context, userID string, planID string) (string, error) {
	// Validate that the user exists
	_, err := s.userRepo.GetUser(ctx, userID)
	if err != nil {
		return "", fmt.Errorf("invalid user: %w", err)
	}

	newTrip := domain.Trip{
		UserID:    userID,
		PlanID:    planID,
		Status:    "ongoing",
		StartTime: time.Now(),
	}
	return s.tripRepo.CreateTrip(ctx, newTrip)
}

func (s *tripService) EndTrip(ctx context.Context, tripID string, distance float64, imageURL string) error {
	now := time.Now()
	updateData := map[string]interface{}{
		"status":   "completed",
		"endTime":  now,
		"distance": distance,
	}
	if imageURL != "" {
		updateData["imageUrl"] = imageURL
	}
	return s.tripRepo.UpdateTrip(ctx, tripID, updateData)
}

func (s *tripService) GetTrips(ctx context.Context, userID string) ([]domain.Trip, error) {
	return s.tripRepo.GetTripsByUserID(ctx, userID)
}

func (s *tripService) GetTripsHistory(ctx context.Context, userID string, startTime time.Time, endTime time.Time) ([]domain.Trip, error) {
	return s.tripRepo.GetTripsHistory(ctx, userID, startTime, endTime)
}

func (s *tripService) UpdateLocation(ctx context.Context, tripID string, loc domain.LocationUpdate) error {
	// 1. อัปเดตลง Redis สำหรับ Real-time Web
	err := s.locationCache.UpdateLatestLocation(ctx, tripID, loc)
	if err != nil {
		return err
	}

	// 2. บันทึกลง Firestore เป็นประวัติย้อนหลัง
	return s.tripRepo.SaveLocationHistory(ctx, tripID, loc)
}

func (s *tripService) GetTripRoute(ctx context.Context, tripID string) ([]domain.LocationUpdate, error) {
	return s.tripRepo.GetLocationHistory(ctx, tripID)
}

func (s *tripService) CheckIn(ctx context.Context, checkIn domain.CheckIn) error {
	return s.tripRepo.SaveCheckIn(ctx, checkIn)
}

func (s *tripService) ReportIssue(ctx context.Context, issue domain.Issue) error {
	issue.ReportedAt = time.Now()
	return s.tripRepo.ReportIssue(ctx, issue)
}
