// TODO: settings.yml to store env variables and secrets
// Settings library to handle the same
package main

import (
	"context"
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/mentalhealthco/common/s3"
	server "github.com/mentalhealthco/internal"
	"github.com/mentalhealthco/internal/db/postgres"
	"github.com/pkg/errors"

	_ "github.com/doug-martin/goqu/v9/dialect/postgres" // dialect registers itself on init
	_ "github.com/lib/pq"
)

func main() {

	gin.SetMode(gin.ReleaseMode)
	ginRouter := gin.New()
	// Use gin.Recovery() to any panic and return a 500 instead of crashing
	// TODO: ZENREACH-23015 - use custom recovery from gin when released
	ginRouter.Use(gin.Recovery())
	// TODO : pass gin context
	// TODO : Add metrics to track number of requests per second
	connectionStr, ok := os.LookupEnv("CONXSTR")
	if !ok {
		fmt.Fprintln(os.Stdout, errors.New("failed to get connection string").Error())
	}
	client, err := postgres.NewClient(context.Background(), connectionStr)
	if err != nil {
		fmt.Fprintln(os.Stdout, errors.Wrap(err, "failed to create postgres client").Error())
	}
	uploader, err := s3.NewS3Uploader()
	if err != nil {
		fmt.Fprintln(os.Stdout, errors.Wrap(err, "failed to create s3 uploader client").Error())
	}
	service := server.NewService(client, uploader)
	service.RegisterRoutes(ginRouter)
	err = ginRouter.Run()
	if err != nil {
		fmt.Fprintln(os.Stdout, errors.Wrap(err, "failed to start service").Error())
	}
}
