package db

import (
	"errors"
	"time"

	"github.com/Maaaarko/go-gas/types"
)

type MemoryDatabase struct {
	Users       map[string]types.User
	GasStations map[string]types.GasStation
	Histories   map[string]types.History
}

func NewMemoryDatabase() *MemoryDatabase {
	return &MemoryDatabase{
		Users:       make(map[string]types.User),
		GasStations: make(map[string]types.GasStation),
		Histories:   make(map[string]types.History),
	}
}

func (m *MemoryDatabase) CreateUser(u *types.User) error {
	m.Users[u.Email] = *u
	return nil
}

func (m *MemoryDatabase) CreateGasStation(g *types.GasStation) error {
	m.GasStations[g.Name] = *g
	return nil
}

func (m *MemoryDatabase) GetGasStation(name string) (*types.GasStation, error) {
	g, ok := m.GasStations[name]
	if !ok {
		return nil, errors.New("gas station not found")
	}
	return &g, nil
}

func (m *MemoryDatabase) AddPriceToGasStation(name string, fuelPrices map[string]float64) error {
	g, ok := m.GasStations[name]
	if !ok {
		return errors.New("gas station not found")
	}

	for fuel, price := range fuelPrices {
		g.Prices[fuel] = price
	}

	var record types.HistoryRecord
	record.Time = time.Now().Unix()
	record.Prices = fuelPrices

	m.Histories[name] = append(m.Histories[name], record)

	return nil
}

func (m *MemoryDatabase) GetAllUsers() map[string]types.User {
	return m.Users
}

func (m *MemoryDatabase) GetAllGasStations() map[string]types.GasStation {
	return m.GasStations
}

func (m *MemoryDatabase) GetAllHistories() map[string]types.History {
	return m.Histories
}
