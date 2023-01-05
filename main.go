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

	begreperRepository := repository.InitRepository()

	begrep := model.Begrep{
		ID:   1,
		Term: "someTerm",
		Def:  "someDef",
	}

	begrepID, err := begreperRepository.InsertBegrep(context.TODO(), begrep)
	if err != nil {
		// handle error
	}

	fmt.Println("Inserted begrep with ID:", begrepID)
}
