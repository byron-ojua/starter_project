package database

type Client struct {
	Name         string
	ContactName  string
	ContactEmail string
}

type Vehicle struct {
	Vin     string
	Client  string
	Mileage int
}

type Weight struct {
	Vin    string
	Weight float32
}
