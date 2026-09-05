package services

import (
	"context"

	"trip-history-api/internal/core/ports"
)

type systemService struct {
	userRepo      ports.UserRepository
	tripRepo      ports.TripRepository
	logRepo       ports.LogRepository
	locationCache ports.LocationCache
	jobPlanRepo   ports.JobPlanRepository
}

func NewSystemService(userRepo ports.UserRepository, tripRepo ports.TripRepository, logRepo ports.LogRepository, locationCache ports.LocationCache, jobPlanRepo ports.JobPlanRepository) ports.SystemService {
	return &systemService{
		userRepo:      userRepo,
		tripRepo:      tripRepo,
		logRepo:       logRepo,
		locationCache: locationCache,
		jobPlanRepo:   jobPlanRepo,
	}
}

func (s *systemService) ClearAllData(ctx context.Context) error {
	if err := s.userRepo.ClearAllUsers(ctx); err != nil {
		return err
	}
	if err := s.tripRepo.ClearAllTrips(ctx); err != nil {
		return err
	}
	if err := s.logRepo.ClearAllLogs(ctx); err != nil {
		return err
	}
	if err := s.locationCache.ClearAll(ctx); err != nil {
		return err
	}
	if err := s.jobPlanRepo.ClearAllJobPlans(ctx); err != nil {
		return err
	}
	return nil
}
