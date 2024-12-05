package database

import (
	"errors"
	"time"
)

// GetWeightsByVin returns the weights of a vehicle given its vin
func (env *Database) GetWeightsByVin(params string) (*[]Weight, error) {
	time.Sleep(time.Millisecond * 750)

	if x, found := env.weight[params]; found {
		return &x, nil
	}

	return nil, errors.New("vehicle weights do not exist")
}
