package main

import (
	"fmt"
	"log"
	"math"
	"net/http"
	"strconv"
)

const (
	earthRadius = 6371000 // Earth radius in meters
	tolerance   = 0.000001 // Tolerance value for floating-point comparisons
)

type Spot struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

func main() {
	http.HandleFunc("/spots", getSpotsHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func getSpotsHandler(w http.ResponseWriter, r *http.Request) {
	latitudeStr := r.URL.Query().Get("latitude")
	longitudeStr := r.URL.Query().Get("longitude")
	radiusStr := r.URL.Query().Get("radius")

	latitude, err := strconv.ParseFloat(latitudeStr, 64)
	if err != nil {
		http.Error(w, "Invalid latitude", http.StatusBadRequest)
		return
	}

	longitude, err := strconv.ParseFloat(longitudeStr, 64)
	if err != nil {
		http.Error(w, "Invalid longitude", http.StatusBadRequest)
		return
	}

	radius, err := strconv.ParseFloat(radiusStr, 64)
	if err != nil || radius <= 0 {
		http.Error(w, "Invalid radius", http.StatusBadRequest)
		return
	}

	// Calculate the bounding box coordinates based on the circle's radius
	minLat, maxLat, minLng, maxLng := calculateBoundingBox(latitude, longitude, radius)

	// Generate some example spots within the bounding box
	spots := generateSpots(minLat, maxLat, minLng, maxLng)

	// Filter the spots within the circle
	spotsInCircle := filterSpotsInCircle(latitude, longitude, radius, spots)

	// Return the spots within the circle as a response
	fmt.Fprintf(w, "%+v", spotsInCircle)
}

func calculateBoundingBox(latitude, longitude, radius float64) (float64, float64, float64, float64) {
	// Convert radius from meters to degrees
	radiusInDegrees := radius / earthRadius * (180 / math.Pi)

	// Calculate the coordinates for the bounding box
	minLat := latitude - radiusInDegrees
	maxLat := latitude + radiusInDegrees
	minLng := longitude - radiusInDegrees
	maxLng := longitude + radiusInDegrees

	return minLat, maxLat, minLng, maxLng
}

func generateSpots(minLat, maxLat, minLng, maxLng float64) []Spot {
	// Here, you can implement your own logic to generate spots within the bounding box.
	// For simplicity, let's generate some example spots within the specified range.

	spots := []Spot{
		{Latitude: minLat + 0.01, Longitude: minLng + 0.01},
		{Latitude: minLat + 0.02, Longitude: minLng + 0.02},
		{Latitude: minLat + 0.03, Longitude: minLng + 0.03},
		{Latitude: maxLat - 0.01, Longitude: maxLng - 0.01},
		{Latitude: maxLat - 0.02, Longitude: maxLng - 0.02},
	}

	return spots
}

func filterSpotsInCircle(centerLat, centerLng, radius float64, spots []Spot) []Spot {
	spotsInCircle := make([]Spot, 0)

	for _, spot := range spots {
		distance := calculateDistance(centerLat, centerLng, spot.Latitude, spot.Longitude)

		// Check if the spot is within the circle
		if distance <= radius {
			spotsInCircle = append(spotsInCircle, spot)
		}
	}

	return spotsInCircle
}

func calculateDistance(lat1, lng1, lat2, lng2 float64) float64 {
	// Convert latitude and longitude from degrees to radians
	lat1Rad := lat1 * (math.Pi / 180)
	lat2Rad := lat2 * (math.Pi / 180)
	lng1Rad := lng1 * (math.Pi / 180)
	lng2Rad := lng2 * (math.Pi / 180)

	// Calculate the differences between coordinates in radians
	deltaLat := lat2Rad - lat1Rad
	deltaLng := lng2Rad - lng1Rad

	// Calculate the central angle between the points using the Haversine formula
	a := math.Pow(math.Sin(deltaLat/2), 2) + math.Cos(lat1Rad)*math.Cos(lat2Rad)*math.Pow(math.Sin(deltaLng/2), 2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	// Calculate the distance using the Earth's radius
	distance := earthRadius * c

	return distance
}
