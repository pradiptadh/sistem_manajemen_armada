package main

import (
	"fmt"
	"os"
)

func main() {
	db := InitDB()
	defer db.Close()

	mqttHost := os.Getenv("MQTT_BROKER")
	mqttPort := os.Getenv("MQTT_PORT")
	if mqttPort == "" {
		mqttPort = "1883"
	}

	mqttURL := fmt.Sprintf("tcp://%s:%s", mqttHost, mqttPort)
	fmt.Println("🔌 Connecting to MQTT broker at", mqttURL)

	ListenAndStoreLocation(mqttURL, db)
}
