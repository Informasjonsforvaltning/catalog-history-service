package main

import (
	"github.com/Informasjonsforvaltning/catalog-history-service/config"
	"github.com/Informasjonsforvaltning/catalog-history-service/handlers"
	"github.com/Informasjonsforvaltning/catalog-history-service/service"
)

func main() {
	config.LoggerSetup()

	router := config.SetupRouter()
	storeUpdate := service.NewUpdateService()
	router.POST("/concepts/:conceptId", handlers.NewUpdateHandler(storeUpdate))
	router.Run(":8080")
	router.Run(":9091")
}
