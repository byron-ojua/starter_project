package main

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/byron-ojua/starter-project/database"
	"github.com/gin-gonic/gin"
)

type ClientWithVehicles struct {
	Name         string `json:"name"`
	ContactName  string `json:"contact_name"`
	ContactEmail string `json:"contact_email"`
	NumVehicles  int    `json:"number_of_vehicles"`
}

type VehicleInfo struct {
	Vin          string `json:"vin"`
	ClientName   string `json:"client_name"`
	ContactName  string `json:"contact_name"`
	ContactEmail string `json:"contact_email"`
	Mileage      int    `json:"mileage"`
	Weights      []int  `json:"weights"`
}

type ClientVehicle struct {
	Vin           string `json:"vin"`
	Mileage       int    `json:"mileage"`
	LargestWeight int    `json:"largest_weight"`
}

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

// getAllClients responds with the list of all albums as JSON.
func getAllClients(c *gin.Context) {
	var wg sync.WaitGroup
	var mu sync.Mutex

	db, errDb := database.New()
	var clients *[]database.Client
	var errCli error
	var vehiclesByClient = make(map[string]int)
	var allClients []ClientWithVehicles

	if errDb != nil {
		fmt.Println("ERROR!!")
		return
	}

	clients, errCli = db.GetAllClients()

	wg.Add(len(*clients))

	for i := 0; i < len(*clients); i++ {
		var clientName string = (*clients)[i].Name
		go func(clientName string) {
			defer wg.Done()
			tVehicles, _ := db.GetVehiclesByClient(clientName)

			mu.Lock()
			vehiclesByClient[clientName] = len(*tVehicles)
			mu.Unlock()
		}(clientName)
	}

	wg.Wait()

	for i := 0; i < len(*clients); i++ {
		var clientName string = (*clients)[i].Name
		allClients = append(allClients, ClientWithVehicles{
			Name:         clientName,
			ContactName:  (*clients)[i].ContactName,
			ContactEmail: (*clients)[i].ContactEmail,
			NumVehicles:  vehiclesByClient[clientName],
		})
	}

	if errCli != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "clients not found"})
	}

	c.IndentedJSON(http.StatusOK, allClients)
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
	var errCli error
	var vehicles = new([]string)
	var errVeh error

	go func() {
		defer wg.Done()
		client, errCli = db.GetClientsByName(id)
	}()

	go func() {
		defer wg.Done()
		vehicles, errVeh = db.GetVehiclesByClient(id)
	}()

	wg.Wait()

	if errCli != nil || errVeh != nil {
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

	db, err := database.New()
	var vehicleVins *[]string
	var errVins error
	// var vehicles []database.Vehicle
	// Create a map whose key is the vin and value is a sli
	var vehicleInfo = make(map[string]ClientVehicle)
	// var vehicleWeights = make(map[string]int)

	if err != nil {
		fmt.Println("ERROR!!")
		return
	}

	vehicleVins, errVins = db.GetVehiclesByClient(id)

	if errVins != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "error retrieving vehicles"})
	}

	// create default objects for vehicleInfo map
	for i := 0; i < len(*vehicleVins); i++ {
		vehicleInfo[(*vehicleVins)[i]] = ClientVehicle{
			Vin:           (*vehicleVins)[i],
			Mileage:       0,
			LargestWeight: 0,
		}
	}

	// Get the vehicle mileage and weight for each vehicle.
	// These calls are made concurrently to speed up the process.
	// wg.Add(len(*vehicleVins))

	for i := 0; i < len(*vehicleVins); i++ {
		wg.Add(2)
		var vin string = (*vehicleVins)[i]
		go func(vin string) {
			defer wg.Done()
			vehicle, _ := db.GetVehicleByVin(vin)

			mu.Lock()
			vInfo := vehicleInfo[vin]
			vInfo.Mileage = vehicle.Mileage
			vehicleInfo[vin] = vInfo
			mu.Unlock()
		}(vin)

		go func(vin string) {
			defer wg.Done()
			weights, _ := db.GetWeightsByVin(vin)
			var largestWeight int

			for i := 0; i < len(*weights); i++ {
				if int((*weights)[i].Weight) > largestWeight {
					largestWeight = int((*weights)[i].Weight)
				}
			}
			mu.Lock()
			vInfo := vehicleInfo[vin]
			vInfo.LargestWeight = largestWeight
			vehicleInfo[vin] = vInfo
			mu.Unlock()
		}(vin)
	}

	wg.Wait()

	// sent the response using the vehicleInfo map
	var clientVehicles ClientVehicles
	clientVehicles.Name = id

	for key := range vehicleInfo {
		clientVehicles.Vehicles = append(clientVehicles.Vehicles, vehicleInfo[key])
	}

	c.IndentedJSON(http.StatusOK, clientVehicles)

}

// getClientByID locates the client whose ID value matches the id
// parameter sent by the client, then returns that client as a response.
func getVehicalByID(c *gin.Context) {
	id := c.Param("id")
	var wg sync.WaitGroup
	wg.Add(2)

	db, errDb := database.New()

	if errDb != nil {
		fmt.Println("ERROR!!")
		return
	}

	var vehicle *database.Vehicle
	var errVeh error
	var weights *[]database.Weight
	var errWei error
	var client *database.Client
	var errCli error

	go func() {
		defer wg.Done()
		vehicle, errVeh = db.GetVehicleByVin(id)
	}()

	go func() {
		defer wg.Done()
		weights, errWei = db.GetWeightsByVin(id)
	}()

	wg.Wait()

	client, errCli = db.GetClientsByName(vehicle.Client)

	if errVeh != nil || errWei != nil || errCli != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "error retrieving vehicle"})
	}

	var intWeights []int
	for i := 0; i < len(*weights); i++ {
		intWeights = append(intWeights, int((*weights)[i].Weight))
	}

	var vehicleInfo = VehicleInfo{
		Vin:          vehicle.Vin,
		ClientName:   client.Name,
		ContactName:  client.ContactName,
		ContactEmail: client.ContactEmail,
		Mileage:      vehicle.Mileage,
		Weights:      intWeights,
	}

	c.IndentedJSON(http.StatusOK, vehicleInfo)
}
