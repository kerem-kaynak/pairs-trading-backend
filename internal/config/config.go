package config

import "os"

type Config struct {
	Port                 string
	DBUser               string
	DBPass               string
	DBName               string
	CloudSQLInstanceName string
	GoogleClientID       string
	GoogleClientSecret   string
	JWTSecret            string
}

func Load() (*Config, error) {
	return &Config{
		Port:                 getEnv("PORT", "8080"),
		DBUser:               os.Getenv("DB_USER"),
		DBPass:               os.Getenv("DB_PASS"),
		DBName:               os.Getenv("DB_NAME"),
		CloudSQLInstanceName: os.Getenv("CLOUD_SQL_INSTANCE_NAME"),
		GoogleClientID:       os.Getenv("GOOGLE_CLIENT_ID"),
		GoogleClientSecret:   os.Getenv("GOOGLE_CLIENT_SECRET"),
		JWTSecret:            os.Getenv("JWT_SECRET"),
	}, nil
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
