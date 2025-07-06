package mqtt

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"os"
	"sistem_manajemen_armada/model"
	"sistem_manajemen_armada/repository"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	amqp "github.com/rabbitmq/amqp091-go"
)

type Subscriber struct {
	repo repository.VehicleRepository
}

func NewSubscriber(repo repository.VehicleRepository) *Subscriber {
	return &Subscriber{repo: repo}
}

func (s *Subscriber) Start(broker string) {
	opts := mqtt.NewClientOptions().AddBroker(broker)
	opts.SetClientID("vehicle-tracker-subscriber")

	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		log.Fatalf("Gagal konek MQTT: %v", token.Error())
	}

	log.Println("Terhubung ke broker MQTT:", broker)

	client.Subscribe("vehicle/location", 0, func(client mqtt.Client, msg mqtt.Message) {
		var loc model.VehicleLocation
		if err := json.Unmarshal(msg.Payload(), &loc); err != nil {
			log.Printf("Gagal unmarshal: %v", err)
			return
		}

		if err := s.repo.SaveLocation(&loc); err != nil {
			log.Println("Gagal simpan data ke DB:", err)
		} else {
			log.Printf("Lokasi kendaraan %s disimpan", loc.VehicleID)
		}

		// Titik geofence target
		geofenceLat := -6.2088
		geofenceLon := 106.8456
		if IsWithinGeofence(loc.Latitude, loc.Longitude, geofenceLat, geofenceLon, 50) {
			log.Println("Kendaraan memasuki geofence, kirim event ke RabbitMQ...")
			go PublishToRabbitMQ(loc)
		}
	})
}

// geoforce buat kurang dari 1km
func IsWithinGeofence(lat1, lon1, lat2, lon2 float64, radiusMeters float64) bool {
	const metersPerDegreeLat = 111_000.0
	const metersPerDegreeLon = 111_000.0

	dLat := (lat2 - lat1) * metersPerDegreeLat
	dLon := (lon2 - lon1) * metersPerDegreeLon * math.Cos(lat1*math.Pi/180)

	distance := math.Sqrt(dLat*dLat + dLon*dLon)
	return distance <= radiusMeters
}

func PublishToRabbitMQ(loc model.VehicleLocation) {

	host := os.Getenv("RABBITMQ_HOST")
	port := os.Getenv("RABBITMQ_PORT")
	user := os.Getenv("RABBITMQ_USER")
	pass := os.Getenv("RABBITMQ_PASS")

	conn, err := amqp.Dial(fmt.Sprintf("amqp://%s:%s@%s:%s/", user, pass, host, port))
	if err != nil {
		log.Printf("Gagal konek RabbitMQ: %v", err)
		return
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Printf("Gagal buka channel: %v", err)
		return
	}
	defer ch.Close()

	err = ch.ExchangeDeclare(
		"fleet.events", // exchange name
		"fanout",       // type
		true,           // durable
		false, false, false, nil,
	)
	if err != nil {
		log.Printf("Gagal deklarasi exchange: %v", err)
		return
	}

	body, _ := json.Marshal(map[string]interface{}{
		"vehicle_id": loc.VehicleID,
		"event":      "geofence_entry",
		"location": map[string]float64{
			"latitude":  loc.Latitude,
			"longitude": loc.Longitude,
		},
		"timestamp": loc.Timestamp,
	})

	err = ch.Publish("fleet.events", "", false, false,
		amqp.Publishing{ContentType: "application/json", Body: body})
	if err != nil {
		log.Printf("Gagal publish event: %v", err)
	}
}
