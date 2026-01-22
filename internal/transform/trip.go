package transform

import (
	"uber-extractor/internal/locations"
	"uber-extractor/internal/parser"
	"uber-extractor/internal/trips"
	"uber-extractor/internal/uberapi"
)

func ProcessTrip(resp *uberapi.GetTripResponse, lp *locations.Processor) (trips.Trip, error) {
	tripData := resp.Data.GetTrip.Trip

	trip := trips.Trip{
		UUID:        tripData.UUID,
		Status:      trips.ParseTripStatus(tripData.Status),
		Driver:      tripData.Driver,
		VehicleType: resp.Data.GetTrip.Receipt.VehicleType,
		Rating:      parser.Rating(resp.Data.GetTrip.Rating),
		Distance:    parser.Distance(resp.Data.GetTrip.Receipt.Distance),
		MapURL:      resp.Data.GetTrip.MapURL,
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

	duration, err := parser.Duration(resp.Data.GetTrip.Receipt.Duration)
	if err == nil {
		trip.Duration = duration.Minutes()
	}

	if len(tripData.Waypoints) > 0 {
		trip.PickupAddress = tripData.Waypoints[0]
		trip.PickupLat, trip.PickupLon, _ = parser.ExtractCoordinates(trip.MapURL, 0)
	}

	if len(tripData.Waypoints) > 1 {
		trip.DropoffAddress = tripData.Waypoints[len(tripData.Waypoints)-1]
		trip.DropoffLat, trip.DropoffLon, _ = parser.ExtractCoordinates(trip.MapURL, 1)
	}

	if tripData.DropoffTime == "" && trips.ParseTripStatus(tripData.Status) != trips.StatusCanceled && len(tripData.Waypoints) > 0 {
		trip.DropoffLat, trip.DropoffLon, _ = parser.ExtractCoordinates(trip.MapURL, 0)
		trip.DropoffAddress = tripData.Waypoints[len(tripData.Waypoints)-1]
	}

	if lp != nil && trips.ParseTripStatus(tripData.Status) == trips.StatusCompleted {
		trip.PickupLocationID = lp.FindOrCreateLocation(trip.PickupAddress, trip.PickupLat, trip.PickupLon)
		trip.DropoffLocationID = lp.FindOrCreateLocation(trip.DropoffAddress, trip.DropoffLat, trip.DropoffLon)
	}

	return trip, nil
}
