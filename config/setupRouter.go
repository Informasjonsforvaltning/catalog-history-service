package config

import (
	"github.com/Informasjonsforvaltning/catalog-history-service/config/security"
	"github.com/Informasjonsforvaltning/catalog-history-service/logging"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"time"

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
	// Validation is performed in the service/repository layer
	e.POST(env.PathValues.Resource, security.RequireWriteAuth(), handlers.StoreUpdate())
	e.GET(env.PathValues.Resource, security.RequireReadAuth(), handlers.GetUpdates())
	e.GET(env.PathValues.ResourceUpdate, security.RequireReadAuth(), handlers.GetUpdate())
	e.GET(env.PathValues.ConceptUpdates, security.RequireReadAuth(), handlers.GetConceptUpdates())
}

func SetupRouter() *gin.Engine {
	router := gin.New()
	router.RedirectFixedPath = false
	router.RedirectTrailingSlash = false
	router.RemoveExtraSlash = false
	router.Use(gin.Recovery())
	router.Use(cors.New(cors.Config{
		AllowOrigins:     env.CorsOriginPatterns(),
		AllowMethods:     []string{"OPTIONS", "GET", "POST"},
		AllowHeaders:     []string{"*"},
		AllowWildcard:    true,
		AllowAllOrigins:  false,
		AllowCredentials: false,
		AllowFiles:       false,
		MaxAge:           1 * time.Hour,
	}))
	InitializeRoutes(router)
	return router
}
