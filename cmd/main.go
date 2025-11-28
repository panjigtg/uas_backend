package main

import (
	"uas/config"
	"path/filepath"
	"log"

	"github.com/joho/godotenv"
)

func main() {
	envPath, _ := filepath.Abs(".env")

	if err := godotenv.Load(envPath); err != nil {
		log.Fatal("Gagal load .env:", err)
	}

	app := config.Bootstrap()
	
	app.Listen(":3000")
}