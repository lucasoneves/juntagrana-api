package main

import (
	"log"
	"os"

	"juntagrana-api/database" // Import do seu pacote de banco de dados
	"juntagrana-api/routes"   // Import do seu pacote de rotas

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
	defer database.CloseDB() // Se CloseDB() em database.go não precisar de argumento

	router := gin.Default()
	routes.SetupRouter(router, db) // Passe o 'router' e o 'db' para SetupRouter

	port := os.Getenv("PORT")

	if port == "" {
		port = "8080" // Porta padrão se não estiver definida no .env
	}

	router.Run(":" + port)
}
