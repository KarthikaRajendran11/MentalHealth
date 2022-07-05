package main

import (
	"context"
	"os"

	"github.com/gin-gonic/gin"
	server "github.com/mentalhealthco/internal"
	"github.com/mentalhealthco/internal/db/postgres"

	_ "github.com/doug-martin/goqu/v9/dialect/postgres" // dialect registers itself on init
	_ "github.com/lib/pq"
)

func main() {

	client, err := postgres.NewClient(context.Background(), os.Getenv("CONXSTR"))
	if err != nil {
		panic(err)
	}

	service := server.NewService(client)

	gin.SetMode(gin.ReleaseMode)
	ginRouter := gin.New()
	// Use gin.Recovery() to any panic and return a 500 instead of crashing
	// TODO: ZENREACH-23015 - use custom recovery from gin when released
	ginRouter.Use(gin.Recovery())

	service.RegisterRoutes(ginRouter)

	ginRouter.Run()

}
