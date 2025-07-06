package mqtt

import (
	"encoding/json"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"log"
	"regexp"
	"sistem_manajemen_armada/model"
	"sistem_manajemen_armada/repository"
)

type Subscriber struct {
	Repo repository.VehicleRepository
}

func NewSubscriber(repo repository.VehicleRepository) *Subscriber {
	return &Subscriber{Repo: repo}
}

func (s *Subscriber) Start(broker string) {
	opts := mqtt.NewClientOptions().AddBroker(broker)
	opts.SetClientID("vehicle-subscriber")

	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		log.Fatalf("MQTT connection error: %v", token.Error())
	}

	topic := "/fleet/vehicle/+/location"
	if token := client.Subscribe(topic, 1, s.messageHandler); token.Wait() && token.Error() != nil {
		log.Fatalf("MQTT subscribe error: %v", token.Error())
	}

	log.Println("MQTT subscriber listening on:", topic)
}

func (s *Subscriber) messageHandler(client mqtt.Client, msg mqtt.Message) {
	log.Println("MQTT Message Received:", string(msg.Payload()))

	// Extract vehicle_id from topic
	topic := msg.Topic()
	re := regexp.MustCompile(`/fleet/vehicle/(.+)/location`)
	match := re.FindStringSubmatch(topic)
	if len(match) < 2 {
		log.Println("Invalid topic format:", topic)
		return
	}
	vehicleID := match[1]

	// Parse JSON
	var payload struct {
		Latitude  float64 `json:"latitude"`
		Longitude float64 `json:"longitude"`
		Timestamp int64   `json:"timestamp"`
	}
	if err := json.Unmarshal(msg.Payload(), &payload); err != nil {
		log.Println("Invalid JSON:", err)
		return
	}

	// Create Location object
	loc := &model.VehicleLocation{
		VehicleID: vehicleID,
		Latitude:  payload.Latitude,
		Longitude: payload.Longitude,
		Timestamp: payload.Timestamp,
	}

	// Save to DB
	if err := s.Repo.SaveLocation(loc); err != nil {
		log.Println("Gagal simpan data ke DB:", err)
	} else {
		log.Printf("Lokasi kendaraan %s disimpan.\n", vehicleID)
	}
}
