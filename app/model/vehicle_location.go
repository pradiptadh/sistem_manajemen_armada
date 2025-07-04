package model

type VehicleLocation struct {
	VehicleID string  `json:"vehicle_id" gorm:"primaryKey"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Timestamp int64   `json:"timestamp" gorm:"primaryKey"`
}
