package handlers

import (
    "net/http"
    "strconv"

    "row-level-security-backend/middleware"
    "row-level-security-backend/models"

    "github.com/gin-gonic/gin"
    "gorm.io/gorm"
)

type PostHandler struct {
    DB *gorm.DB
}

type CreatePostRequest struct {
    Title    string `json:"title" binding:"required"`
    Content  string `json:"content" binding:"required"`
    IsPublic bool   `json:"is_public"`
}

func NewPostHandler(db *gorm.DB) *PostHandler {
    return &PostHandler{DB: db}
}

func (h *PostHandler) CreatePost(c *gin.Context) {
    var req CreatePostRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    userID := c.GetUint("user_id")
    tenantID := c.GetUint("tenant_id")

    middleware.SetRLSContext(c, h.DB)

    post := models.Post{
        Title:    req.Title,
        Content:  req.Content,
        UserID:   userID,
        TenantID: tenantID,
        IsPublic: req.IsPublic,
    }

    if err := h.DB.Create(&post).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create post"})
        return
    }

    c.JSON(http.StatusCreated, post)
}

func (h *PostHandler) GetPosts(c *gin.Context) {
    tenantID := c.GetUint("tenant_id")
    userID := c.GetUint("user_id")

    middleware.SetRLSContext(c, h.DB)

    var posts []models.Post
    query := h.DB.Preload("User").Preload("Tenant")

    query = query.Where("tenant_id = ? AND (is_public = true OR user_id = ?)", tenantID, userID)

    if err := query.Find(&posts).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch posts"})
        return
    }

    c.JSON(http.StatusOK, posts)
}

func (h *PostHandler) GetPost(c *gin.Context) {
    id, err := strconv.ParseUint(c.Param("id"), 10, 32)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID"})
        return
    }

    tenantID := c.GetUint("tenant_id")
    userID := c.GetUint("user_id")

    middleware.SetRLSContext(c, h.DB)

    var post models.Post
    query := h.DB.Preload("User").Preload("Tenant")

    query = query.Where("id = ? AND tenant_id = ? AND (is_public = true OR user_id = ?)", id, tenantID, userID)

    if err := query.First(&post).Error; err != nil {
        if err == gorm.ErrRecordNotFound {
            c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
        } else {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch post"})
        }
        return
    }

    c.JSON(http.StatusOK, post)
}

func (h *PostHandler) UpdatePost(c *gin.Context) {
    id, err := strconv.ParseUint(c.Param("id"), 10, 32)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID"})
        return
    }

    var req CreatePostRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    tenantID := c.GetUint("tenant_id")
    userID := c.GetUint("user_id")

    middleware.SetRLSContext(c, h.DB)

    var post models.Post
    if err := h.DB.Where("id = ? AND tenant_id = ? AND user_id = ?", id, tenantID, userID).First(&post).Error; err != nil {
        if err == gorm.ErrRecordNotFound {
            c.JSON(http.StatusNotFound, gin.H{"error": "Post not found or access denied"})
        } else {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch post"})
        }
        return
    }

    post.Title = req.Title
    post.Content = req.Content
    post.IsPublic = req.IsPublic

    if err := h.DB.Save(&post).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update post"})
        return
    }

    c.JSON(http.StatusOK, post)
}

func (h *PostHandler) DeletePost(c *gin.Context) {
    id, err := strconv.ParseUint(c.Param("id"), 10, 32)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID"})
        return
    }

    tenantID := c.GetUint("tenant_id")
    userID := c.GetUint("user_id")

    middleware.SetRLSContext(c, h.DB)

    result := h.DB.Where("id = ? AND tenant_id = ? AND user_id = ?", id, tenantID, userID).Delete(&models.Post{})
    if result.Error != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete post"})
        return
    }

    if result.RowsAffected == 0 {
        c.JSON(http.StatusNotFound, gin.H{"error": "Post not found or access denied"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Post deleted successfully"})
}