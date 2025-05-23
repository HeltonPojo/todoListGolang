package routes

import (
	"todo-api/controllers"
	"todo-api/middleware"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterRoutes(r *gin.Engine, db *gorm.DB) {
	public := r.Group("/")
	{
		public.POST("/register", func(c *gin.Context) {
			controllers.Register(c, db)
		})
		public.POST("/login", func(c *gin.Context) {
			controllers.Login(c, db)
		})
	}

	protected := r.Group("/tasks")
	protected.Use(middleware.AuthMiddleware())
	{
	}
}
