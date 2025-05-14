package main

import (
	"log"
	"os"

	// Importe seu pacote de handlers
	"juntagrana-api/database" // Importe seu pacote de banco de dados
	"juntagrana-api/routes"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Erro ao carregar as variáveis de ambiente:", err)
	}

	db, err := database.ConnectDB()
	if err != nil {
		log.Fatal("Erro ao conectar ao banco de dados:", err)
	}
	defer database.CloseDB()

	router := gin.Default()

	// Defina suas rotas aqui
	routes.SetupRoutes(router, db)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Porta padrão se não estiver definida no .env
	}
	router.Run(":" + port)
}
