package config

import (
	"log"
	"os"
)

type Config struct {
	Port            string
	PostgresDSN     string
	RedisAddr       string
	JwtSecret       string
	GoogleClientID  string
	GoogleSecret    string
	GoogleCallback  string
	FacebookClientID string
	FacebookSecret   string
	FacebookCallback string
}

func LoadFromEnv() *Config {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	pgHost := os.Getenv("POSTGRES_HOST")
	pgPort := os.Getenv("POSTGRES_PORT")
	if pgPort == "" {
		pgPort = "5432"
	}
	pgUser := os.Getenv("POSTGRES_USER")
	pgPass := os.Getenv("POSTGRES_PASSWORD")
	pgDB := os.Getenv("POSTGRES_DB")
	if pgHost == "" || pgUser == "" || pgPass == "" || pgDB == "" {
		log.Fatalf("postgres env required")
	}
	dsn := "host=" + pgHost + " port=" + pgPort + " user=" + pgUser + " password=" + pgPass + " dbname=" + pgDB + " sslmode=disable"

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatalf("JWT_SECRET is required")
	}

	return &Config{
		Port:            port,
		PostgresDSN:     dsn,
		RedisAddr:       getOrDefault("REDIS_ADDR", "redis:6379"),
		JwtSecret:       jwtSecret,
		GoogleClientID:  os.Getenv("GOOGLE_CLIENT_ID"),
		GoogleSecret:    os.Getenv("GOOGLE_CLIENT_SECRET"),
		GoogleCallback:  os.Getenv("GOOGLE_REDIRECT_URL"),
		FacebookClientID: os.Getenv("FACEBOOK_CLIENT_ID"),
		FacebookSecret:   os.Getenv("FACEBOOK_CLIENT_SECRET"),
		FacebookCallback: os.Getenv("FACEBOOK_REDIRECT_URL"),
	}
}

func getOrDefault(k, def string) string {
	v := os.Getenv(k)
	if v == "" {
		return def
	}
	return v
}
