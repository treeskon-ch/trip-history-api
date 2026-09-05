package firestore

import (
	"context"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"
	"trip-history-api/internal/core/domain"
	"trip-history-api/internal/core/ports"
)

type logRepository struct {
	client *firestore.Client
}

func NewLogRepository(client *firestore.Client) ports.LogRepository {
	return &logRepository{
		client: client,
	}
}

func (r *logRepository) SaveLog(ctx context.Context, log domain.SystemLog) error {
	_, _, err := r.client.Collection("system_logs").Add(ctx, log)
	return err
}

func (r *logRepository) ClearAllLogs(ctx context.Context) error {
	iter := r.client.Collection("system_logs").Documents(ctx)
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return err
		}
		doc.Ref.Delete(ctx)
	}
	return nil
}
