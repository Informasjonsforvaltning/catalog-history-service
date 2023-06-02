package handlers

import (
	"net/http"
	"strconv"

	"github.com/Informasjonsforvaltning/catalog-history-service/logging"
	"github.com/Informasjonsforvaltning/catalog-history-service/service"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func GetConceptUpdatesHandler() func(c *gin.Context) {
	updateService := service.InitUpdateService()
	return func(c *gin.Context) {
		catalogId := c.Param("catalogId")
		conceptId := c.Param("conceptId")

		page, err := strconv.Atoi(c.Query("page"))
		if err != nil {
			page = 1
		}

		size, err := strconv.Atoi(c.Query("size"))
		if err != nil {
			size = 10
		}

		concepts, status := updateService.GetConceptUpdates(c.Request.Context(), catalogId, conceptId, page, size, c.Query("sort_by"), c.Query("sort_order"))
		if status == http.StatusOK {
			c.JSON(status, concepts)
		} else {
			c.Status(status)
		}
	}
}

func GetConceptUpdateHandler() func(c *gin.Context) {
	updateService := service.InitUpdateService()
	return func(c *gin.Context) {
		catalogId := c.Param("catalogId")
		conceptId := c.Param("conceptId")
		updateId := c.Param("updateId")
		logrus.Infof("Get update %s for concept %s", updateId, conceptId)

		concept, status := updateService.GetConceptUpdate(c.Request.Context(), catalogId, conceptId, updateId)
		if status == http.StatusOK {
			c.JSON(status, concept)
		} else {
			c.Status(status)
		}
	}
}

func PostConceptUpdate() func(c *gin.Context) {
	updateService := service.InitUpdateService()
	return func(c *gin.Context) {
		catalogId := c.Param("catalogId")
		conceptId := c.Param("conceptId")
		logrus.Infof("Update for concept %s received.", conceptId)
		bytes, err := c.GetRawData()

		if err != nil {
			logrus.Errorf("Unable to get bytes from request.")
			logging.LogAndPrintError(err)

			c.JSON(http.StatusBadRequest, err.Error())
		} else {
			newId, err := updateService.StoreConceptUpdate(c.Request.Context(), bytes, catalogId, conceptId)
			if err == nil {
				c.Writer.Header().Add("Location", "/"+catalogId+"/concepts/"+conceptId+"/updates/"+*newId)
				c.JSON(http.StatusCreated, nil)
			} else {
				c.JSON(http.StatusInternalServerError, err.Error())
			}
		}
	}
}
