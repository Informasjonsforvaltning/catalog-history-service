package handlers

import (
	"net/http"

	"github.com/Informasjonsforvaltning/catalog-history-service/service"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func GetAllHandler() func(c *gin.Context) {
	service := service.InitService()
	return func(c *gin.Context) {
		logrus.Info("Getting all concepts")

		concepts, status := service.GetConceptUpdates(c.Request.Context(), nil)
		if status == http.StatusOK {
			c.JSON(status, concepts)
		} else {
			c.Status(status)
		}
	}
}

func GetUpdateHandler() func(c *gin.Context) {
	service := service.InitService()
	return func(c *gin.Context) {
		id := c.Param("id")
		logrus.Infof("Get concept update with id: %s", id)

		concept, status := service.GetConceptUpdate(c.Request.Context(), id)
		if status == http.StatusOK {
			c.JSON(status, concept)
		} else {
			c.Status(status)
		}
	}
}
