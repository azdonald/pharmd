package main

import (
	"log"

	"github.com/azdonald/pharmd/backend/server"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, using environment variables")
	}

	server.Run()
}
