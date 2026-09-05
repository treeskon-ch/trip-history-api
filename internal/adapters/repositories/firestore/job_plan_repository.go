package firestore

import (
	"context"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"
	"trip-history-api/internal/core/domain"
	"trip-history-api/internal/core/ports"
)

type jobPlanRepository struct {
	client *firestore.Client
}

func NewJobPlanRepository(client *firestore.Client) ports.JobPlanRepository {
	return &jobPlanRepository{
		client: client,
	}
}

func (r *jobPlanRepository) CreateJobPlan(ctx context.Context, plan domain.JobPlan) (string, error) {
	docRef, _, err := r.client.Collection("job_plans").Add(ctx, plan)
	if err != nil {
		return "", err
	}
	return docRef.ID, nil
}

func (r *jobPlanRepository) GetJobPlans(ctx context.Context) ([]domain.JobPlan, error) {
	var plans []domain.JobPlan
	iter := r.client.Collection("job_plans").OrderBy("createdAt", firestore.Desc).Documents(ctx)
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}
		var plan domain.JobPlan
		if err := doc.DataTo(&plan); err == nil {
			plan.PlanID = doc.Ref.ID
			plans = append(plans, plan)
		}
	}
	return plans, nil
}

func (r *jobPlanRepository) GetJobPlansByDriver(ctx context.Context, driverID string) ([]domain.JobPlan, error) {
	var plans []domain.JobPlan
	iter := r.client.Collection("job_plans").Where("driverId", "==", driverID).OrderBy("createdAt", firestore.Desc).Documents(ctx)
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}
		var plan domain.JobPlan
		if err := doc.DataTo(&plan); err == nil {
			plan.PlanID = doc.Ref.ID
			plans = append(plans, plan)
		}
	}
	return plans, nil
}

func (r *jobPlanRepository) UpdateJobPlanStatus(ctx context.Context, planID string, status string) error {
	_, err := r.client.Collection("job_plans").Doc(planID).Update(ctx, []firestore.Update{
		{Path: "status", Value: status},
	})
	return err
}

func (r *jobPlanRepository) ClearAllJobPlans(ctx context.Context) error {
	iter := r.client.Collection("job_plans").Documents(ctx)
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
