package transform

import (
	"uber-extractor/internal/locations"
	"uber-extractor/internal/parser"
	"uber-extractor/internal/trips"
	"uber-extractor/internal/uberapi"
)

func ProcessTrip(response *uberapi.GetTripResponse, lp *locations.Processor) (trips.Trip, error) {
	tripData := response.Data.GetTrip.Trip

	trip := trips.Trip{
		UUID:        tripData.UUID,
		Status:      trips.ParseTripStatus(tripData.Status),
		Driver:      tripData.Driver,
		VehicleType: response.Data.GetTrip.Receipt.VehicleType,
		Rating:      parser.Rating(response.Data.GetTrip.Rating),
		Distance:    parser.Distance(response.Data.GetTrip.Receipt.Distance),
		MapURL:      response.Data.GetTrip.MapURL,
	}

	beginTime, err := parser.Time(tripData.BeginTripTime)
	if err == nil {
		trip.BeginTime = beginTime
	}

	endTime, err := parser.Time(tripData.DropoffTime)
	if err == nil {
		trip.EndTime = endTime
	}

	trip.Fare = parser.Fare(tripData.Fare)

	duration, err := parser.Duration(response.Data.GetTrip.Receipt.Duration)
	if err == nil {
		trip.Duration = duration.Minutes()
	}

	if len(tripData.Waypoints) > 0 {
		trip.PickupAddress = tripData.Waypoints[0]
		trip.PickupLat, trip.PickupLon = parser.ExtractCoordinates(trip.MapURL, 0)
	}

	if len(tripData.Waypoints) > 1 {
		trip.DropoffAddress = tripData.Waypoints[len(tripData.Waypoints)-1]
		trip.DropoffLat, trip.DropoffLon = parser.ExtractCoordinates(trip.MapURL, 1)
	}

	if tripData.DropoffTime == "" && trips.ParseTripStatus(tripData.Status) != trips.StatusCanceled && len(tripData.Waypoints) > 0 {
		trip.DropoffLat, trip.DropoffLon = parser.ExtractCoordinates(trip.MapURL, 0)
		trip.DropoffAddress = tripData.Waypoints[len(tripData.Waypoints)-1]
	}

	if lp != nil && trips.ParseTripStatus(tripData.Status) == trips.StatusCompleted {
		trip.PickupLocationID = lp.FindOrCreateLocation(trip.PickupAddress, trip.PickupLat, trip.PickupLon)
		trip.DropoffLocationID = lp.FindOrCreateLocation(trip.DropoffAddress, trip.DropoffLat, trip.DropoffLon)
	}

	return trip, nil
}
