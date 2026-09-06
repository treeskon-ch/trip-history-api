package services

import (
	"context"
	"fmt"
	"math"
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

	// Calculate actual distance from GPS points
	locations, err := s.tripRepo.GetLocationHistory(ctx, tripID)
	if err == nil && len(locations) > 0 {
		calculatedDistance := calculateTotalDistance(locations)
		if calculatedDistance > 0 {
			distance = calculatedDistance
		}
	}

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

func calculateTotalDistance(locations []domain.LocationUpdate) float64 {
	if len(locations) < 2 {
		return 0.0
	}

	var totalDistance float64
	for i := 0; i < len(locations)-1; i++ {
		totalDistance += haversine(locations[i].Lat, locations[i].Lng, locations[i+1].Lat, locations[i+1].Lng)
	}
	return totalDistance
}

func haversine(lat1, lon1, lat2, lon2 float64) float64 {
	const R = 6371 // Earth's radius in kilometers
	
	dLat := (lat2 - lat1) * (math.Pi / 180.0)
	dLon := (lon2 - lon1) * (math.Pi / 180.0)

	lat1 = lat1 * (math.Pi / 180.0)
	lat2 = lat2 * (math.Pi / 180.0)

	a := math.Sin(dLat/2)*math.Sin(dLat/2) + math.Sin(dLon/2)*math.Sin(dLon/2)*math.Cos(lat1)*math.Cos(lat2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return R * c
}
