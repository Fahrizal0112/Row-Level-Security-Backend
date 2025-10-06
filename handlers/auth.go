package handlers

import (
    "net/http"

    "row-level-security-backend/models"
    "row-level-security-backend/utils"

    "github.com/gin-gonic/gin"
    "gorm.io/gorm"
)

type AuthHandler struct {
    DB *gorm.DB
}

type LoginRequest struct {
    Email    string `json:"email" binding:"required,email"`
    Password string `json:"password" binding:"required"`
}

type RegisterRequest struct {
    Email    string `json:"email" binding:"required,email"`
    Password string `json:"password" binding:"required,min=6"`
    Name     string `json:"name" binding:"required"`
    TenantID uint   `json:"tenant_id" binding:"required"`
}

func NewAuthHandler(db *gorm.DB) *AuthHandler {
    return &AuthHandler{DB: db}
}

func (h *AuthHandler) Login(c *gin.Context) {
    var req LoginRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    var user models.User
    if err := h.DB.Preload("Tenant").Where("email = ?", req.Email).First(&user).Error; err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
        return
    }

    if !utils.CheckPasswordHash(req.Password, user.Password) {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
        return
    }

    token, err := utils.GenerateJWT(user.ID, user.TenantID, user.Role)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "token": token,
        "user":  user,
    })
}

func (h *AuthHandler) Register(c *gin.Context) {
    var req RegisterRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    var tenant models.Tenant
    if err := h.DB.First(&tenant, req.TenantID).Error; err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tenant"})
        return
    }

    hashedPassword, err := utils.HashPassword(req.Password)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
        return
    }

    user := models.User{
        Email:    req.Email,
        Password: hashedPassword,
        Name:     req.Name,
        TenantID: req.TenantID,
        Role:     "user",
    }

    if err := h.DB.Create(&user).Error; err != nil {
        c.JSON(http.StatusConflict, gin.H{"error": "User already exists"})
        return
    }

    token, err := utils.GenerateJWT(user.ID, user.TenantID, user.Role)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
        return
    }

    c.JSON(http.StatusCreated, gin.H{
        "token": token,
        "user":  user,
    })
}