package main

import (
	"github.com/byron-ojua/starter-project/internal/api"
)

func main() {
	api, err := api.New()
	if err != nil {
		panic(err)
	}

	err = api.RunLocal()
	if err != nil {
		panic(err)
	}
}

// corsMiddleware is a middleware function that adds the necessary headers to allow CORS requests.
// func corsMiddleware() gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
// 		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE")
// 		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

// 		if c.Request.Method == "OPTIONS" {
// 			c.AbortWithStatus(204)
// 			return
// 		}

// 		c.Next()
// 	}
// }
