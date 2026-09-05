package services

import (
	"context"

	"trip-history-api/internal/core/domain"
	"trip-history-api/internal/core/ports"
)

type logService struct {
	logRepo ports.LogRepository
}

func NewLogService(logRepo ports.LogRepository) ports.LogService {
	return &logService{
		logRepo: logRepo,
	}
}

func (s *logService) SaveLog(ctx context.Context, log domain.SystemLog) error {
	return s.logRepo.SaveLog(ctx, log)
}
