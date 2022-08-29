package config

import (
	"fmt"
	"time"
)

type AppConfig struct {
	Address      string
	JWTSecretKey string
	JWTDuration  time.Duration
	UploadsPath  string
}

type ClientsConfig struct {
	DatabaseURL string

	RedisAddress      string
	RedisPassword     string
	RedisDB           int
	RedisWriteTimeout time.Duration
	RedisReadTimeout  time.Duration
}

type Config struct {
	App     AppConfig
	Clients ClientsConfig
}

func New() *Config {
	var c Config

	port := LookupEnv("PORT", 5000)
	c.App = AppConfig{
		Address:      LookupEnv("APP_ADDRESS", fmt.Sprintf(":%d", port)),
		JWTSecretKey: LookupEnv("JWT_SECRET_KEY", "secret"),
		JWTDuration: LookupEnv("JWT_DURATION", time.Duration(
			time.Now().Add(time.Hour*24*30).Unix())),
		UploadsPath: LookupEnv("UPLOADS_PATH", "public/uploads"),
	}

	c.Clients = ClientsConfig{
		DatabaseURL: LookupEnv("DATABASE_URL", ""),

		RedisAddress:      LookupEnv("REDIS_ADDRESS", "redis://127.0.0.1:6379"),
		RedisPassword:     LookupEnv("REDIS_PASSWORD", ""),
		RedisDB:           LookupEnv("REDIS_DB", 0),
		RedisWriteTimeout: LookupEnv("REDIS_WRITE_TIMEOUT", time.Second),
		RedisReadTimeout:  LookupEnv("REDIS_READ_TIMEOUT", time.Second),
	}

	return &c
}
