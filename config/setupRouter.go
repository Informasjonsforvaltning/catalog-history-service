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
	e.POST(env.PathValues.Resource, security.RequireWriteAuth(), handlers.StoreUpdate())
	e.GET(env.PathValues.Resource, security.RequireReadAuth(), handlers.GetUpdates())
	e.GET(env.PathValues.ResourceUpdate, security.RequireReadAuth(), handlers.GetUpdate())
}

func SetupRouter() *gin.Engine {
	router := gin.New()
	router.Use(gin.Recovery())

	InitializeRoutes(router)
	return router
}
