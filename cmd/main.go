package main

import (
	"uas/config"
	"path/filepath"
	"log"
	_ "uas/cmd/docs"

	"github.com/joho/godotenv"
	"github.com/gofiber/swagger"
)

// @title           UAS Achievement API
// @version         1.0
// @description     API Sistem Prestasi Mahasiswa
// @BasePath        /api/v1
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	envPath, _ := filepath.Abs(".env")

	if err := godotenv.Load(envPath); err != nil {
		log.Fatal("Gagal load .env:", err)
	}

	app := config.Bootstrap()
	app.Get("/swagger/*", swagger.HandlerDefault)
	app.Listen(":3000")
}