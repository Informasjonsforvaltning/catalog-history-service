package config

import (
	"github.com/Informasjonsforvaltning/catalog-history-service/config/security"
	"github.com/Informasjonsforvaltning/catalog-history-service/logging"
	"github.com/gin-gonic/gin"

	"github.com/Informasjonsforvaltning/catalog-history-service/config/env"
	"github.com/Informasjonsforvaltning/catalog-history-service/handlers"
)

func InitializeRoutes(e *gin.Engine) {
	err := e.SetTrustedProxies(nil)
	if err != nil {
		logging.LogAndPrintError(err)
	}
	e.GET(env.PathValues.Ping, handlers.PingHandler())
	e.GET(env.PathValues.Ready, handlers.ReadyHandler())
	e.POST(env.PathValues.Concept, security.RequireWriteAuth(), handlers.PostConceptUpdate())
	e.GET(env.PathValues.Concept, security.RequireReadAuth(), handlers.GetConceptUpdatesHandler())
	e.GET(env.PathValues.ConceptUpdate, security.RequireReadAuth(), handlers.GetConceptUpdateHandler())
}

func SetupRouter() *gin.Engine {
	router := gin.New()
	router.Use(gin.Recovery())

	InitializeRoutes(router)
	return router
}
