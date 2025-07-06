package repository

import (
	"sistem_manajemen_armada/model"
)

type VehicleRepository interface {
	GetLastLocation(vehicleID string) (*model.VehicleLocation, error)
	GetLocationHistory(vehicleID string, start, end int64) ([]model.VehicleLocation, error)
	SaveLocation(loc *model.VehicleLocation) error
}
