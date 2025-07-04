package handler

import (
	"net/http"
	"sistem_manajemen_armada/repository"
	"strconv"

	"github.com/gin-gonic/gin"
)

type VehicleHandler struct {
	Repo repository.VehicleRepository
}

func NewVehicleHandler(repo repository.VehicleRepository) *VehicleHandler {
	return &VehicleHandler{Repo: repo}
}

func (h *VehicleHandler) GetLastLocation(c *gin.Context) {
	vehicleID := c.Param("vehicle_id")

	location, err := h.Repo.GetLastLocation(vehicleID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Data tidak ditemukan"})
		return
	}
	c.JSON(http.StatusOK, location)
}

func (h *VehicleHandler) GetHistory(c *gin.Context) {
	vehicleID := c.Param("vehicle_id")
	startStr := c.Query("start")
	endStr := c.Query("end")

	start, err1 := strconv.ParseInt(startStr, 10, 64)
	end, err2 := strconv.ParseInt(endStr, 10, 64)

	if err1 != nil || err2 != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Query start dan end harus valid timestamp"})
		return
	}

	history, err := h.Repo.GetLocationHistory(vehicleID, start, end)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data riwayat"})
		return
	}
	c.JSON(http.StatusOK, history)
}
