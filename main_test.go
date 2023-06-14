package main

import (
	"math"
	"net/http"
	"net/http/httptest"
	"testing"
)

const tolerance = 0.000001 // Tolerance value for floating-point comparisons

func TestGetSpotsHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/spots?latitude=40.7128&longitude=-74.0060&radius=1000", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(getSpotsHandler)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected status %d but got %d", http.StatusOK, rr.Code)
	}

	expectedResponse := `[{Latitude:40.71380678394081 Longitude:-74.00499321605918} {Latitude:40.71179321605919 Longitude:-74.00700678394082}]`
	if rr.Body.String() != expectedResponse {
		t.Errorf("expected response body %s but got %s", expectedResponse, rr.Body.String())
	}
}

func TestCalculateBoundingBox(t *testing.T) {
	latitude := 40.7128
	longitude := -74.0060
	radius := 1000.0

	minLat, maxLat, minLng, maxLng := calculateBoundingBox(latitude, longitude, radius)

	expectedMinLat := 40.703807
	expectedMaxLat := 40.721793
	expectedMinLng := -74.014993
	expectedMaxLng := -73.997007

	if !floatEquals(minLat, expectedMinLat, tolerance) ||
		!floatEquals(maxLat, expectedMaxLat, tolerance) ||
		!floatEquals(minLng, expectedMinLng, tolerance) ||
		!floatEquals(maxLng, expectedMaxLng, tolerance) {
		t.Errorf("expected bounding box coordinates (%f, %f, %f, %f) but got (%f, %f, %f, %f)",
			expectedMinLat, expectedMaxLat, expectedMinLng, expectedMaxLng, minLat, maxLat, minLng, maxLng)
	}
}

func TestFilterSpotsInCircle(t *testing.T) {
	centerLat := 40.7128
	centerLng := -74.0060
	radius := 1000.0

	spots := []Spot{
		{Latitude: 40.7028, Longitude: -73.996},
		{Latitude: 40.7228, Longitude: -74.016},
		{Latitude: 40.7328, Longitude: -74.026},
	}

	filteredSpots := filterSpotsInCircle(centerLat, centerLng, radius, spots)

	if len(filteredSpots) != 0 {
		t.Errorf("expected 2 spots in circle but got %d", len(filteredSpots))
	}

	expectedSpot1 := Spot{Latitude: 40.7028, Longitude: -73.996}
	expectedSpot2 := Spot{Latitude: 40.7228, Longitude: -74.016}

	if !containsSpot(filteredSpots, expectedSpot1) || !containsSpot(filteredSpots, expectedSpot2) {
		t.Errorf("filtered spots do not contain expected spots")
	}
}

func containsSpot(spots []Spot, target Spot) bool {
	for _, spot := range spots {
		if floatEquals(spot.Latitude, target.Latitude, tolerance) &&
			floatEquals(spot.Longitude, target.Longitude, tolerance) {
			return true
		}
	}
	return false
}


func TestCalculateDistance(t *testing.T) {
	lat1 := 40.7128
	lng1 := -74.0060
	lat2 := 40.7028
	lng2 := -73.996

	expectedDistance := 1395.322686

	distance := calculateDistance(lat1, lng1, lat2, lng2)

	if !floatEquals(distance, expectedDistance, tolerance) {
		t.Errorf("expected distance %f but got %f", expectedDistance, distance)
	}
}

func floatEquals(a, b, tolerance float64) bool {
	return math.Abs(a-b) <= tolerance
}
