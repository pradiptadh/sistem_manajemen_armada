package main

import (
	"encoding/json"
	"fmt"
	"github.com/eclipse/paho.mqtt.golang"
	"time"
)

type LocationMessage struct {
	VehicleID string  `json:"vehicle_id"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Timestamp int64   `json:"timestamp"`
}

func main() {
	opts := mqtt.NewClientOptions().AddBroker("tcp://mosquitto:1883")
	opts.SetClientID("publisher")

	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	defer client.Disconnect(250)

	for {
		msg := LocationMessage{
			VehicleID: "B1234XYZ",
			Latitude:  -6.2088,
			Longitude: 106.8456,
			Timestamp: time.Now().Unix(),
		}

		data, _ := json.Marshal(msg)
		token := client.Publish("/fleet/vehicle/B1234XYZ/location", 0, false, data)
		token.Wait()

		fmt.Println("Published mock location")
		time.Sleep(2 * time.Second)
	}
}
