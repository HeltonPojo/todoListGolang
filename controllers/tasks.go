package controllers

import (
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
}
