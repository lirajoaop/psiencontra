package main

import (
	"log"

	"github.com/joaop/psiencontra/api/config"
	"github.com/joaop/psiencontra/api/handler"
	"github.com/joaop/psiencontra/api/router"
)

func main() {
	if err := config.Init(); err != nil {
		log.Fatalf("Failed to initialize config: %v", err)
	}

	handler.Init()

	r := router.New()

	port := config.GetEnv("PORT", "8080")
	config.Log.Info.Printf("Server starting on port %s", port)

	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
