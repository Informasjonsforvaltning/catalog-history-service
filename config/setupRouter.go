package config

import (
	"github.com/gin-gonic/gin"

	"github.com/Informasjonsforvaltning/catalog-history-service/config/env"
	"github.com/Informasjonsforvaltning/catalog-history-service/handlers"
)

func InitializeRoutes(e *gin.Engine) {
	e.SetTrustedProxies(nil)
	e.POST(env.PathValues.Concept, handlers.PostConceptUpdate())
	e.GET(env.PathValues.Concept, handlers.GetConceptUpdateHandler())
	e.GET(env.PathValues.ConceptUpdate, handlers.GetConceptUpdatesHandler())
}

func SetupRouter() *gin.Engine {
	router := gin.New()
	router.Use(gin.Recovery())

	InitializeRoutes(router)
	return router
}
