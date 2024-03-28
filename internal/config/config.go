package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	DBConnect string
	Port int
	NatsUrl string
}

func Load() (*Config, error) {
	err := godotenv.Load("../../.env")
	if err != nil {
		return nil, err
	}
	port, err := strconv.Atoi(os.Getenv("PORT"))
	if err != nil {
		return nil, err
	}
	
	
	return &Config{
		DBConnect: os.Getenv("DB_CONNECT"),
		Port: port,
		NatsUrl: os.Getenv("NATS_URL"),
	}, nil
}