package fly

import (
	"fmt"
	"golife/gps"
	"golife/utils"
	"math/rand"
)

type Position struct {
	Longitude float64 `json:"longitude"`
	Latitude  float64 `json:"latitude"`
}

func (p *Position) IsSame(other Position) bool {
	if p == nil {
		return false
	}
	return p.Longitude == other.Longitude && p.Latitude == other.Latitude
}

func (p *Position) ToString() string {
	if p == nil {
		return ""
	}
	return fmt.Sprintf("%f_%f", p.Longitude, p.Latitude)
}

// Move function to update the position towards the target position by a certain distance
func (p *Position) Move(target Position, meter float64) bool {
	if p == nil {
		return false
	}

	// Calculate the distance between the current position and the target
	distance := gps.CalculateDistance(p.Latitude, p.Longitude, target.Latitude, target.Longitude)

	// Calculate the azimuth (direction) towards the target
	azimuth := gps.CalculateAzimuth(p.Latitude, p.Longitude, target.Latitude, target.Longitude)

	// If the specified meter distance is greater than the distance to the target, impossible
	if meter >= distance {
		return false
	} else {
		// Move the current position towards the target by the specified distance
		newLatitude, newLongitude := gps.MoveTowards(p.Latitude, p.Longitude, azimuth, meter)

		// Update the position of the object
		p.Latitude = newLatitude
		p.Longitude = newLongitude
	}
	return true
}

func RandPosition(r ...*rand.Rand) Position {
	lat, long := utils.RandLatitudeLongitude(true, r...)
	return Position{
		Longitude: long,
		Latitude:  lat,
	}
}
