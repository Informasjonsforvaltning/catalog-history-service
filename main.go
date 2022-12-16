package main

import "github.com/Informasjonsforvaltning/catalog-history-service/config"

func main() {
	router := config.SetupRouter()
	router.Run(":8080")
}