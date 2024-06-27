package database

import (
	"errors"
	"time"
)

// GetAllClients returns a list of all available clients
func (env *Database) GetAllClients() (*[]Client, error) {
	time.Sleep(time.Millisecond * 750)

	var results []Client
	for key := range env.clients {
		results = append(results, env.clients[key])
	}

	return &results, nil
}

// GetClientsByName returns the client given its name
func (env *Database) GetClientsByName(params string) (*Client, error) {
	time.Sleep(time.Millisecond * 750)

	if x, found := env.clients[params]; found {
		return &x, nil
	}

	return nil, errors.New("client does not exist")

}
