package repository

import (
	"sistem_manajemen_armada/db"
	"sistem_manajemen_armada/model"
)

type GormRepo struct{}

func NewPostgresRepo() VehicleRepository {
	return &GormRepo{}
}

func (r *GormRepo) GetLastLocation(vehicleID string) (*model.VehicleLocation, error) {
	var loc model.VehicleLocation
	err := db.DB.
		Where("vehicle_id = ?", vehicleID).
		Order("timestamp DESC").
		First(&loc).Error
	if err != nil {
		return nil, err
	}
	return &loc, nil
}

func (r *GormRepo) GetLocationHistory(vehicleID string, start, end int64) ([]model.VehicleLocation, error) {
	var locations []model.VehicleLocation
	err := db.DB.
		Where("vehicle_id = ? AND timestamp BETWEEN ? AND ?", vehicleID, start, end).
		Order("timestamp ASC").
		Find(&locations).Error
	if err != nil {
		return nil, err
	}
	return locations, nil
}

func (r *GormRepo) SaveLocation(loc *model.VehicleLocation) error {
	return db.DB.Create(loc).Error
}
