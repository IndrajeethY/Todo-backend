package router

import (
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"todo-backend/handlers"
	"todo-backend/middleware"
)

func Setup(r *gin.Engine, db *gorm.DB, jwtExpiry time.Duration) {
	api := r.Group("/api")
	handlers.RegisterAuthRoutes(api, jwtExpiry)
	protected := api.Group("/")
	protected.Use(middleware.JWTAuthMiddleware())
	handlers.RegisterTodoRoutes(protected.Group("/todos"), db)
}
