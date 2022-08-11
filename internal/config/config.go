package config

import "time"

type AppConfig struct {
	Address      string
	JWTSecretKey string
	JWTDuration  time.Duration
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
		Address:      LookupEnv("APP_ADDRESS", "127.0.0.1:5000"),
		JWTSecretKey: LookupEnv("JWT_SECRET_KEY", "secret"),
		JWTDuration: LookupEnv("JWT_DURATION", time.Duration(
			time.Now().Add(time.Hour*24*30).Unix())),
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
		RedisWriteTimeout: LookupEnv("REDIS_WRITE_TIMEOUT", time.Second),
		RedisReadTimeout:  LookupEnv("REDIS_READ_TIMEOUT", time.Second),
	}

	return &c
}
