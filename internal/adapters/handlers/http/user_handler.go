package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"trip-history-api/internal/core/ports"
)

type UserHandler struct {
	userService ports.UserService
}

func NewUserHandler(userService ports.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

func (h *UserHandler) RegisterUser(c *gin.Context) {
	var req struct {
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
		Name     string `json:"name" binding:"required"`
		Role     string `json:"role"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, link, err := h.userService.RegisterUser(c.Request.Context(), req.Email, req.Password, req.Name, req.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register user: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "User registered", 
		"userId": user.UserID,
		"verificationLink": link,
	})
}

func (h *UserHandler) DeleteUser(c *gin.Context) {
	uid := c.Param("id")
	if uid == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "uid is required"})
		return
	}

	err := h.userService.DeleteUser(c.Request.Context(), uid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete user: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "user deleted successfully"})
}

func (h *UserHandler) UpdateFCMToken(c *gin.Context) {
	uid := c.Param("id")
	if uid == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "uid is required"})
		return
	}

	var req struct {
		Token string `json:"fcmToken" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payload"})
		return
	}

	err := h.userService.UpdateFCMToken(c.Request.Context(), uid, req.Token)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update FCM token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "token updated"})
}

func (h *UserHandler) GetAllUsers(c *gin.Context) {
	users, err := h.userService.GetAllUsers(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get users"})
		return
	}
	c.JSON(http.StatusOK, users)
}
