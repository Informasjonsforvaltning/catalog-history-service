package handlers

import (
	"net/http"

	"github.com/Informasjonsforvaltning/catalog-history-service/service"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func GetConceptUpdatesHandler() func(c *gin.Context) {
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

func GetConceptUpdateHandler() func(c *gin.Context) {
	service := service.InitService()
	return func(c *gin.Context) {
		conceptId := c.Param("conceptId")
		updateId := c.Param("updateId")
		logrus.Infof("Get concept update with id: %s", conceptId)

		concept, status := service.GetConceptUpdate(c.Request.Context(), conceptId, updateId)
		if status == http.StatusOK {
			c.JSON(status, concept)
		} else {
			c.Status(status)
		}
	}
}

func PostConceptUpdate() func(c *gin.Context) {
	service := service.InitService()
	return func(c *gin.Context) {
		logrus.Infof("Concept update received.")
		bytes, err := c.GetRawData()

		if err != nil {
			logrus.Errorf("Unable to get bytes from request.")

			c.JSON(http.StatusBadRequest, err.Error())
		} else {
			newId, err := service.StoreConceptUpdate(c.Request.Context(), bytes, c.Param("conceptId"))
			if err == nil {
				c.Writer.Header().Add("Location", "/concepts/"+*newId)
				c.JSON(http.StatusCreated, nil)
			} else {
				c.JSON(http.StatusInternalServerError, err.Error())
			}
		}
	}
}
