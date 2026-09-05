package services

import (
	"context"
	"log"
	"time"

	"trip-history-api/internal/core/domain"
	"trip-history-api/internal/core/ports"
)

type jobPlanService struct {
	repo         ports.JobPlanRepository
	userRepo     ports.UserRepository
	notifService ports.NotificationService
}

func NewJobPlanService(repo ports.JobPlanRepository, userRepo ports.UserRepository, notifService ports.NotificationService) ports.JobPlanService {
	return &jobPlanService{
		repo:         repo,
		userRepo:     userRepo,
		notifService: notifService,
	}
}

func (s *jobPlanService) CreateJobPlan(ctx context.Context, plan domain.JobPlan) (string, error) {
	plan.Status = "pending"
	plan.CreatedAt = time.Now()
	
	planID, err := s.repo.CreateJobPlan(ctx, plan)
	if err != nil {
		return "", err
	}

	// Fetch driver user to get FCM token
	driver, err := s.userRepo.GetUser(ctx, plan.DriverID)
	if err != nil {
		log.Printf("CreateJobPlan: Failed to get driver %s: %v", plan.DriverID, err)
	} else if driver.FCMToken != "" {
		errNotif := s.notifService.SendJobPlanNotification(ctx, driver.FCMToken, "งานใหม่เข้า!", "รับ-ส่ง: "+plan.PickupPoint.Address)
		if errNotif != nil {
			log.Printf("CreateJobPlan: Failed to send FCM: %v", errNotif)
		} else {
			log.Printf("CreateJobPlan: Successfully sent FCM to driver %s", plan.DriverID)
		}
	} else {
		log.Printf("CreateJobPlan: Driver %s has no FCMToken", plan.DriverID)
	}

	return planID, nil
}

func (s *jobPlanService) GetAllJobPlans(ctx context.Context) ([]domain.JobPlan, error) {
	return s.repo.GetJobPlans(ctx)
}

func (s *jobPlanService) GetJobPlansByDriver(ctx context.Context, driverID string) ([]domain.JobPlan, error) {
	return s.repo.GetJobPlansByDriver(ctx, driverID)
}

func (s *jobPlanService) UpdateJobPlanStatus(ctx context.Context, planID string, status string) error {
	return s.repo.UpdateJobPlanStatus(ctx, planID, status)
}
