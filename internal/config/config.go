package config

import "time"

type AppConfig struct {
	Address string
}

type ClientsConfig struct {
	PostgresUser     string
	PostgresPassword string
	PostgresHost     string
	PostgresPort     int
	PostgresDB       string

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

	c.App = AppConfig{
		Address: LookupEnv("APP_ADDRESS", "127.0.0.1:5000"),
	}

	c.Clients = ClientsConfig{
		PostgresUser:     LookupEnv("POSTGRES_USER", "postgres"),
		PostgresPassword: LookupEnv("POSTGRES_PASSWORD", ""),
		PostgresHost:     LookupEnv("POSTGRES_HOST", "127.0.0.1"),
		PostgresPort:     LookupEnv("POSTGRES_PORT", 5432),
		PostgresDB:       LookupEnv("POSTGRES_DB", "test"),

		RedisAddress:      LookupEnv("REDIS_ADDRESS", "redis://127.0.0.1:6379"),
		RedisPassword:     LookupEnv("REDIS_PASSWORD", ""),
		RedisDB:           LookupEnv("REDIS_DB", 0),
		RedisWriteTimeout: time.Duration(LookupEnv("REDIS_WRITE_TIMEOUT", 0)) * time.Second,
		RedisReadTimeout:  time.Duration(LookupEnv("REDIS_READ_TIMEOUT", 0)) * time.Second,
	}

	return &c
}