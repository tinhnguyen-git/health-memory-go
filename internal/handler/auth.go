package handler

import (
	"github.com/redis/go-redis/v9"
	"github.com/tinhnguyen-git/health-memory-go/internal/config"
	"github.com/tinhnguyen-git/health-memory-go/internal/model"
	"github.com/tinhnguyen-git/health-memory-go/internal/repo"
	"github.com/tinhnguyen-git/health-memory-go/internal/service"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Handler struct definition added
type Handler struct {
	cfg       *config.Config
	authSvc   *service.AuthService
	userRepo  *repo.UserRepo
	redis     *redis.Client
	jwtSecret string
}

func NewHandler(cfg *config.Config) (*Handler, error) {
	db, err := gorm.Open(postgres.Open(cfg.PostgresDSN), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// migrate model, not repo
	if err := db.AutoMigrate(&model.User{}); err != nil {
		return nil, err
	}

	uRepo := repo.NewUserRepo(db)
	authSvc := service.NewAuthService(uRepo)

	rdb := redis.NewClient(&redis.Options{
		Addr: cfg.RedisAddr,
	})

	return &Handler{
		cfg:       cfg,
		authSvc:   authSvc,
		userRepo:  uRepo,
		redis:     rdb,
		jwtSecret: cfg.JwtSecret,
	}, nil
}
