package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type LocationPayload struct {
	VehicleID string  `json:"vehicle_id"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Timestamp int64   `json:"timestamp"`
}

func ListenAndStoreLocation(broker string, db *sql.DB) {
	opts := mqtt.NewClientOptions().
		AddBroker(broker).
		SetClientID("golang-mqtt-subscriber")

	opts.OnConnect = func(client mqtt.Client) {
		if token := client.Subscribe("/fleet/vehicle/+/location", 0, func(c mqtt.Client, m mqtt.Message) {
			topic := m.Topic()
			parts := strings.Split(topic, "/")
			if len(parts) != 5 || parts[1] != "fleet" || parts[2] != "vehicle" {
				log.Println("⚠️ Topik tidak sesuai format:", topic)
				return
			}

			var payload LocationPayload
			if err := json.Unmarshal(m.Payload(), &payload); err != nil {
				log.Println("JSON tidak valid:", err)
				return
			}

			if err := validatePayload(payload); err != nil {
				log.Println("Data tidak valid:", err)
				return
			}

			// Simpan ke database
			_, err := db.Exec(`
                INSERT INTO vehicle_locations (vehicle_id, latitude, longitude, timestamp)
                VALUES ($1, $2, $3, $4)
                ON CONFLICT (vehicle_id, timestamp) DO NOTHING;
            `, payload.VehicleID, payload.Latitude, payload.Longitude, payload.Timestamp)

			if err != nil {
				log.Println("Gagal insert DB:", err)
			} else {
				log.Printf("Data tersimpan: %+v\n", payload)
			}
		}); token.Wait() && token.Error() != nil {
			log.Fatal("Gagal subscribe:", token.Error())
		}
	}

	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		log.Fatal("Gagal koneksi MQTT:", token.Error())
	}

	select {} // keep app running
}

func validatePayload(data LocationPayload) error {
	if data.VehicleID == "" {
		return fmt.Errorf("vehicle_id kosong")
	}
	if data.Latitude < -90 || data.Latitude > 90 {
		return fmt.Errorf("latitude tidak valid")
	}
	if data.Longitude < -180 || data.Longitude > 180 {
		return fmt.Errorf("longitude tidak valid")
	}
	if data.Timestamp <= 0 {
		return fmt.Errorf("timestamp tidak valid")
	}
	return nil
}
