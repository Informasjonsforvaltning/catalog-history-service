package main

import (
	"context"
	"fmt"

	"github.com/Informasjonsforvaltning/catalog-history-service/config"
	"github.com/Informasjonsforvaltning/catalog-history-service/model"
	"github.com/Informasjonsforvaltning/catalog-history-service/repository"
)

func main() {
	config.LoggerSetup()

	router := config.SetupRouter()
	router.Run(":8080")
	router.Run(":9091")

	conceptsRepository := repository.InitRepository()

	concept := model.Concept{
		ID:   1,
		Term: "someTerm",
		Def:  "someDef",
	}

	conceptID, err := conceptsRepository.InsertConcept(context.TODO(), concept)
	if err != nil {
		// handle error
	}

	fmt.Println("Inserted concept with ID:", conceptID)
}
