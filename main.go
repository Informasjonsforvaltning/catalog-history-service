package main

import (
	"context"
	"log"

	"github.com/Informasjonsforvaltning/catalog-history-service/config/connection"
	"github.com/Informasjonsforvaltning/catalog-history-service/repository"
	"github.com/Informasjonsforvaltning/catalog-history-service/service"
	"github.com/gin-gonic/gin"
)

// define a struct to hold the mock JSON document
type Book struct {
	Title  string `json:"title"`
	Author string `json:"author"`
	Year   int    `json:"year"`
}

func main() {
	// create a new gin router
	router := gin.Default()

	// create a new service and repository
	service := service.InitService()
	repository := repository.InitRepository()

	// apply a JSON Patch to a mock JSON document stored in a MongoDB collection
	router.PATCH("/", func(c *gin.Context) {
		if err := service.ApplyJSONPatch(context.Background(), connection.MongoClient()); err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, gin.H{"status": "success"})
	})

	// get all data sources from the repository
	router.GET("/datasources", repository.GetAllDataSources)

	// start the server
	log.Fatal(router.Run())
}
