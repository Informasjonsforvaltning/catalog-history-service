package handlers

import (
	"net/http"

	"github.com/Informasjonsforvaltning/catalog-history-service/service"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/ulule/paging"
)

func GetConceptUpdatesHandler(service *service.UpdateService) func(c *gin.Context) {
	return func(c *gin.Context) {
		conceptId := c.Param("conceptId")
		logrus.Info("Getting all updates for concepts with id: %s", conceptId)
		page, size := pagination.GetPageAndSize(c)
		updates, total, status := service.GetConceptUpdates(c.Request.Context(), conceptId, page, size)
		if status == http.StatusOK {
			c.JSON(status, gin.H{
				"updates": updates,
				"total":   total,
			})
		} else {
			c.AbortWithStatus(status)
		}
	}
}

func GetConceptUpdateHandler() func(c *gin.Context) {
	service := service.InitService()
	return func(c *gin.Context) {
		conceptId := c.Param("conceptId")
		updateId := c.Param("updateId")
		logrus.Infof("Get update %s for concept %s", updateId, conceptId)

		concept, status := service.GetConceptUpdate(c.Request.Context(), conceptId, updateId)
		if status == http.StatusOK {
			c.JSON(status, concept)
		} else {
			c.Status(status)
		}
	}
}

// function to get diff between two updates
func GetConceptDiffHandler() func(c *gin.Context) {
	service := service.InitService()
	return func(c *gin.Context) {
		conceptId := c.Param("conceptId")
		updateId := c.Param("updateId")
		logrus.Infof("Get diff for update %s for concept %s", updateId, conceptId)

		diff, status := service.GetConceptUpdateDiff(c.Request.Context(), conceptId, updateId)
		if status == http.StatusOK {
			c.JSON(status, diff)
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
				c.Writer.Header().Add("Location", "/concepts/"+conceptId+"/"+*newId)
				c.JSON(http.StatusCreated, nil)
			} else {
				c.JSON(http.StatusInternalServerError, err.Error())
			}
		}
	}
}
