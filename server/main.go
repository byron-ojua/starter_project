package main

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/byron-ojua/starter-project/database"
	"github.com/gin-gonic/gin"
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

func main() {
	router := gin.Default()
	router.Use(corsMiddleware())
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
func getAllClients(c *gin.Context) {
	var wg sync.WaitGroup
	var mu sync.Mutex

	db, err_db := database.New()
	var clients *[]database.Client
	var err_client error
	var vehicles_by_client = make(map[string]int)
	var all_clients []ClientWithVehicles

	if err_db != nil {
		fmt.Println("ERROR!!")
		return
	}

	clients, err_client = db.GetAllClients()

	if err_client != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "clients not found"})
	}

	wg.Add(len(*clients))

	for i := 0; i < len(*clients); i++ {
		var client_name string = (*clients)[i].Name
		go func(clientName string) {
			defer wg.Done()
			tVehicles, _ := db.GetVehiclesByClient(clientName)

			mu.Lock()
			vehicles_by_client[clientName] = len(*tVehicles)
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

	if err_client != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "clients not found"})
	}

	c.IndentedJSON(http.StatusOK, all_clients)
}

// getClientByID locates the client whose ID value matches the id
// parameter sent by the client, then returns that client as a response.
func getClientByID(c *gin.Context) {
	id := c.Param("id")
	var wg sync.WaitGroup
	wg.Add(2)

	db, err := database.New()

	if err != nil {
		fmt.Println("ERROR!!")
		return
	}

	var client *database.Client
	var err_client error
	var vehicles = new([]string)
	var err_vehicles error

	go func() {
		defer wg.Done()
		client, err_client = db.GetClientsByName(id)
	}()

	go func() {
		defer wg.Done()
		vehicles, err_vehicles = db.GetVehiclesByClient(id)
	}()

	wg.Wait()

	if err_client != nil || err_vehicles != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "error retrieving client"})
	}

	var clientInfo = ClientWithVehicles{
		Name:         client.Name,
		ContactName:  client.ContactName,
		ContactEmail: client.ContactEmail,
		NumVehicles:  len(*vehicles),
	}

	c.IndentedJSON(http.StatusOK, clientInfo)
}

// getVehicalByID locates the vehicle whoses ID value matches the id
// parameter sent by the client, then returns that vehicle as a response.
func getClientVehicles(c *gin.Context) {
	id := c.Param("id")
	var wg sync.WaitGroup
	var mu sync.Mutex

	db, err_db := database.New()
	var vehicle_vins *[]string
	var err_vins error
	var vehicle_info = make(map[string]ClientVehicle)

	if err_db != nil {
		fmt.Println("ERROR!!")
		return
	}

	vehicle_vins, err_vins = db.GetVehiclesByClient(id)

	if err_vins != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "error retrieving vehicles"})
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
			vehicle, _ := db.GetVehicleByVin(vin)

			mu.Lock()
			vInfo := vehicle_info[vin]
			vInfo.Mileage = vehicle.Mileage
			vehicle_info[vin] = vInfo
			mu.Unlock()
		}(vin)

		go func(vin string) {
			defer wg.Done()
			weights, _ := db.GetWeightsByVin(vin)
			var largest_weight int

			for i := 0; i < len(*weights); i++ {
				if int((*weights)[i].Weight) > largest_weight {
					largest_weight = int((*weights)[i].Weight)
				}
			}
			mu.Lock()
			vInfo := vehicle_info[vin]
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
func getVehicalByID(c *gin.Context) {
	id := c.Param("id")
	var wg sync.WaitGroup
	wg.Add(2)

	db, err_db := database.New()

	if err_db != nil {
		fmt.Println("ERROR!!")
		return
	}

	var vehicle *database.Vehicle
	var err_vehicle error
	var weights *[]database.Weight
	var err_weight error
	var client *database.Client
	var err_client error

	go func() {
		defer wg.Done()
		vehicle, err_vehicle = db.GetVehicleByVin(id)
	}()

	go func() {
		defer wg.Done()
		weights, err_weight = db.GetWeightsByVin(id)
	}()

	wg.Wait()

	client, err_client = db.GetClientsByName(vehicle.Client)

	if err_vehicle != nil || err_weight != nil || err_client != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "error retrieving vehicle"})
	}

	var int_weights []int
	for i := 0; i < len(*weights); i++ {
		int_weights = append(int_weights, int((*weights)[i].Weight))
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
