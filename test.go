package main_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func TestGetSpotsHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/spots", nil)
	if err != nil {
		t.Fatal(err)
	}

	q := req.URL.Query()
	q.Add("latitude", "10.12345")
	q.Add("longitude", "20.54321")
	q.Add("radius", "100")
	q.Add("type", "circle")
	req.URL.RawQuery = q.Encode()

	rr := httptest.NewRecorder()

	router := mux.NewRouter()
	router.HandleFunc("/spots", main.GetSpotsHandler).Methods("GET")

	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var spots []main.Spot
	err = json.Unmarshal(rr.Body.Bytes(), &spots)
	assert.NoError(t, err, "Error unmarshaling JSON response")

	assert.Equal(t, 2, len(spots), "Unexpected number of spots returned")

	assert.Equal(t, "Spot 1", spots[0].Name, "Unexpected spot name")
	assert.Equal(t, "http://example.com", spots[0].Website, "Unexpected spot website")
	assert.Equal(t, "10.12345,20.54321", spots[0].Coordinates, "Unexpected spot coordinates")
	assert.Equal(t, "Description 1", spots[0].Description, "Unexpected spot description")
	assert.Equal(t, 4.5, spots[0].Rating, "Unexpected spot rating")
}

func TestSortSpots(t *testing.T) {
	spots := []main.Spot{
		{
			ID:          uuid.NewV4(),
			Name:        "Spot 1",
			Website:     "http://example.com",
			Coordinates: "10.12345,20.54321",
			Description: "Description 1",
			Rating:      4.5,
		},
		{
			ID:          uuid.NewV4(),
			Name:        "Spot 2",
			Website:     "http://example.com",
			Coordinates: "10.6789,20.9876",
			Description: "Description 2",
			Rating:      3.2,
		},
	}

	main.SortSpots(10.12345, 20.54321, spots)

	assert.Equal(t, "Spot 1", spots[0].Name, "Unexpected spot order")
	assert.Equal(t, "Spot 2", spots[1].Name, "Unexpected spot order")
}

func TestDistance(t *testing.T) {
	// Create some test spots
	spots := []main.Spot{
		{
			ID:          uuid.NewV4(),
			Name:        "Spot 1",
			Website:     "http://example.com",
			Coordinates: "10.12345,20.54321",
			Description: "Description 1",
			Rating:      4.5,
		},
		{
			ID:          uuid.NewV4(),
			Name:        "Spot 2",
			Website:     "http://example.com",
			Coordinates: "10.6789,20.9876",
			Description: "Description 2",
			Rating:      3.2,
		},
	}

	distance := main.Distance(spots[0], spots[1])

	expectedDistance := 0.0
	assert.InEpsilon(t, expectedDistance, distance, 0.0001, "Unexpected distance value")
}

func TestSortBy(t *testing.T) {
	spots := []main.Spot{
		{
			ID:          uuid.NewV4(),
			Name:        "Spot 1",
			Website:     "http://example.com",
			Coordinates: "10.12345,20.54321",
			Description: "Description 1",
			Rating:      4.5,
		},
		{
			ID:          uuid.NewV4(),
			Name:        "Spot 2",
			Website:     "http://example.com",
			Coordinates: "10.6789,20.9876",
			Description: "Description 2",
			Rating:      3.2,
		},
	}

	byRating := func(i, j int) bool {
		return spots[i].Rating > spots[j].Rating
	}

	main.SortBy(byRating).Sort(spots)

	assert.Equal(t, "Spot 1", spots[0].Name, "Unexpected spot order")
	assert.Equal(t, "Spot 2", spots[1].Name, "Unexpected spot order")
}
