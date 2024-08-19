/*
* @file main.go
* @author Byron Ojua-Nice
* @version 1.0
*
* @section DESCRIPTION
*
* This file contains the main function and the handlers for the API endpoints.
 */

package main

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/byron-ojua/starter-project/database"
	"github.com/gin-gonic/gin"

	_ "github.com/byron-ojua/starter-project/docs"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// ClientWithVehicles is a struct that represents a client and the number of vehicles they have.
type ClientWithVehicles struct {
	Name         string `json:"name"`
	ContactName  string `json:"contact_name"`
	ContactEmail string `json:"contact_email"`
	NumVehicles  int    `json:"number_of_vehicles"`
}

// VehicleInfo is a struct that represents a vehicle and its owner's information.
type VehicleInfo struct {
	Vin          string `json:"vin"`
	ClientName   string `json:"client_name"`
	ContactName  string `json:"contact_name"`
	ContactEmail string `json:"contact_email"`
	Mileage      int    `json:"mileage"`
	Weights      []int  `json:"weights"`
}

// ClientVehicle is a struct that represents a vehicle and its basic information.
type ClientVehicle struct {
	Vin           string `json:"vin"`
	Mileage       int    `json:"mileage"`
	LargestWeight int    `json:"largest_weight"`
}

// ClientVehicles is a struct that represents a client and their vehicles.
type ClientVehicles struct {
	Name     string          `json:"name"`
	Vehicles []ClientVehicle `json:"vehicles"`
}

// @title Simple API
// @version 1
// @description This is a simple API that retrieves information about clients and their vehicles.

// @contact.name Byron Ojua-Nice
// @contact.url https://github.com/byron-ojua
// @contact.email byron.n@air-weigh.com

// @securityDefinitions.apikey bearerToken
// @in header
// @name Authorization

// @host localhost:8080
// @BasePath /
func main() {
	router := gin.Default()
	router.Use(corsMiddleware())

	// Swagger documentation
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	router.GET("/clients", getAllClients)
	router.GET("/clients/:id", getClientByID)
	router.GET("/clients/:id/vehicles", getClientVehicles)
	router.GET("/vehicles/:id", getVehicalByID)

	router.Run("localhost:8080")
}

// corsMiddleware is a middleware function that adds the necessary headers to allow CORS requests.
func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

// getAllClients responds with the list of all clients as JSON.
// @Summary Get all clients
// @Description Get all clients and the number of vehicles they have
// @Tags clients
// @Success 200 {array} ClientWithVehicles
// @Router /clients [get]
func getAllClients(c *gin.Context) {
	var wg sync.WaitGroup
	var mu sync.Mutex

	db, err_db := database.New()
	var clients *[]database.Client
	var err_client error
	var vehicles_by_client = make(map[string]int)
	var all_clients []ClientWithVehicles

	if err_db != nil {
		fmt.Println(err_db)
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "error retrieving clients"})
		return
	}

	clients, err_client = db.GetAllClients()

	if err_client != nil {
		fmt.Println(err_client)
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "clients not found"})
		return
	}

	wg.Add(len(*clients))

	// Use Goroutines to speed up the process of getting the number of vehicles for each client.
	for i := 0; i < len(*clients); i++ {
		var client_name string = (*clients)[i].Name
		go func(client_name string) {
			defer wg.Done()
			temp_vehicles, err_vehicle := db.GetVehiclesByClient(client_name)

			if err_vehicle != nil {
				fmt.Println(err_vehicle)
				return
			}

			// Maps are not thread-safe, so we need to use a mutex to prevent
			mu.Lock()
			vehicles_by_client[client_name] = len(*temp_vehicles)
			mu.Unlock()
		}(client_name)
	}

	wg.Wait()

	for i := 0; i < len(*clients); i++ {
		var client_name string = (*clients)[i].Name
		all_clients = append(all_clients, ClientWithVehicles{
			Name:         client_name,
			ContactName:  (*clients)[i].ContactName,
			ContactEmail: (*clients)[i].ContactEmail,
			NumVehicles:  vehicles_by_client[client_name],
		})
	}

	c.IndentedJSON(http.StatusOK, all_clients)
}

// getClientByID locates the client whose ID value matches the id
// parameter sent by the client, then returns that client as a response.
// @Summary Get a client by ID
// @Description Get a client by their ID and the number of vehicles they have
// @Tags clients
// @Param id path string true "Client ID"
// @Success 200 {object} ClientWithVehicles
// @Router /clients/{id} [get]
func getClientByID(c *gin.Context) {
	id := c.Param("id")
	var wg sync.WaitGroup
	var client *database.Client
	var err_client error
	var vehicles = new([]string)
	var err_vehicles error
	db, err_db := database.New()

	if err_db != nil {
		fmt.Println(err_db)
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "error retrieving client"})
		return
	}

	wg.Add(2)

	// Use Goroutines to speed up the process of getting the client and their vehicles.
	go func() {
		defer wg.Done()
		client, err_client = db.GetClientsByName(id)
	}()

	go func() {
		defer wg.Done()
		vehicles, err_vehicles = db.GetVehiclesByClient(id)
	}()

	wg.Wait()

	// Error handling
	if err_client != nil {
		fmt.Println(err_client)
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": err_client.Error()})
		return
	}

	var numVehicles int
	if err_vehicles != nil {
		numVehicles = 0
	} else {
		numVehicles = len(*vehicles)
	}

	var clientInfo = ClientWithVehicles{
		Name:         client.Name,
		ContactName:  client.ContactName,
		ContactEmail: client.ContactEmail,
		NumVehicles:  numVehicles,
	}

	c.IndentedJSON(http.StatusOK, clientInfo)
}

// getVehicalByID locates the vehicle whoses ID value matches the id
// parameter sent by the client, then returns that vehicle as a response.
// @Summary Get a vehicle by ID
// @Description Get a vehicle by its ID and its owner's information
// @Tags vehicles
// @Param id path string true "Vehicle ID"
// @Success 200 {object} VehicleInfo
// @Router /vehicles/{id} [get]
func getClientVehicles(c *gin.Context) {
	id := c.Param("id")
	var wg sync.WaitGroup
	var mu sync.Mutex
	var vehicle_vins *[]string
	var err_vins error
	var vehicle_info = make(map[string]ClientVehicle)

	db, err_db := database.New()

	if err_db != nil {
		fmt.Println(err_db)
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "error retrieving vehicles"})
		return
	}

	vehicle_vins, err_vins = db.GetVehiclesByClient(id)

	if err_vins != nil {
		fmt.Println(err_vins)
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": err_vins.Error()})
		return
	}

	// create default objects for vehicleInfo map
	for i := 0; i < len(*vehicle_vins); i++ {
		vehicle_info[(*vehicle_vins)[i]] = ClientVehicle{
			Vin:           (*vehicle_vins)[i],
			Mileage:       0,
			LargestWeight: 0,
		}
	}

	// Get the vehicle mileage and weight for each vehicle.
	// These calls are made using Goroutines to speed up the process.
	for i := 0; i < len(*vehicle_vins); i++ {
		wg.Add(2)
		var vin string = (*vehicle_vins)[i]
		go func(vin string) {
			defer wg.Done()
			vehicle, err_vehicle := db.GetVehicleByVin(vin)

			if err_vehicle != nil {
				fmt.Println(err_vehicle)
				return
			}

			// Update the vehicleInfo map with the vehicle mileage.
			// Maps are not thread-safe, so we need to use a mutex to prevent
			mu.Lock()
			vInfo := vehicle_info[vin] // You can't update a field in a struct in a map directly, so you need to get the struct first.
			vInfo.Mileage = vehicle.Mileage
			vehicle_info[vin] = vInfo
			mu.Unlock()
		}(vin)

		go func(vin string) {
			defer wg.Done()
			weights, err_weights := db.GetWeightsByVin(vin)

			if err_weights != nil {
				fmt.Println(err_weights)
				return
			}

			var largest_weight int

			for i := 0; i < len(*weights); i++ {
				if int((*weights)[i].Weight) > largest_weight {
					largest_weight = int((*weights)[i].Weight)
				}
			}

			// Update the vehicleInfo map with the largest weight.
			// Maps are not thread-safe, so we need to use a mutex to prevent
			mu.Lock()
			vInfo := vehicle_info[vin] // You can't update a field in a struct in a map directly, so you need to get the struct first.
			vInfo.LargestWeight = largest_weight
			vehicle_info[vin] = vInfo
			mu.Unlock()
		}(vin)
	}

	wg.Wait()

	// sent the response using the vehicleInfo map
	var client_vehicles ClientVehicles
	client_vehicles.Name = id

	for key := range vehicle_info {
		client_vehicles.Vehicles = append(client_vehicles.Vehicles, vehicle_info[key])
	}

	c.IndentedJSON(http.StatusOK, client_vehicles)
}

// getClientByID locates the client whose ID value matches the id
// parameter sent by the client, then returns that client as a response.
// @Summary Get a client by ID
// @Description Get a client by their ID and the number of vehicles they have
// @Tags clients
// @Param id path string true "Client ID"
// @Success 200 {object} ClientWithVehicles
// @Router /clients/{id} [get]
func getVehicalByID(c *gin.Context) {
	id := c.Param("id")
	var wg sync.WaitGroup
	var vehicle *database.Vehicle
	var err_vehicle error
	var weights *[]database.Weight
	var err_weight error
	var client *database.Client
	var err_client error

	db, err_db := database.New()

	if err_db != nil {
		fmt.Println(err_db)
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "error retrieving vehicle"})
		return
	}

	wg.Add(2)

	// Use Goroutines to speed up the process of getting the vehicle, its weights, and its client.
	go func() {
		defer wg.Done()
		vehicle, err_vehicle = db.GetVehicleByVin(id)
	}()

	go func() {
		defer wg.Done()
		weights, err_weight = db.GetWeightsByVin(id)
	}()

	wg.Wait()

	// Error handling
	if err_vehicle != nil {
		fmt.Println(err_vehicle, err_weight)
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": err_vehicle.Error()})
		return
	}

	client, err_client = db.GetClientsByName(vehicle.Client)

	if err_client != nil {
		fmt.Println(err_vehicle, err_weight, err_client)
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "error retrieving vehicle"})
		return
	}

	var int_weights []int
	if *weights != nil {
		for i := 0; i < len(*weights); i++ {
			int_weights = append(int_weights, int((*weights)[i].Weight))
		}
	}

	var vehicle_info = VehicleInfo{
		Vin:          vehicle.Vin,
		ClientName:   client.Name,
		ContactName:  client.ContactName,
		ContactEmail: client.ContactEmail,
		Mileage:      vehicle.Mileage,
		Weights:      int_weights,
	}

	c.IndentedJSON(http.StatusOK, vehicle_info)
}
