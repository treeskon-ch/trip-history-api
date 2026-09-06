package http

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"trip-history-api/internal/adapters/handlers/ws"
	"trip-history-api/internal/core/domain"
	"trip-history-api/internal/core/ports"
)

type TripHandler struct {
	tripService ports.TripService
	trackingHub *ws.TrackingHub // Used for WebSocket broadcast
}

func NewTripHandler(tripService ports.TripService, trackingHub *ws.TrackingHub) *TripHandler {
	return &TripHandler{
		tripService: tripService,
		trackingHub: trackingHub,
	}
}

func (h *TripHandler) StartTrip(c *gin.Context) {
	var req struct {
		UserID string `json:"userId" binding:"required"`
		PlanID string `json:"planId"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tripID, err := h.tripService.StartTrip(c.Request.Context(), req.UserID, req.PlanID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create trip"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Trip started successfully", "tripId": tripID})
}

func (h *TripHandler) EndTrip(c *gin.Context) {
	tripID := c.Param("id")
	var req struct {
		Distance float64 `json:"distance"`
		ImageURL string  `json:"imageUrl"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	err := h.tripService.EndTrip(c.Request.Context(), tripID, req.Distance, req.ImageURL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to end trip"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "Trip ended"})
}

func (h *TripHandler) GetTrips(c *gin.Context) {
	userID := c.Query("userId")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "userId query param is required"})
		return
	}
	
	trips, err := h.tripService.GetTrips(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get trips"})
		return
	}
	c.JSON(http.StatusOK, trips)
}

func (h *TripHandler) GetHistory(c *gin.Context) {
	userID := c.Query("userId")
	if userID == "all" {
		userID = "" // Empty means fetch for all
	}
	startStr := c.Query("start")
	endStr := c.Query("end")
	if startStr == "" || endStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "start and end query params are required"})
		return
	}

	startTime, err := time.Parse(time.RFC3339, startStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid start time format, expected RFC3339"})
		return
	}
	endTime, err := time.Parse(time.RFC3339, endStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid end time format, expected RFC3339"})
		return
	}

	trips, err := h.tripService.GetTripsHistory(c.Request.Context(), userID, startTime, endTime)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get history trips: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, trips)
}

func (h *TripHandler) UpdateLocation(c *gin.Context) {
	tripID := c.Param("id")
	var req struct {
		UserID    string  `json:"userId" binding:"required"`
		Lat       float64 `json:"lat" binding:"required"`
		Lng       float64 `json:"lng" binding:"required"`
		Speed     float64 `json:"speed"`
		Timestamp string  `json:"timestamp"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	timestamp := req.Timestamp
	if timestamp == "" {
		timestamp = time.Now().Format(time.RFC3339)
	}

	loc := domain.LocationUpdate{
		Lat:       req.Lat,
		Lng:       req.Lng,
		Speed:     req.Speed,
		Timestamp: timestamp,
	}

	err := h.tripService.UpdateLocation(c.Request.Context(), tripID, loc)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update location"})
		return
	}

	// Broadcast to WebSocket
	if h.trackingHub != nil {
		// Create a dynamic map to include userId along with location data
		payloadMap := map[string]interface{}{
			"userId": req.UserID,
			"lat":    req.Lat,
			"lng":    req.Lng,
			"speed":  req.Speed,
		}
		payload, _ := json.Marshal(payloadMap)
		h.trackingHub.Broadcast <- ws.Message{
			UserID:  req.UserID,
			Payload: payload,
		}
	}

	c.JSON(http.StatusOK, gin.H{"status": "Location updated"})
}

func (h *TripHandler) GetTripRoute(c *gin.Context) {
	tripID := c.Param("id")
	locations, err := h.tripService.GetTripRoute(c.Request.Context(), tripID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get route"})
		return
	}
	c.JSON(http.StatusOK, locations)
}

func (h *TripHandler) CheckIn(c *gin.Context) {
	tripID := c.Param("id")
	var checkIn domain.CheckIn
	if err := c.ShouldBindJSON(&checkIn); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	checkIn.TripID = tripID

	err := h.tripService.CheckIn(c.Request.Context(), checkIn)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check-in"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "Check-in saved"})
}

func (h *TripHandler) ReportIssue(c *gin.Context) {
	tripID := c.Param("id")
	var req struct {
		UserID      string `json:"userId" binding:"required"`
		Title       string `json:"title" binding:"required"`
		Description string `json:"description"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	issue := domain.Issue{
		TripID:      tripID,
		Title:       req.Title,
		Description: req.Description,
		ReportedAt:  time.Now(),
	}

	err := h.tripService.ReportIssue(c.Request.Context(), issue)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to report issue"})
		return
	}

	// Broadcast issue to WebSocket
	if h.trackingHub != nil {
		payloadMap := map[string]interface{}{
			"type":        "issue",
			"userId":      req.UserID,
			"tripId":      tripID,
			"title":       req.Title,
			"description": req.Description,
			"reportedAt":  issue.ReportedAt.Format(time.RFC3339),
		}
		payload, _ := json.Marshal(payloadMap)
		h.trackingHub.Broadcast <- ws.Message{
			UserID:  req.UserID,
			Payload: payload,
		}
	}

	c.JSON(http.StatusOK, gin.H{"status": "Issue reported successfully"})
}
