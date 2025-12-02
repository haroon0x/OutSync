package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseUrl  string
	GeminiAPIKey string
}

func LoadConfig() Config {
	godotenv.Load()
	config := Config{}
	config.DatabaseUrl = os.Getenv("DATABASE_URL")
	config.GeminiAPIKey = os.Getenv("GEMINI_API_KEY")
	fmt.Println("Config Loaded: ", config)
	return config
}
