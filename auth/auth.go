package auth

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

const SessionName = "user_session" // Definição da constante SessionName

// GoogleAuthConfig configura os parâmetros de autenticação do Google OAuth 2.0
var GoogleAuthConfig *oauth2.Config

func init() {
	GoogleAuthConfig = &oauth2.Config{
		ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),     // Obtido do Google Cloud Console
		ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"), // Obtido do Google Cloud Console
		RedirectURL:  os.Getenv("GOOGLE_REDIRECT_URL"),  // Sua URL de callback
		Scopes: []string{"https://www.googleapis.com/auth/userinfo.profile",
			"https://www.googleapis.com/auth/userinfo.email"},
		Endpoint: google.Endpoint,
	}
}

const googleUserInfoURL = "https://www.googleapis.com/oauth2/v3/userinfo"

// GoogleLoginHandler redireciona o usuário para a página de login do Google
func GoogleLoginHandler(c *gin.Context) {
	url := GoogleAuthConfig.AuthCodeURL("state") // "state" é usado para proteção CSRF
	c.Redirect(http.StatusTemporaryRedirect, url)
}

// GoogleCallbackHandler processa o callback do Google após o login
func GoogleCallbackHandler(c *gin.Context) {
	code := c.Query("code")
	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Código de autorização ausente"})
		return
	}

	token, err := GoogleAuthConfig.Exchange(context.Background(), code)
	if err != nil {
		log.Println("Erro ao trocar o código por token:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Falha ao trocar o código por token"})
		return
	}

	client := GoogleAuthConfig.Client(context.Background(), token)
	resp, err := client.Get(googleUserInfoURL)
	if err != nil {
		log.Println("Erro ao obter informações do usuário:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Falha ao obter informações do usuário"})
		return
	}
	defer resp.Body.Close()

	var googleUser struct {
		ID    string `json:"sub"`
		Name  string `json:"name"`
		Email string `json:"email"`
		// Adicione outros campos relevantes
	}
	if err := json.NewDecoder(resp.Body).Decode(&googleUser); err != nil {
		log.Println("Erro ao decodificar informações do usuário:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Falha ao decodificar informações do usuário"})
		return
	}

	// Salvar informações do usuário na sessão
	session := sessions.Default(c)
	session.Set("user_id", googleUser.ID)
	session.Set("user_name", googleUser.Name)
	session.Set("user_email", googleUser.Email)
	if err := session.Save(); err != nil {
		log.Println("Erro ao salvar a sessão:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Falha ao salvar a sessão"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Login do Google bem-sucedido", "user": googleUser})
	// TODO: Redirecionar para a página principal do seu aplicativo
}

// LogoutHandler limpa a sessão do usuário
func LogoutHandler(c *gin.Context) {
	session := sessions.Default(c)
	session.Clear()
	if err := session.Save(); err != nil {
		log.Println("Erro ao limpar a sessão:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Falha ao limpar a sessão"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Logout bem-sucedido"})
	// TODO: Redirecionar para a página de login
}

// AuthMiddleware é um middleware para verificar se o usuário está autenticado
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		userID := session.Get("user_id")
		if userID == nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Não autenticado"})
			return
		}
		c.Next()
	}
}
