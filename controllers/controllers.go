package controllers

import "github.com/gin-gonic/gin"

// SetupRoutes configura as rotas da API

func Greetings(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "Hello World!",
	})
}
