package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/Informasjonsforvaltning/catalog-history-service/logging"
	"github.com/Informasjonsforvaltning/catalog-history-service/service"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func GetUpdates() func(c *gin.Context) {
	updateService := service.InitUpdateService()
	return func(c *gin.Context) {
		catalogId := c.Param("catalogId")
		resourceId := c.Param("resourceId")

		// Parse query parameters - validation happens in service layer
		page, err := strconv.Atoi(c.Query("page"))
		if err != nil {
			page = 0
		}

		size, err := strconv.Atoi(c.Query("size"))
		if err != nil {
			size = 10
		}

		sortBy := c.Query("sort_by")
		sortOrder := c.Query("sort_order")

		updates, status := updateService.GetUpdates(c.Request.Context(), catalogId, resourceId, page, size, sortBy, sortOrder)
		if status == http.StatusOK {
			c.JSON(status, updates)
		} else {
			c.Status(status)
		}
	}
}

func GetConceptUpdates() func(c *gin.Context) {
	updateService := service.InitUpdateService()
	return func(c *gin.Context) {
		catalogId := c.Param("catalogId")

		// Parse query parameters - validation happens in service layer
		page, err := strconv.Atoi(c.Query("page"))
		if err != nil {
			page = 0
		}

		size, err := strconv.Atoi(c.Query("size"))
		if err != nil {
			size = 10
		}

		sortBy := c.Query("sort_by")
		sortOrder := c.Query("sort_order")

		updates, status := updateService.GetConceptUpdates(c.Request.Context(), catalogId, page, size, sortBy, sortOrder)
		if status == http.StatusOK {
			c.JSON(status, updates)
		} else {
			c.Status(status)
		}
	}
}

func GetUpdate() func(c *gin.Context) {
	updateService := service.InitUpdateService()
	return func(c *gin.Context) {
		catalogId := c.Param("catalogId")
		resourceId := c.Param("resourceId")
		updateId := c.Param("updateId")
		// Validation happens in service/repository layer
		logrus.Infof("Get update %s for resource %s", updateId, resourceId)

		update, status := updateService.GetUpdate(c.Request.Context(), catalogId, resourceId, updateId)
		if status == http.StatusOK {
			c.JSON(status, update)
		} else {
			c.Status(status)
		}
	}
}

func StoreUpdate() func(c *gin.Context) {
	updateService := service.InitUpdateService()
	return func(c *gin.Context) {
		catalogId := c.Param("catalogId")
		resourceId := c.Param("resourceId")
		logrus.Infof("Update for resource %s received.", resourceId)
		bytes, err := c.GetRawData()

		if err != nil {
			logrus.Errorf("Unable to get bytes from request.")
			logging.LogAndPrintError(err)

			c.JSON(http.StatusBadRequest, err.Error())
		} else {
			newId, err := updateService.StoreUpdate(c.Request.Context(), bytes, catalogId, resourceId)
			if err == nil {
				c.Writer.Header().Add("Location", fmt.Sprintf("/%s/%s/updates/%s", catalogId, resourceId, newId))
				c.JSON(http.StatusCreated, nil)
			} else {
				c.JSON(http.StatusInternalServerError, err.Error())
			}
		}
	}
}
