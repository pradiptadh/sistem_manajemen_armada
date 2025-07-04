package db

import (
	"fmt"
	"log"
	"os"
	"time"

	"sistem_manajemen_armada/model"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Init() {
	host := os.Getenv("POSTGRES_HOST")
	port := os.Getenv("POSTGRES_PORT")
	user := os.Getenv("POSTGRES_USER")
	pass := os.Getenv("POSTGRES_PASSWORD")
	name := os.Getenv("POSTGRES_DB")

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, pass, name)

	var err error
	var db *gorm.DB

	for i := 0; i < 10; i++ {
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err == nil {
			DB = db
			log.Println("Sukses konek ke DB")
			break
		}
		log.Printf("Gagal konek DB (percobaan %d/10): %v", i+1, err)
		time.Sleep(3 * time.Second)
	}

	if err != nil {
		log.Fatalf("Gagal konek DB setelah 10 percobaan: %v", err)
	}

	// 🔧 Auto migrate model
	if err := DB.AutoMigrate(&model.VehicleLocation{}); err != nil {
		log.Fatalf("Gagal migrate schema: %v", err)
	}
	log.Println("Sukses auto-migrate schema VehicleLocation")
}
