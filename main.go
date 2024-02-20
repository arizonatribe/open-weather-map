package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Printf("Failed to load environment variables from '.env' file:\n\t%s\n", err)
	}

	apiKey := os.Getenv("OPENWEATHER_APIKEY")

	// This can also go in the .env file.
	// It isn't sensitive, but some URLs can change per environment and as a common practice are derived from env vars
	baseUrlOpenWeather := "https://api.openweathermap.org/data/2.5"

	routes := SetupRoutes(baseUrlOpenWeather, apiKey)
	routes.Run(":8080")
}
