package main

import "github.com/Informasjonsforvaltning/catalog-history-service/config"

func main() {
	config.LoggerSetup()
	router := config.SetupRouter()
	router.Run(":8080")
}
