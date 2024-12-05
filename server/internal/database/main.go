package database

type Database struct {
	clients  map[string]Client
	vehicles map[string]Vehicle
	weight   map[string][]Weight
}

func New() (*Database, error) {
	database := getdata()
	return database, nil
}

func getdata() *Database {
	return &Database{
		clients: map[string]Client{
			"Bobs Burgers": {
				Name:         "Bobs Burgers",
				ContactName:  "Bob Belcher",
				ContactEmail: "bob@bestburgers.com",
			},
			"Dunder Mifflin": {
				Name:         "Dunder Mifflin",
				ContactName:  "Michael Scott",
				ContactEmail: "bestboss@dunermifflin.com",
			},
			"CIA": {
				Name:         "CIA",
				ContactName:  "Stan Smith",
				ContactEmail: "stan@cia.com",
			},
		},
		vehicles: map[string]Vehicle{
			"123456789G": {
				Vin:     "123456789G",
				Client:  "Bobs Burgers",
				Mileage: 100783,
			},
			"123E456789G": {
				Vin:     "123E456789G",
				Client:  "Bobs Burgers",
				Mileage: 107598,
			},
			"23E456789G": {
				Vin:     "23E456789G",
				Client:  "Bobs Burgers",
				Mileage: 178783,
			},
			"23EFU456789G": {
				Vin:     "23EFU456789G",
				Client:  "Dunder Mifflin",
				Mileage: 124783,
			},
			"23EFU4FW56789G": {
				Vin:     "23EFU4FW56789G",
				Client:  "Dunder Mifflin",
				Mileage: 10783,
			},
			"23EFfwU4FW56789G": {
				Vin:     "23EFfwU4FW56789G",
				Client:  "Dunder Mifflin",
				Mileage: 14783,
			},
			"23EFU4FW5fe6789G": {
				Vin:     "23EFU4FW5fe6789G",
				Client:  "Dunder Mifflin",
				Mileage: 1100783,
			},
			"23EFU4FW5678f39G": {
				Vin:     "23EFU4FW5678f39G",
				Client:  "CIA",
				Mileage: 103,
			},
			"23EFU4FW5678ff39G": {
				Vin:     "23EFU4FW5678ff39G",
				Client:  "CIA",
				Mileage: 0,
			},
		},
		weight: map[string][]Weight{
			"123456789G": {
				{
					Vin:    "123456789G",
					Weight: 32.1,
				},
				{
					Vin:    "123456789G",
					Weight: 106,
				},
				{
					Vin:    "123456789G",
					Weight: 5.36,
				},
			},
			"123E456789G": {
				{
					Vin:    "123E456789G",
					Weight: 104,
				},
				{
					Vin:    "123E456789G",
					Weight: 2342,
				},
			},
			"23E456789G": {
				{
					Vin:    "23E456789G",
					Weight: 9182,
				},
				{
					Vin:    "23E456789G",
					Weight: 2346,
				},
				{
					Vin:    "23E456789G",
					Weight: 56856,
				},
			},
			"23EFU456789G": {
				{
					Vin:    "23EFU456789G",
					Weight: 10.236,
				},
				{
					Vin:    "23EFU456789G",
					Weight: 10234.6,
				},
				{
					Vin:    "23EFU456789G",
					Weight: 5347890,
				},
			},
			"23EFU4FW56789G": {
				{
					Vin:    "23EFU4FW56789G",
					Weight: 0.2,
				},
				{
					Vin:    "23EFU4FW56789G",
					Weight: 23467,
				},
				{
					Vin:    "23EFU4FW56789G",
					Weight: 10.6,
				},
				{
					Vin:    "23EFU4FW56789G",
					Weight: 786,
				},
			},
			"23EFfwU4FW56789G": {
				{
					Vin:    "23EFfwU4FW56789G",
					Weight: 14,
				},
				{
					Vin:    "23EFfwU4FW56789G",
					Weight: 1564,
				},
				{
					Vin:    "23EFfwU4FW56789G",
					Weight: 134,
				},
				{
					Vin:    "23EFfwU4FW56789G",
					Weight: 1442,
				},
			},
			"23EFU4FW5fe6789G": {
				{
					Vin:    "23EFU4FW5fe6789G",
					Weight: 10.36,
				},
				{
					Vin:    "23EFU4FW5fe6789G",
					Weight: 16,
				},
			},
			"23EFU4FW5678f39G": {
				{
					Vin:    "23EFU4FW5678f39G",
					Weight: 17,
				},
			},
			"23EFU4FW5678ff39G": {
				{
					Vin:    "23EFU4FW5678ff39G",
					Weight: 10.6,
				},
				{
					Vin:    "23EFU4FW5678ff39G",
					Weight: 11000,
				},
			},
		},
	}
}
