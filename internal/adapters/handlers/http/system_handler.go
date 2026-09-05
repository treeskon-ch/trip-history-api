package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"trip-history-api/internal/core/ports"
)

type SystemHandler struct {
	systemService ports.SystemService
}

func NewSystemHandler(systemService ports.SystemService) *SystemHandler {
	return &SystemHandler{
		systemService: systemService,
	}
}

func (h *SystemHandler) ClearAllData(c *gin.Context) {
	err := h.systemService.ClearAllData(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to clear all data: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "All data cleared successfully"})
}
