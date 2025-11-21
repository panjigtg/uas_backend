package database

import (
    "database/sql"
    "log"
	"os"
	"fmt"

    "github.com/joho/godotenv"
    _ "github.com/lib/pq"
)

func PostgresConnections() *sql.DB {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Gagal memuat file .env")
	}
	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Can't Connect: %v", err)
	}

	if err := db.Ping(); err != nil {
		log.Fatalf("Can't Connect: %v", err)
	}
	log.Println("Succes Connect")
	return db
}
