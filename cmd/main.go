package main

import (
	"log"
	"os"

	"github.com/tinhnguyen-git/health-memory-go/internal/config"
	"github.com/tinhnguyen-git/health-memory-go/internal/http"
)

func main() {
	// load config (from env)
	cfg := config.LoadFromEnv()

	// start server
	srv := http.NewServer(cfg)
	if err := srv.Run(); err != nil {
		log.Fatalf("server exit: %v", err)
	}
}
