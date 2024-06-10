package db

import (
	"github.com/Maaaarko/go-gas/types"
)

type Databaser interface {
	CreateUser(u *types.User) error
	CreateGasStation(g *types.GasStation) error
	GetGasStation(name string) (*types.GasStation, error)
	AddPriceToGasStation(name string, fuelPrices map[string]float64) error
	GetAllUsers() map[string]types.User
	GetAllGasStations() map[string]types.GasStation
	GetAllHistories() map[string]types.History
}
