package routes

import (
	"juntagrana-api/auth"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin" // Ou "gorm.io/gorm" se estiver usando GORM
	"gorm.io/gorm"
)

// SetupRouter configura todas as rotas da aplicação
func SetupRouter(router *gin.Engine, db *gorm.DB) { // Ou *gorm.DB se estiver usando GORM
	// Configuração das sessões
	store := cookie.NewStore([]byte("secret")) // Substitua "secret" por uma chave secreta forte
	router.Use(sessions.Sessions(auth.SessionName, store))

	router.GET("/api/auth/google/login", auth.GoogleLoginHandler)
	router.GET("/api/auth/google/callback", auth.GoogleCallbackHandler)
	router.POST("/api/auth/logout", auth.LogoutHandler)

	// Rotas protegidas por autenticação (exemplo)
	protected := router.Group("/api/protected")
	protected.Use(auth.AuthMiddleware())
	{
		protected.GET("/dashboard", func(c *gin.Context) {
			session := sessions.Default(c)
			userID := session.Get("user_id")
			userName := session.Get("user_name")
			c.JSON(http.StatusOK, gin.H{"message": "Dashboard", "user_id": userID, "user_name": userName})
		})
		// Adicione outras rotas protegidas aqui
	}
}
