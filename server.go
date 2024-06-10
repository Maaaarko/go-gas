package main

import (
	"encoding/json"
	"log"
	"math"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/Maaaarko/go-gas/db"
	"github.com/Maaaarko/go-gas/types"
)

type ServerConfig struct {
	Addr           string
	ServerCert     string
	ServerKey      string
	UpdateInterval time.Duration
}

type apiError struct {
	Err    string
	Status int
}

func (e *apiError) Error() string {
	return e.Err
}

type apiFunc func(http.ResponseWriter, *http.Request, db.Databaser) error

func makeHandler(fn apiFunc, db db.Databaser) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := fn(w, r, db); err != nil {
			if e, ok := err.(*apiError); ok {
				WriteJSON(w, e.Status, e)
				return
			}
			WriteJSON(w, http.StatusInternalServerError, apiError{Err: "Error", Status: http.StatusInternalServerError})
			return
		}
	}
}

func createUser(w http.ResponseWriter, r *http.Request, db db.Databaser) error {
	var body types.CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		return &apiError{Err: "Invalid JSON", Status: http.StatusBadRequest}
	}

	u := types.NewUser(body.Name, body.Email, body.Password)

	if err := db.CreateUser(u); err != nil {
		return &apiError{Err: "Error creating user", Status: http.StatusInternalServerError}
	}

	return WriteJSON(w, http.StatusCreated, u.ToResponse())
}

func getUsers(w http.ResponseWriter, r *http.Request, db db.Databaser) error {
	users := make([]types.UserResponse, 0, len(db.GetAllUsers()))
	for _, u := range db.GetAllUsers() {
		users = append(users, *u.ToResponse())
	}
	return WriteJSON(w, http.StatusOK, users)
}

func createGasStation(w http.ResponseWriter, r *http.Request, db db.Databaser) error {
	var body types.GasStation
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		return &apiError{Err: "Invalid JSON", Status: http.StatusBadRequest}
	}

	if err := db.CreateGasStation(&body); err != nil {
		return &apiError{Err: "Error creating gas station", Status: http.StatusInternalServerError}
	}

	return WriteJSON(w, http.StatusCreated, body)
}

func getNearbyGasStations(w http.ResponseWriter, r *http.Request, db db.Databaser) error {
	lat := r.URL.Query().Get("lat")
	lon := r.URL.Query().Get("lon")

	latFloat, err := strconv.ParseFloat(lat, 64)
	if err != nil {
		return &apiError{Err: "Invalid lat", Status: http.StatusBadRequest}
	}

	lonFloat, err := strconv.ParseFloat(lon, 64)
	if err != nil {
		return &apiError{Err: "Invalid lon", Status: http.StatusBadRequest}
	}

	gasStations := db.GetAllGasStations()
	distances := make(map[string]float64)

	for _, g := range db.GetAllGasStations() {
		dx := g.Location.Lat - latFloat
		dy := g.Location.Lon - lonFloat

		d := math.Sqrt(dx*dx + dy*dy)
		distances[g.Name] = d
	}

	nearby := make([]types.GasStation, 0, min(3, len(gasStations)))
	for i := 0; i < min(3, len(gasStations)); i++ {
		var minName string
		var minDist float64
		for name, dist := range distances {
			if minName == "" || dist < minDist {
				minName = name
				minDist = dist
			}
		}
		nearby = append(nearby, gasStations[minName])
		delete(distances, minName)
	}

	return WriteJSON(w, http.StatusOK, nearby)
}

func getHistory(w http.ResponseWriter, r *http.Request, db db.Databaser) error {
	name := r.PathValue("name")

	g, err := db.GetGasStation(name)
	if err != nil {
		return err
	}

	var gasStationWithHistory types.GasStationWithHistory
	gasStationWithHistory.GasStation = g
	gasStationWithHistory.History = db.GetAllHistories()[name]

	return WriteJSON(w, http.StatusOK, gasStationWithHistory)
}

func generatePriceLoop(db db.Databaser, updateInterval time.Duration) {
	for {
		for _, g := range db.GetAllGasStations() {
			fuelPrices := make(map[string]float64)
			for fuel := range g.Prices {
				price := g.Prices[fuel]
				price += (rand.Float64() - 0.5) * 0.1
				fuelPrices[fuel] = price
			}

			if err := db.AddPriceToGasStation(g.Name, fuelPrices); err != nil {
				log.Printf("Error adding price to gas station %s: %s", g.Name, err)
			}

			log.Printf("Added price to gas station %s: %v", g.Name, fuelPrices)
		}

		<-time.After(updateInterval)
	}
}

func runServer(config ServerConfig) error {
	db := db.NewMemoryDatabase()
	http.HandleFunc("POST /users", makeHandler(createUser, db))
	http.HandleFunc("GET /users", makeHandler(getUsers, db))

	http.HandleFunc("POST /gas-stations", makeHandler(createGasStation, db))
	http.HandleFunc("GET /gas-stations/nearby", makeHandler(getNearbyGasStations, db))
	http.HandleFunc("GET /gas-stations/{name}", makeHandler(getHistory, db))

	go generatePriceLoop(db, config.UpdateInterval)

	return http.ListenAndServeTLS(config.Addr, config.ServerCert, config.ServerKey, nil)
}

func main() {
	config := ServerConfig{
		Addr:           ":443",
		ServerCert:     "server.crt",
		ServerKey:      "server.key",
		UpdateInterval: 5 * time.Second,
	}

	err := runServer(config)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
