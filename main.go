package main

import (
	"fmt"
	"log"

	"caretop-backend/config"
	"caretop-backend/database"
	"caretop-backend/routes"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// 加载本地 .env；若不存在则继续使用系统环境变量和默认值
	_ = godotenv.Load()

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
