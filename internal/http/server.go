package http

import (
	"fmt"
	"log"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/tinhnguyen-git/health-memory-go/internal/config"
	"github.com/tinhnguyen-git/health-memory-go/internal/handler"
	"github.com/tinhnguyen-git/health-memory-go/internal/model"
)

type Server struct {
	cfg *config.Config
}

func NewServer(cfg *config.Config) *Server {
	return &Server{cfg: cfg}
}

func (s *Server) Run() error {
	// init DB
	db, err := gorm.Open(postgres.Open(s.cfg.PostgresDSN), &gorm.Config{})
	if err != nil {
		return err
	}
	// AutoMigrate
	if err := db.AutoMigrate(&model.User{}); err != nil {
		log.Printf("migrate err: %v", err)
	}

	// init redis
	rdb := redis.NewClient(&redis.Options{
		Addr:     s.cfg.RedisAddr,
		Password: "",
		DB:       0,
	})
	if err := rdb.Ping(rdb.Context()).Err(); err != nil {
		return err
	}

	// gin
	r := gin.Default()
	r.Use(cors.Default())

	// create repos/services/handler
	h := handler.NewBasicHandler(s.cfg, db, rdb) // we'll provide simpler constructor below
	// register routes
	h.RegisterRoutes(r)

	// oauth endpoints (google)
	if s.cfg.GoogleClientID != "" && s.cfg.GoogleSecret != "" && s.cfg.GoogleCallback != "" {
		// register google endpoints by using oauth package directly (omitted here for brevity)
	}

	addr := fmt.Sprintf(":%s", s.cfg.Port)
	log.Printf("Listening on %s", addr)
	return r.Run(addr)
}
