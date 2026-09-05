package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"trip-history-api/internal/core/domain"
	"trip-history-api/internal/core/ports"
)

type JobPlanHandler struct {
	service ports.JobPlanService
}

func NewJobPlanHandler(service ports.JobPlanService) *JobPlanHandler {
	return &JobPlanHandler{
		service: service,
	}
}

func (h *JobPlanHandler) CreateJobPlan(c *gin.Context) {
	var plan domain.JobPlan
	if err := c.ShouldBindJSON(&plan); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	planID, err := h.service.CreateJobPlan(c.Request.Context(), plan)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create job plan: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Job plan created successfully",
		"planId":  planID,
	})
}

func (h *JobPlanHandler) GetJobPlans(c *gin.Context) {
	driverID := c.Query("driverId")
	
	var plans []domain.JobPlan
	var err error

	if driverID != "" {
		plans, err = h.service.GetJobPlansByDriver(c.Request.Context(), driverID)
	} else {
		plans, err = h.service.GetAllJobPlans(c.Request.Context())
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get job plans"})
		return
	}

	c.JSON(http.StatusOK, plans)
}

func (h *JobPlanHandler) UpdateJobPlanStatus(c *gin.Context) {
	planID := c.Param("id")
	if planID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "plan id is required"})
		return
	}

	var req struct {
		Status string `json:"status" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	err := h.service.UpdateJobPlanStatus(c.Request.Context(), planID, req.Status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update status"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "Job plan updated"})
}
