package main

import (
	"fmt"
	"log"

	"caretop-backend/config"
	"caretop-backend/database"
	"caretop-backend/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.Load()

	database.Connect(cfg)
	database.Migrate()

	r := gin.Default()

	routes.Setup(r)

	addr := fmt.Sprintf(":%s", cfg.ServerPort)
	log.Printf("Server starting on %s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
