package database

import (
	"errors"
	"time"
)

// GetVehiclesByClient returns a list of VINs associated with a client
func (env *Database) GetVehiclesByClient(client string) (*[]string, error) {
	time.Sleep(time.Millisecond * 750)

	var results []string
	for key := range env.vehicles {
		if env.vehicles[key].Client == client {
			results = append(results, env.vehicles[key].Vin)
		}
	}

	return &results, nil
}

// GetVehicleByName returns the vehilc given its vin
func (env *Database) GetVehicleByVin(params string) (*Vehicle, error) {
	time.Sleep(time.Millisecond * 750)

	if x, found := env.vehicles[params]; found {
		return &x, nil
	}

	return nil, errors.New("vehicle does not exist")
}
