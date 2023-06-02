package main

import (
	"github.com/Informasjonsforvaltning/catalog-history-service/config"
	"github.com/Informasjonsforvaltning/catalog-history-service/logging"
)

func main() {
	logging.LoggerSetup()

	router := config.SetupRouter()
	err := router.Run(":8080")
	if err != nil {
		logging.LogAndPrintError(err)
	}
}
