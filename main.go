package main

import (
	"log"
	"os"

	"be-porto-v3/database"
	"be-porto-v3/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	// Load environment variables
	database.LoadEnv()

	// Initialize database
	db, err := database.ConnectDB()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Setup Gin router
	r := gin.Default()

	// Melayani file statis dari folder uploads
	r.Static("/uploads", "./uploads")

	routes.SetupRoutes(r, db)

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Server running on port %s", port)
	r.Run(":" + port)
}
