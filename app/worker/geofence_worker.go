package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sistem_manajemen_armada/db"
	"sistem_manajemen_armada/model"

	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	db.Init()

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
		log.Fatalf("Gagal buka channel: %v", err)
	}
	defer ch.Close()

	err = ch.ExchangeDeclare(
		"fleet.events", // name
		"fanout",       // type
		true,           // durable
		false,          // auto-deleted
		false,          // internal
		false,          // no-wait
		nil,            // arguments
	)
	if err != nil {
		log.Fatalf("Gagal deklarasi exchange: %v", err)
	}

	q, err := ch.QueueDeclare(
		"geofence_alerts", // name
		true,              // durable
		false,             // delete when unused
		false,             // exclusive
		false,             // no-wait
		nil,               // arguments
	)
	if err != nil {
		log.Fatalf("Gagal deklarasi queue: %v", err)
	}

	err = ch.QueueBind(q.Name, "", "fleet.events", false, nil)
	if err != nil {
		log.Fatalf("Gagal bind queue: %v", err)
	}

	msgs, err := ch.Consume(
		q.Name,
		"",
		true,  // auto-ack
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil,
	)
	if err != nil {
		log.Fatalf("Gagal konsumsi queue: %v", err)
	}

	log.Println("Geofence worker aktif, menunggu pesan...")

	for msg := range msgs {
		var payload model.GeofenceEvent
		if err := json.Unmarshal(msg.Body, &payload); err != nil {
			log.Printf("Gagal decode pesan: %v", err)
			continue
		}
		log.Printf(" Event Geofence diterima: %+v", payload)
	}
}
