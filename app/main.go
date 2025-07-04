package main

import (
	"sistem_manajemen_armada/db"
	"sistem_manajemen_armada/handler"
	"sistem_manajemen_armada/mqtt"
	"sistem_manajemen_armada/repository"

	"github.com/gin-gonic/gin"
)

func main() {
	db.Init()

	repo := repository.NewPostgresRepo()
	h := handler.NewVehicleHandler(repo)

	// Jalankan MQTT subscriber di goroutine
	go func() {
		sub := mqtt.NewSubscriber(repo)
		sub.Start("tcp://mosquitto:1883")
	}()

	r := gin.Default()
	r.GET("/vehicles/:vehicle_id/location", h.GetLastLocation)
	r.GET("/vehicles/:vehicle_id/history", h.GetHistory)

	r.Run(":8080")
}
