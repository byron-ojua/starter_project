package api

import (
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	_ "github.com/byron-ojua/starter-project/internal/api/docs"
	"github.com/byron-ojua/starter-project/internal/database"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

//go:generate swag init --parseDependency --parseInternal

// @title Starter Project
// @version 1.1
// @description This is a simple API that retrieves information about clients and their vehicles.
// @termsOfService TBD
//
// @contact.name Byron Ojua-Nice
// @contact.url http://firstlaunch.dev
// @contact.email byronojua@firstlaunch.dev
//
// @license.name TBD
// @license.url TBD
//
// @host localhost:8080
// @BasePath /
//
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Use the format: `Bearer <your_token>` (Bearer must be added before the token)
//
// @externalDocs.description  OpenAPI
// @externalDocs.url          https://swagger.io/resources/open-api/

const API_VERSION = "1.1"

// error declarations
const (
	//the error returned to the user when bad user token being provided
	Errunauthorized = "bad user token"

	//the error returned to the user when no token is provided
	ErrNoToke = "no authentication token provided"

	//the error returned to the user when a user is not authorized to perform an action
	ErrForbidden = "forbidden"

	//the error returned when no resource ID is provided for a request pertaining to a specific resource
	ErrorNoResourceID = "no resource id provided"

	//the error returned to the user when the requested record in out of date
	ErrCalibrationOutOfDate = "out of date record"

	//error signifying that a feautre is not supported for a specific product
	ErrUnsupportedFeature = "unsupported product feature"

	//err signifying that a vehicle does not have an attached device
	ErrNoVehicleDevice = "no device attached to vehicle"
)

type Api interface {
	RunLocal() error
}

// ErrorsResponse represents the structure for API error responses.
type ErrorsResponse struct {
	Errors []string `json:"errors"` // Array of error messages
}

// TODO Add logger - Zap is a good option
type env struct {
	api       *gin.Engine         `validate:"required"`
	db        *database.Database  `validate:"required"`
	validator *validator.Validate `validate:"required"`
}

func (e env) RunLocal() error {
	return http.ListenAndServe("localhost:8080", e.api)
}

func New() (Api, error) {
	// Create a new validator
	validator := validator.New()

	// Create a new database connection
	db, err := database.New()
	if err != nil {
		// TODO Add log statement here
		return nil, err
	}

	// Create a new env
	e := env{
		db:        db,
		validator: validator,
	}

	// validate the env
	err = validator.Struct(e)
	if err != nil {
		// TODO Add log statement here
		return nil, err
	}

	r := gin.Default()

	//CORS setup
	r.Use(cors.New(cors.Config{
		AllowHeaders:    []string{"Authorization", "Content-Type"},
		AllowAllOrigins: true,
		AllowMethods:    []string{"GET", "POST", "PUT", "PATCH", "DELETE"},
	}))

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	///////////////////////////////////////////////////////////////////////////
	// 							Client Routes								 //
	///////////////////////////////////////////////////////////////////////////
	r.GET("/clients", e.getClients)
	r.GET("/clients/:id", e.getClientById)
	r.GET("/clients/:id/vehicles", e.getClientVehicles)

	///////////////////////////////////////////////////////////////////////////
	// 							Vehicle Routes								 //
	///////////////////////////////////////////////////////////////////////////
	r.GET("/vehicles/:id", e.getVehicalById)

	e.api = r

	return e, nil
}
