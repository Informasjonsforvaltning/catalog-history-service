package config

import (
	"github.com/gin-gonic/gin"

	"github.com/Informasjonsforvaltning/catalog-history-service/config/env"
	"github.com/Informasjonsforvaltning/catalog-history-service/handlers"
)

func InitializeRoutes(e *gin.Engine) {
	e.SetTrustedProxies(nil)
	e.GET(env.PathValues.Ping, handlers.PingHandler())
	e.GET(env.PathValues.Ready, handlers.ReadyHandler())
	e.POST(env.PathValues.Concept, handlers.PostConceptUpdate())
	e.GET(env.PathValues.Concept, handlers.GetConceptUpdatesHandler())
	e.GET(env.PathValues.ConceptUpdate, handlers.GetConceptUpdateHandler())
	e.GET(env.PathValues.ConceptDiff, handlers.GetConceptDiffHandler())
}

func SetupRouter() *gin.Engine {
	router := gin.New()
	router.Use(gin.Recovery())

	InitializeRoutes(router)
	return router
}
