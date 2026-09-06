package firestore

import (
	"context"
	"time"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"
	"trip-history-api/internal/core/domain"
	"trip-history-api/internal/core/ports"
)

type tripRepository struct {
	client *firestore.Client
}

func NewTripRepository(client *firestore.Client) ports.TripRepository {
	return &tripRepository{
		client: client,
	}
}

func (r *tripRepository) CreateTrip(ctx context.Context, trip domain.Trip) (string, error) {
	docRef, _, err := r.client.Collection("trips").Add(ctx, trip)
	if err != nil {
		return "", err
	}
	return docRef.ID, nil
}

func (r *tripRepository) UpdateTrip(ctx context.Context, tripID string, updateData map[string]interface{}) error {
	var updates []firestore.Update
	for k, v := range updateData {
		updates = append(updates, firestore.Update{Path: k, Value: v})
	}
	_, err := r.client.Collection("trips").Doc(tripID).Update(ctx, updates)
	return err
}

func (r *tripRepository) GetTripsByUserID(ctx context.Context, userID string) ([]domain.Trip, error) {
	trips := make([]domain.Trip, 0) // Initialize as empty slice to avoid null JSON
	iter := r.client.Collection("trips").Where("userId", "==", userID).Documents(ctx)
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}
		var trip domain.Trip
		if err := doc.DataTo(&trip); err != nil {
			return nil, err
		}
		
		// Inject the document ID into the struct
		trip.ID = doc.Ref.ID
		
		trips = append(trips, trip)
	}
	return trips, nil
}

func (r *tripRepository) GetTripsHistory(ctx context.Context, userID string, startTime time.Time, endTime time.Time) ([]domain.Trip, error) {
	trips := make([]domain.Trip, 0)
	query := r.client.Collection("trips").
		Where("startTime", ">=", startTime).
		Where("startTime", "<=", endTime)

	if userID != "" {
		query = query.Where("userId", "==", userID)
	}

	iter := query.Documents(ctx)
		
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}
		var trip domain.Trip
		if err := doc.DataTo(&trip); err != nil {
			return nil, err
		}
		
		trip.ID = doc.Ref.ID
		trips = append(trips, trip)
	}
	return trips, nil
}

func (r *tripRepository) SaveLocationHistory(ctx context.Context, tripID string, loc domain.LocationUpdate) error {
	_, _, err := r.client.Collection("trips").Doc(tripID).Collection("locations").Add(ctx, loc)
	return err
}

func (r *tripRepository) GetLocationHistory(ctx context.Context, tripID string) ([]domain.LocationUpdate, error) {
	locations := make([]domain.LocationUpdate, 0)
	iter := r.client.Collection("trips").Doc(tripID).Collection("locations").OrderBy("timestamp", firestore.Asc).Documents(ctx)
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}
		var loc domain.LocationUpdate
		if err := doc.DataTo(&loc); err != nil {
			return nil, err
		}
		locations = append(locations, loc)
	}
	return locations, nil
}

func (r *tripRepository) SaveCheckIn(ctx context.Context, checkIn domain.CheckIn) error {
	_, _, err := r.client.Collection("trips").Doc(checkIn.TripID).Collection("checkins").Add(ctx, checkIn)
	return err
}

func (r *tripRepository) ReportIssue(ctx context.Context, issue domain.Issue) error {
	_, _, err := r.client.Collection("trips").Doc(issue.TripID).Collection("issues").Add(ctx, issue)
	return err
}

func (r *tripRepository) ClearAllTrips(ctx context.Context) error {
	// 1. Delete all locations across all trips (Collection Group query)
	iterLoc := r.client.CollectionGroup("locations").Documents(ctx)
	for {
		doc, err := iterLoc.Next()
		if err == iterator.Done || err != nil {
			break
		}
		doc.Ref.Delete(ctx)
	}

	// 2. Delete all checkins across all trips (Collection Group query)
	iterChk := r.client.CollectionGroup("checkins").Documents(ctx)
	for {
		doc, err := iterChk.Next()
		if err == iterator.Done || err != nil {
			break
		}
		doc.Ref.Delete(ctx)
	}

	// 3. Delete all issues across all trips (Collection Group query)
	iterIssues := r.client.CollectionGroup("issues").Documents(ctx)
	for {
		doc, err := iterIssues.Next()
		if err == iterator.Done || err != nil {
			break
		}
		doc.Ref.Delete(ctx)
	}

	// 4. Delete all trip documents
	iterTrips := r.client.Collection("trips").Documents(ctx)
	for {
		doc, err := iterTrips.Next()
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
