package controllers

import (
	"net/http"
	"strconv"
	"todo-api/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetTask(c *gin.Context, db *gorm.DB) {
	userID := c.GetUint("userID")

	var tasks []models.Tasks
	var total int64

	statusParam := c.Query("status")
	query := db.Where(" user_id = ?", userID)
	if statusParam != "" {
		query = query.Where("status = ?", statusParam)
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset := (page - 1) * limit

	query.Model(&models.Tasks{}).Count(&total)
	query.Order("created_at desc").Limit(limit).Offset(offset).Find(&tasks)

	c.JSON(http.StatusOK, gin.H{
		"data":  tasks,
		"total": total,
		"page":  page,
	})
}

func CreateTask(c *gin.Context, db *gorm.DB) {
	userID := c.GetUint("userID")

	var input struct {
		Title       string `json:"title" binding:"required"`
		Description string `json:"description"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Input"})
		return
	}

	task := models.Tasks{
		Title:       input.Title,
		Description: input.Description,
		UserID:      userID,
	}
	if err := db.Create(&task).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not creat task"})
		return
	}

	c.JSON(http.StatusCreated, task)
}

func UpdateTask(c *gin.Context, db *gorm.DB) {
	userID := c.GetUint("userID")
	id := c.Param("id")

	var task models.Tasks
	if err := db.Where("id = ? AND user_id = ?", id, userID).Find(&task).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Could not find the task"})
		return
	}

	var input struct {
		Title       *string `json:"title"`
		Description *string `json:"description"`
		Status      *string `json:"status"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Input"})
		return
	}

	if input.Title != nil {
		task.Title = *input.Title
	}
	if input.Description != nil {
		task.Description = *input.Description
	}
	if input.Status != nil {
		task.Status = *input.Status
	}

	db.Save(&task)

	c.JSON(http.StatusOK, task)
}

func DeleteTask(c *gin.Context, db *gorm.DB) {
	userID := c.GetUint("userID")
	id := c.Param("id")

	var task models.Tasks
	if err := db.Where("id = ? AND user_id = ?", id, userID).Find(&task).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Could not find the task"})
		return
	}

	db.Delete(&task)
	c.JSON(http.StatusOK, gin.H{"message": "Task Deleted"})
}
