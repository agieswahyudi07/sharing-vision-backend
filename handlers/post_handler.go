package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"sharing-vision-backend/config"
	"sharing-vision-backend/models"
)

// CreatePost handles POST /article
func CreatePost(c *gin.Context) {
	var post models.Post
	if err := c.ShouldBindJSON(&post); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Validate fields
	if validationErrs := post.Validate(); validationErrs != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errors": validationErrs})
		return
	}

	// Save to DB
	if err := config.DB.Create(&post).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Response must be {}
	c.JSON(http.StatusCreated, gin.H{})
}

// GetPosts handles GET /article/:limit/:offset
func GetPosts(c *gin.Context) {
	limitStr := c.Param("limit")
	offsetStr := c.Param("offset")

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 0 {
		limit = 10
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil || offset < 0 {
		offset = 0
	}

	// Also support query param filtering by status (useful for frontend)
	statusFilter := c.Query("status")

	var posts []models.Post
	var total int64

	query := config.DB.Model(&models.Post{})
	if statusFilter != "" {
		query = query.Where("Status = ?", statusFilter)
	}

	// Count total
	if err := query.Count(&total).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Fetch page
	if err := query.Limit(limit).Offset(offset).Order("Created_date DESC").Find(&posts).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Expose total count header
	c.Header("X-Total-Count", strconv.FormatInt(total, 10))
	c.Header("Access-Control-Expose-Headers", "X-Total-Count")

	// Response must be array of articles
	c.JSON(http.StatusOK, posts)
}

// GetPostByID handles GET /article/:id
func GetPostByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	var post models.Post
	if err := config.DB.First(&post, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Article not found"})
		return
	}

	c.JSON(http.StatusOK, post)
}

// UpdatePost handles PUT/PATCH/POST /article/:id
func UpdatePost(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	var existingPost models.Post
	if err := config.DB.First(&existingPost, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Article not found"})
		return
	}

	var inputPost models.Post
	if err := c.ShouldBindJSON(&inputPost); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Validate fields
	if validationErrs := inputPost.Validate(); validationErrs != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errors": validationErrs})
		return
	}

	// Update columns
	existingPost.Title = inputPost.Title
	existingPost.Content = inputPost.Content
	existingPost.Category = inputPost.Category
	existingPost.Status = inputPost.Status

	if err := config.DB.Save(&existingPost).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Response must be {}
	c.JSON(http.StatusOK, gin.H{})
}

// DeletePost handles DELETE /article/:id (or POST /article/:id)
func DeletePost(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	var post models.Post
	if err := config.DB.First(&post, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Article not found"})
		return
	}

	if err := config.DB.Delete(&post).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Response must be {}
	c.JSON(http.StatusOK, gin.H{})
}
