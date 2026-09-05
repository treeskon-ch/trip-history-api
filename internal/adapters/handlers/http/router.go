package http

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"trip-history-api/internal/adapters/handlers/ws"
)

func SetupRouter(
	tripHandler *TripHandler,
	userHandler *UserHandler,
	logHandler *LogHandler,
	systemHandler *SystemHandler,
	jobPlanHandler *JobPlanHandler,
	trackingHub *ws.TrackingHub,
) *gin.Engine {
	r := gin.Default()

	// CORS config
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	api := r.Group("/api")
	{
		// Trips
		api.POST("/trips/start", tripHandler.StartTrip)
		api.POST("/trips/:id/end", tripHandler.EndTrip)
		api.POST("/trips/:id/checkin", tripHandler.CheckIn)
		api.GET("/trips", tripHandler.GetTrips)
		api.POST("/trips/:id/locations", tripHandler.UpdateLocation)
		api.GET("/trips/:id/route", tripHandler.GetTripRoute)
		api.POST("/trips/:id/issues", tripHandler.ReportIssue)
		api.GET("/history", tripHandler.GetHistory)

		// Users & Auth
		api.POST("/users/register", userHandler.RegisterUser)
		api.GET("/users", userHandler.GetAllUsers)
		api.DELETE("/users/:id", userHandler.DeleteUser)
		api.PATCH("/users/:id/fcm-token", userHandler.UpdateFCMToken)

		// Logs
		api.POST("/logs", logHandler.SaveSystemLog)

		// System & Data Clearing
		api.POST("/system/clear-data", systemHandler.ClearAllData)

		// Job Plans
		api.POST("/plans", jobPlanHandler.CreateJobPlan)
		api.GET("/plans", jobPlanHandler.GetJobPlans)
		api.PATCH("/plans/:id/status", jobPlanHandler.UpdateJobPlanStatus)
	}

	// WebSockets
	wsGroup := r.Group("/ws")
	{
		wsGroup.GET("/tracking/:userId", trackingHub.ServeWS)
	}

	return r
}
