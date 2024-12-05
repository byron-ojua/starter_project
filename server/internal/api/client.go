package api

import (
	"fmt"
	"net/http"
	"sort"
	"sync"

	"github.com/byron-ojua/starter-project/internal/database"
	"github.com/gin-gonic/gin"
)

// ClientInfo is a struct that represents a client and the number of vehicles they have.
type ClientInfo struct {
	Name         string `json:"name"`               // The name of the client
	ContactName  string `json:"contact_name"`       // The name of the contact person for the client
	ContactEmail string `json:"contact_email"`      // The email of the contact person for the client
	NumVehicles  int    `json:"number_of_vehicles"` // The number of vehicles the client has
}

// BasicVehicle is a struct that represents a vehicle and its basic information.
type BasicVehicle struct {
	Vin           string `json:"vin"`            // The VIN of the vehicle
	Mileage       int    `json:"mileage"`        // The mileage of the vehicle
	LargestWeight int    `json:"largest_weight"` // The largest weight the vehicle has carried
}

// ClientVehicles is a struct that represents a client and their vehicles.
type ClientVehicles struct {
	Name     string         `json:"name"`     // The name of the client
	Vehicles []BasicVehicle `json:"vehicles"` // The vehicles the client has
}

type getAllClientsOutput struct {
	Clients   []ClientInfo `json:"data"`      // The list of clients
	next_page string       `json:"next_page"` // The URL for the next page of clients
}

// This example route demonstrates how to use Goroutines to speed up the process of getting data from the database.
// This uses a sync.Map to store the client information and a channel to communicate errors between Goroutines.
// getClients retrieves all clients and the number of vehicles they have.
// @Summary Get all clients
// @Description Get all clients and the number of vehicles they have
// @Produce json
// @Tags Clients
// @Success 200 {array} getAllClientsOutput
// @Failure 404 {object} ErrorsResponse "Clients not found"
// @Failure 500 {object} ErrorsResponse "Internal Server Error"
// @Router /clients [get]
func (env *env) getClients(c *gin.Context) {
	var wg sync.WaitGroup

	// This is a thread-safe map that will store client information. Thread-safe maps are
	// safe to use in Goroutines. Another approach to get a similar result is to use a
	// a channel to communicate between Goroutines, similar to the err channel in this code.
	var clientInfo = sync.Map{}

	var output getAllClientsOutput

	clients, err := env.db.GetAllClients()
	output.next_page = "something" // This is a placeholder for the next page URL. This is usally returned in a database response.

	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusNotFound, ErrorsResponse{Errors: []string{"could not retrieve clients"}})
		return
	}

	// Channels are used to communicate between Goroutines. In this case, we are using a channel to
	// communicate errors between Goroutines. This is a common pattern in Go.
	var errChan = make(chan error, len(*clients))

	// Use Goroutines to speed up the process of getting the number of vehicles for each client.
	// This creates a bar bones goroutine pool. You can use a library like "ants" to create a more
	// sophisticated goroutine pool.
	for _, client := range *clients {
		wg.Add(1)
		mClient := client // When using Goroutines in a for loop, you need to create a new variable to avoid the loop variable being overwritten.

		// Run the Goroutine
		go func() {
			defer wg.Done()

			vehiclesByClient, err := env.db.GetVehiclesByClient(mClient.Name)

			if err != nil {
				errChan <- err
				return
			}

			cInfo := ClientInfo{
				Name:         mClient.Name,
				ContactName:  mClient.ContactName,
				ContactEmail: mClient.ContactEmail,
				NumVehicles:  len(*vehiclesByClient),
			}

			clientInfo.Store(mClient.Name, cInfo)
		}()
	}

	// Wait for all Goroutines to finish
	wg.Wait()

	// Close the error channel. After closing the channel, you can no longer send values to it,
	// but you can still receive values from it.
	close(errChan)

	// Loop over the error channel to check for any errors that occurred during the Goroutines.
	// Range over a channel will block until the channel is closed.
	for err := range errChan {
		fmt.Println(err)
	}

	// Load the client information from the sync.Map and append it to the allClients slice.
	clientInfo.Range(func(key, value interface{}) bool {
		output.Clients = append(output.Clients, value.(ClientInfo))
		return true
	})

	// sort the clients by name
	sort.Slice(output.Clients, func(i, j int) bool {
		return output.Clients[i].Name < output.Clients[j].Name
	})

	c.JSON(http.StatusOK, output)
}

type getClientByIdOutput struct {
	Client ClientInfo `json:"data"` // The client
}

// This example route demonstrates how to use Goroutines to speed up the process of getting data from the database.
// This uses channels to communicate between Goroutines.
// getClientById gets a client by their ID and the number of vehicles they have.
// @Summary Get client by ID
// @Description Get a client by their ID and the number of vehicles they have.
// @Tags Clients
// @Accept json
// @Produce json
// @Param id path string true "Client ID (name)"
// @Success 200 {object} getClientByIdOutput
// @Failure 400 {object} ErrorsResponse "Bad request - invalid ID"
// @Failure 404 {object} ErrorsResponse "Client not found"
// @Failure 500 {object} ErrorsResponse "Internal Server Error"
// @Router /clients/{id} [get]
func (env *env) getClientById(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, ErrorsResponse{Errors: []string{"invalid ID"}})
		return
	}

	var wg sync.WaitGroup
	var errChan = make(chan error, 2)
	var clientChan = make(chan *database.Client, 1)
	var vehiclesChan = make(chan *[]string, 1)
	var client *database.Client
	var vehicles *[]string

	// Use Goroutines to speed up the process of getting the client and their vehicles.
	// Goroutine 1: Get the client by ID
	wg.Add(1)
	go func() {
		defer wg.Done()
		client, err := env.db.GetClientsByName(id)
		if err != nil {
			errChan <- err
			return
		}
		clientChan <- client

		// Close the client channel after sending the client.
		close(clientChan)
	}()

	// Goroutine 2: Get the vehicles by client ID
	wg.Add(1)
	go func() {
		defer wg.Done()
		vehicles, err := env.db.GetVehiclesByClient(id)
		if err != nil {
			errChan <- err
			return
		}
		vehiclesChan <- vehicles

		// Close the vehicles channel after sending the vehicles.
		close(vehiclesChan)
	}()

	// Wait for all Goroutines to finish
	wg.Wait()

	// Close the error channel. After closing the channel, you can no longer send values to it,
	// but you can still receive values from it.
	close(errChan)

	// Loop over the error channel to check for any errors that occurred during the Goroutines.
	// Range over a channel will block until the channel is closed.
	for err := range errChan {
		fmt.Println(err)
	}

	// Load the client and vehicles from the channels.
	client = <-clientChan
	vehicles = <-vehiclesChan

	// Send the response.
	var output getClientByIdOutput

	output.Client = ClientInfo{
		Name:         client.Name,
		ContactName:  client.ContactName,
		ContactEmail: client.ContactEmail,
		NumVehicles:  len(*vehicles),
	}

	c.JSON(http.StatusOK, output)
}

type getClientVehiclesOutput struct {
	Vehicle   []BasicVehicle `json:"data"`      // The vehicle
	next_page string         `json:"next_page"` // The URL for the next page of vehicles
}

// getClientVehicles gets all vehicles for a client by their ID.
// @Summary Get client vehicles
// @Description Gets all vehicles for a client by their ID
// @Tags Clients
// @Accept json
// @Produce json
// @Param id path string true "Client ID (name)"
// @Success 200 {object} getClientVehiclesOutput
// @Failure 400 {object} ErrorsResponse "Bad request - invalid ID"
// @Failure 404 {object} ErrorsResponse "Vehicles not found"
// @Failure 500 {object} ErrorsResponse "Internal Server Error"
// @Router /clients/{id}/vehicles [get]
func (env *env) getClientVehicles(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, ErrorsResponse{Errors: []string{"id is required"}})
		return
	}

	vins, err := env.db.GetVehiclesByClient(id)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusNotFound, ErrorsResponse{Errors: []string{"could not retrieve vehicles"}})
		return
	}

	var wg sync.WaitGroup
	var errChan = make(chan error, len(*vins))
	var vehicleChan = make(chan BasicVehicle, len(*vins))

	// Use Goroutines to speed up the process of getting the vehicle information.
	for _, vin := range *vins {
		wg.Add(1)
		mVin := vin // When using Goroutines in a for loop, you need to create a new variable to avoid the loop variable being overwritten.

		// Run the Goroutine
		go func() {
			defer wg.Done()

			vehicle, err := env.db.GetVehicleByVin(mVin)
			if err != nil {
				errChan <- err
				return
			}

			vehicleChan <- BasicVehicle{
				Vin:           vehicle.Vin,
				Mileage:       vehicle.Mileage,
				LargestWeight: 0, // This is a placeholder for the largest weight. You can add this functionality if needed.
			}
		}()
	}

	// Wait for all Goroutines to finish
	go func() {
		wg.Wait()
		close(vehicleChan)
	}()

	// Close the error channel. After closing the channel, you can no longer send values to it,
	// but you can still receive values from it.
	close(errChan)

	// Loop over the error channel to check for any errors that occurred during the Goroutines.
	// Range over a channel will block until the channel is closed.
	for err := range errChan {
		fmt.Println(err)
	}

	// Load the vehicle information from the channel and append it to the allVehicles slice.
	var output getClientVehiclesOutput
	output.next_page = "something" // This is a placeholder for the next page URL. This is usally returned in a database response.
	for vehicle := range vehicleChan {
		output.Vehicle = append(output.Vehicle, vehicle)
	}

	// sort the vehicles by vin
	sort.Slice(output.Vehicle, func(i, j int) bool {
		return output.Vehicle[i].Vin < output.Vehicle[j].Vin
	})

	c.JSON(http.StatusOK, output)
}
