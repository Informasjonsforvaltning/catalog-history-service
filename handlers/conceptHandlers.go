package handlers

import (
	"net/http"
	"strconv"

	"github.com/Informasjonsforvaltning/catalog-history-service/service"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func GetConceptUpdatesHandler() func(c *gin.Context) {
	service := service.InitService()
	return func(c *gin.Context) {
		conceptId := c.Param("conceptId")
		logrus.Info("Getting all updates for concepts with id: %s", conceptId)
		concepts, status := service.GetConceptUpdates(c.Request.Context(), conceptId)
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
		logrus.Infof("Get update %s for concept %s", updateId, conceptId)

		// Get the skip and limit values from the query parameters
		skip := c.Query("skip")
		limit := c.Query("limit")

		// Convert the skip and limit values to integers
		skipInt, err := strconv.Atoi(skip)
		if err != nil {
			skipInt = 0
		}

		limitInt, err := strconv.Atoi(limit)
		if err != nil {
			limitInt = 10
		}

		concept, status := service.GetConceptUpdate(c.Request.Context(), conceptId, updateId, skipInt, limitInt)
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
		conceptId := c.Param("conceptId")
		logrus.Infof("Update for concept %s received.", conceptId)
		bytes, err := c.GetRawData()

		if err != nil {
			logrus.Errorf("Unable to get bytes from request.")

			c.JSON(http.StatusBadRequest, err.Error())
		} else {
			newId, err := service.StoreConceptUpdate(c.Request.Context(), bytes, conceptId)
			if err == nil {
				c.Writer.Header().Add("Location", "/concepts/"+conceptId+"/updates/"+*newId)
				c.JSON(http.StatusCreated, nil)
			} else {
				c.JSON(http.StatusInternalServerError, err.Error())
			}
		}
	}
}
