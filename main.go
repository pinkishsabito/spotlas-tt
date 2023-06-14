package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"sort"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"github.com/satori/go.uuid"
)

// Spot represents a spot in the database
type Spot struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Website     string    `json:"website"`
	Coordinates string    `json:"coordinates"`
	Description string    `json:"description"`
	Rating      float64   `json:"rating"`
}

// Database connection parameters
const (
	host     = "localhost"
	port     = 5432
	user     = "your_username"
	password = "your_password"
	dbname   = "your_database"
)

// GetSpotsHandler handles the GET /spots endpoint
func GetSpotsHandler(w http.ResponseWriter, r *http.Request) {
	// Retrieve the query parameters
	latStr := r.URL.Query().Get("latitude")
	lonStr := r.URL.Query().Get("longitude")
	radiusStr := r.URL.Query().Get("radius")
	recruitmentType := r.URL.Query().Get("type")

	// Parse query parameters
	lat, err := strconv.ParseFloat(latStr, 64)
	if err != nil {
		http.Error(w, "Invalid latitude", http.StatusBadRequest)
		return
	}
	lon, err := strconv.ParseFloat(lonStr, 64)
	if err != nil {
		http.Error(w, "Invalid longitude", http.StatusBadRequest)
		return
	}
	radius, err := strconv.ParseFloat(radiusStr, 64)
	if err != nil {
		http.Error(w, "Invalid radius", http.StatusBadRequest)
		return
	}

	// Connect to the database
	db, err := sql.Open("postgres", fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname))
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Query the database to get spots within the specified area
	var spots []Spot
	var query string

	switch recruitmentType {
	case "circle":
		query = "SELECT * FROM MY_TABLE WHERE ST_DWithin(coordinates, ST_SetSRID(ST_MakePoint($1, $2), 4326), $3)"
	case "square":
		query = "SELECT * FROM MY_TABLE WHERE coordinates && ST_MakeEnvelope($1-$3/111320, $2-$3/111320, $1+$3/111320, $2+$3/111320, 4326)"
	default:
		http.Error(w, "Invalid recruitment type", http.StatusBadRequest)
		return
	}

	rows, err := db.Query(query, lon, lat, radius)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	// Iterate over the result rows and populate the spots slice
	for rows.Next() {
		var spot Spot
		err := rows.Scan(&spot.ID, &spot.Name, &spot.Website, &spot.Coordinates, &spot.Description, &spot.Rating)
		if err != nil {
			log.Fatal(err)
		}
		spots = append(spots, spot)
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}

	// Sort the spots by distance and rating if distance is smaller than 50m
	sortSpots(lat, lon, spots)

	// Convert the spots slice to JSON
	jsonData, err := json.Marshal(spots)
	if err != nil {
		http.Error(w, "Error converting spots to JSON", http.StatusInternalServerError)
		return
	}

	// Set the response headers and write the JSON data
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
}

// Sorts the spots by distance and rating if distance is smaller than 50m
func sortSpots(lat, lon float64, spots []Spot) {
	distance := func(spot Spot) float64 {
		// TODO: Calculate the distance between lat/lon and spot.Coordinates
		// Implement your own distance calculation method here
		return 0
	}

	distanceComparator := func(i, j int) bool {
		distI := distance(spots[i])
		distJ := distance(spots[j])

		if distI < 50 && distJ < 50 {
			return spots[i].Rating > spots[j].Rating
		}

		return distI < distJ
	}

	// Sort the spots using the custom distanceComparator
	SortBy(distanceComparator).Sort(spots)
}

// SortBy is a generic type that allows sorting slices using a custom comparator
type SortBy func(i, j int) bool

// Sort sorts the slice using the custom comparator
func (by SortBy) Sort(slice []Spot) {
	sorter := &sorter{
		slice: slice,
		by:    by,
	}
	sort.Sort(sorter)
}

type sorter struct {
	slice []Spot
	by    func(i, j int) bool
}

func (s *sorter) Len() int {
	return len(s.slice)
}

func (s *sorter) Swap(i, j int) {
	s.slice[i], s.slice[j] = s.slice[j], s.slice[i]
}

func (s *sorter) Less(i, j int) bool {
	return s.by(i, j)
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/spots", GetSpotsHandler).Methods("GET")

	http.Handle("/", r)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
