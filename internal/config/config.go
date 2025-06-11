package config

import (
	"sync/atomic"
	"fmt"
	"os"
	"database/sql"
	"github.com/wexlerdev/chirpy/internal/database"
	"github.com/joho/godotenv"
	"log"
)

type ApiConfig struct {
	FileserverHits atomic.Int32
	DbQueries *database.Queries
	platform string
	JwtSecret	string
	PolkaKey	string
}

func NewConfig() *ApiConfig {
	var cfg ApiConfig
	err := godotenv.Load()
	if err != nil {
		// It's common to log this as a warning for development,
		// but not necessarily fatal, as env vars might be set externally in production.
		log.Println("Error loading .env file (this is fine if using external env vars):", err)
	}
	dbURL := os.Getenv("DB_URL")
	fmt.Printf("Value of DB_URL environment variable: '%s'\n", dbURL) // <-- Add this line

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		fmt.Println("conn err: ", err)
		os.Exit(1)
	}

	platformString := os.Getenv("PLATFORM")
	jwtSecret := os.Getenv("JWT_SECRET")
	polkaKey := os.Getenv("POLKA_KEY")

	dbQueries := database.New(db)
	cfg = ApiConfig{
		DbQueries: dbQueries,
		platform:		platformString,
		JwtSecret: jwtSecret,
		PolkaKey: polkaKey,
	}



	return &cfg

}

func (cfg * ApiConfig) GetPlatform() string {
	return cfg.platform
}
