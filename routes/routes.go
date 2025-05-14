package routes

import (
	"juntagrana-api/controllers"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRoutes(router *gin.Engine, db *gorm.DB) { // Ou *gorm.DB se estiver usando GORM
	router.GET("/api", controllers.Greetings)
}

// pingHandler é um handler de exemplo
