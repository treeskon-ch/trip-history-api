package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"trip-history-api/internal/core/domain"
	"trip-history-api/internal/core/ports"
)

type LogHandler struct {
	logService ports.LogService
}

func NewLogHandler(logService ports.LogService) *LogHandler {
	return &LogHandler{
		logService: logService,
	}
}

func (h *LogHandler) SaveSystemLog(c *gin.Context) {
	var log domain.SystemLog
	if err := c.ShouldBindJSON(&log); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.logService.SaveLog(c.Request.Context(), log)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save log"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "Log saved"})
}
