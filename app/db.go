package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

func InitDB() *sql.DB {
	host := os.Getenv("POSTGRES_HOST")
	port := os.Getenv("POSTGRES_PORT")
	user := os.Getenv("POSTGRES_USER")
	pass := os.Getenv("POSTGRES_PASSWORD")
	name := os.Getenv("POSTGRES_DB")

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, pass, name)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal("DB Connection failed:", err)
	}

	// Create table if not exists
	_, err = db.Exec(`
        CREATE TABLE IF NOT EXISTS vehicle_locations (
			  vehicle_id TEXT NOT NULL,
			  latitude DOUBLE PRECISION NOT NULL,
			  longitude DOUBLE PRECISION NOT NULL,
			  timestamp TIMESTAMP NOT NULL,
			  PRIMARY KEY (vehicle_id, timestamp)
		);
    `)
	if err != nil {
		log.Fatal("Gagal buat tabel:", err)
	}

	return db
}
