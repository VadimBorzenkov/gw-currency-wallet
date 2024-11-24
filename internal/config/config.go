package config

import (
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
)

// Config содержит параметры конфигурации.
type Config struct {
	Port                   string
	DBHost                 string
	DBPort                 string
	DBUser                 string
	DBPass                 string
	DBName                 string
	JWTSecret              string
	AuthJWTPublicKeyPath   string
	AuthJWTPrivateKeyPath  string
	AccessTokenExpiration  time.Duration
	RefreshTokenExpiration time.Duration
	GRPCExchangeHost       string
	GRPCExchangePort       string
}

// LoadConfig загружает переменные конфигурации из файла .env.
func LoadConfig() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file:", err)
	}

	accessTokenExpiration, err := time.ParseDuration(os.Getenv("ACCESS_TOKEN_EXPIRATION"))
	if err != nil {
		log.Fatalf("Invalid ACCESS_TOKEN_EXPIRATION format: %v", err)
	}

	refreshTokenExpiration, err := time.ParseDuration(os.Getenv("REFRESH_TOKEN_EXPIRATION"))
	if err != nil {
		log.Fatalf("Invalid REFRESH_TOKEN_EXPIRATION format: %v", err)
	}

	return &Config{
		Port:                   os.Getenv("PORT"),
		DBHost:                 os.Getenv("DB_HOST"),
		DBPort:                 os.Getenv("DB_PORT"),
		DBUser:                 os.Getenv("DB_USER"),
		DBPass:                 os.Getenv("DB_PASSWORD"),
		DBName:                 os.Getenv("DB_NAME"),
		JWTSecret:              os.Getenv("JWT_SECRET"),
		AuthJWTPublicKeyPath:   os.Getenv("AUTH_JWT_PUBLIC_KEY_PATH"),
		AuthJWTPrivateKeyPath:  os.Getenv("AUTH_JWT_PRIVATE_KEY_PATH"),
		AccessTokenExpiration:  accessTokenExpiration,
		RefreshTokenExpiration: refreshTokenExpiration,
		GRPCExchangeHost:       os.Getenv("GRPC_EXCHANGE_HOST"),
		GRPCExchangePort:       os.Getenv("GRPC_EXCHANGE_PORT"),
	}, nil
}

// GetAuthJWTPublicKeyPath возвращает путь к публичному ключу JWT.
func (cfg *Config) GetAuthJWTPublicKeyPath() string {
	return cfg.AuthJWTPublicKeyPath
}

// GetAuthJWTPrivateKeyPath возвращает путь к приватному ключу JWT.
func (cfg *Config) GetAuthJWTPrivateKeyPath() string {
	return cfg.AuthJWTPrivateKeyPath
}

// GetAccessTokenExpiration возвращает срок действия Access токена.
func (cfg *Config) GetAccessTokenExpiration() time.Duration {
	return cfg.AccessTokenExpiration
}

// GetRefreshTokenExpiration возвращает срок действия Refresh токена.
func (cfg *Config) GetRefreshTokenExpiration() time.Duration {
	return cfg.RefreshTokenExpiration
}
