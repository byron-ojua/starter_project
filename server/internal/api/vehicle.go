package api

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/byron-ojua/starter-project/internal/database"
	"github.com/gin-gonic/gin"
)

// VehicleInfo is a struct that represents a vehicle and its owner's information.
type VehicleInfo struct {
	Vin          string `json:"vin"`
	ClientName   string `json:"client_name"`
	ContactName  string `json:"contact_name"`
	ContactEmail string `json:"contact_email"`
	Mileage      int    `json:"mileage"`
	Weights      []int  `json:"weights"`
}

type getVehicleByIdOutput struct {
	Vehicle VehicleInfo `json:"data"`
}

// getVehicalById gets a vehicle by its ID.
// @Summary Get vehicle by ID
// @Description Get a vehicle by its ID.
// @Tags Vehicles
// @Accept json
// @Produce json
// @Param id path string true "Vehicle ID"
// @Success 200 {object} getVehicleByIdOutput
// @Failure 400 {object} ErrorsResponse "Bad request - Id is required"
// @Failure 404 {object} ErrorsResponse "Vehicle not found"
// @Router /vehicles/{id} [get]
func (env *env) getVehicalById(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, ErrorsResponse{Errors: []string{"id is required"}})
		return
	}

	var wg sync.WaitGroup
	var errChan = make(chan error, 2)
	var vehicleChan = make(chan *database.Vehicle, 1)
	var weightsChan = make(chan *[]database.Weight, 1)
	var vehicle *database.Vehicle
	var weights *[]database.Weight

	wg.Add(1)
	go func() {
		defer wg.Done()
		vehicle, err := env.db.GetVehicleByVin(id)
		if err != nil {
			errChan <- err
			return
		}

		vehicleChan <- vehicle
		close(vehicleChan)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		weights, err := env.db.GetWeightsByVin(id)
		if err != nil {
			errChan <- err
			return
		}

		weightsChan <- weights
		close(weightsChan)
	}()

	wg.Wait()

	close(errChan)
	for err := range errChan {
		fmt.Println(err)
	}

	vehicle = <-vehicleChan
	weights = <-weightsChan

	client, err := env.db.GetClientsByName(vehicle.Client)
	if err != nil {
		c.JSON(http.StatusNotFound, ErrorsResponse{Errors: []string{"error retrieving vehicle"}})
		return
	}

	var intWeights []int
	if *weights != nil {
		for i := 0; i < len(*weights); i++ {
			intWeights = append(intWeights, int((*weights)[i].Weight))
		}
	}

	vehicleInfo := VehicleInfo{
		Vin:          vehicle.Vin,
		ClientName:   client.Name,
		ContactName:  client.ContactName,
		ContactEmail: client.ContactEmail,
		Mileage:      vehicle.Mileage,
		Weights:      intWeights,
	}

	fmt.Println(vehicleInfo)

	c.JSON(http.StatusOK, getVehicleByIdOutput{Vehicle: vehicleInfo})
}
