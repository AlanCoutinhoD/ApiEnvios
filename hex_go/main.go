package main

import (
	shippingInfrastructure "demo/src/shipping/infrastructure"
	"fmt"
	"log"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	defer shippingInfrastructure.CloseConnections()

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	host := os.Getenv("HOST")
	if host == "" {
		host = "localhost"
	}

	router := gin.Default()

	// Configuración de CORS
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{"*"}
	config.AllowCredentials = true

	router.Use(cors.New(config))

	// Inicializar infraestructura de shipping
	shippingInfrastructure.Init(router)

	serverAddr := fmt.Sprintf("%s:%s", host, port)
	log.Printf("🚀 Servidor corriendo en http://%s", serverAddr)
	log.Printf("📝 Endpoints disponibles:")
	log.Printf("   POST http://%s/shipping", serverAddr)

	router.Run(":" + port)
}
