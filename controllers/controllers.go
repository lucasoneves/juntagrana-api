package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// SetupRoutes configura as rotas da API

func Greetings(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Hello World!",
	})
}
