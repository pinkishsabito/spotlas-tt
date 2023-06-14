package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/google/uuid"
)

func TestGetSpotsHandler(t *testing.T) {
	// Create a mock database connection
	db := &mockDB{}
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		GetSpotsHandler(w, r, db)
	})

	// Create a sample request with test parameters
	req, err := http.NewRequest("GET", "/spots?latitude=12.34&longitude=56.78&radius=100&type=circle", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Set the expected spots to be returned by the mock DB
	expectedSpots := []Spot{
		{
			ID:          uuid.MustParse("00000000-0000-0000-0000-000000000001"),
			Name:        "Spot 1",
			Website:     "https://www.spot1.com",
			Coordinates: "POINT(12.34 56.78)",
			Description: "Spot 1 description",
			Rating:      4.5,
		},
		{
			ID:          uuid.MustParse("00000000-0000-0000-0000-000000000002"),
			Name:        "Spot 2",
			Website:     "https://www.spot2.com",
			Coordinates: "POINT(12.35 56.79)",
			Description: "Spot 2 description",
			Rating:      3.8,
		},
	}
	db.MockGetSpots = func(lat, lon, radius float64) ([]Spot, error) {
		return expectedSpots, nil
	}

	// Perform the request
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	// Check the response status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v, want %v", status, http.StatusOK)
	}

	// Check the response content type
	expectedContentType := "application/json"
	if ct := rr.Header().Get("Content-Type"); ct != expectedContentType {
		t.Errorf("handler returned wrong content type: got %v, want %v", ct, expectedContentType)
	}

	// Parse the response JSON into Spot objects
	var responseSpots []Spot
	err = json.Unmarshal(rr.Body.Bytes(), &responseSpots)
	if err != nil {
		t.Fatalf("error parsing response JSON: %v", err)
	}

	// Compare the response spots with the expected spots
	if !reflect.DeepEqual(responseSpots, expectedSpots) {
		t.Errorf("handler returned unexpected spots: got %v, want %v", responseSpots, expectedSpots)
	}
}

// mockDB implements the DB interface with mock methods
type mockDB struct {
	MockGetSpots func(lat, lon, radius float64) ([]Spot, error)
}

// GetSpots retrieves spots from the mock database
func (db *mockDB) GetSpots(lat, lon, radius float64) ([]Spot, error) {
	if db.MockGetSpots != nil {
		return db.MockGetSpots(lat, lon, radius)
	}
	return nil, nil
}
