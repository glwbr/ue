package trips

import (
	"time"

	"uber-extractor/internal/uberapi"
)

type Trip struct {
	UUID              string     `json:"uuid"`
	BeginTime         time.Time  `json:"beginTime"`
	EndTime           time.Time  `json:"endTime"`
	Status            TripStatus `json:"status"`
	Fare              float64    `json:"fare"`
	Driver            string     `json:"driver"`
	VehicleType       string     `json:"vehicleType"`
	Distance          float64    `json:"distance"`
	Duration          float64    `json:"duration"`
	PickupAddress     string     `json:"pickupAddress"`
	DropoffAddress    string     `json:"dropoffAddress"`
	PickupLat         float64    `json:"pickupLat"`
	PickupLon         float64    `json:"pickupLon"`
	DropoffLat        float64    `json:"dropoffLat"`
	DropoffLon        float64    `json:"dropoffLon"`
	Rating            int        `json:"rating"`
	MapURL            string     `json:"mapUrl"`
	PickupLocationID  string     `json:"pickupLocationID"`
	DropoffLocationID string     `json:"dropoffLocationID"`
}

type TripSummary struct {
	Count      int
	TotalFare  float64
	Activities []uberapi.Activity
}
