package handlers

import (
	"net/http"
	"row-level-security-backend/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type TenantHandler struct {
	DB *gorm.DB
}

type CreateTenantRequest struct {
	Name   string `json:"name" binding:"required"`
	Domain string `json:"domain" binding:"required"`
}

func NewTenantHandler(db *gorm.DB) *TenantHandler {
	return &TenantHandler{DB: db}
}

func (h *TenantHandler) CreateTenant(c *gin.Context) {
	var req CreateTenantRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	tenant := models.Tenant{
		Name:   req.Name,
		Domain: req.Domain,
	}

	if err := h.DB.Create(&tenant).Error; err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Tenant domain already exists"})
		return
	}

	var user models.User
	if err := h.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	user.TenantID = &tenant.ID
	user.Role = "admin"
	if err := h.DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"tenant": tenant,
		"user":   user,
	})
}

func (h *TenantHandler) GetMyTenant(c *gin.Context) {
	userID, _ := c.Get("user_id")
	
	var user models.User
	if err := h.DB.Preload("Tenant").First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	if user.TenantID == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User has no tenant"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"tenant": user.Tenant})
}