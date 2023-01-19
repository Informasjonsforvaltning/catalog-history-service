package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/Informasjonsforvaltning/catalog-history-service/model"
	"github.com/Informasjonsforvaltning/catalog-history-service/service"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func UpdateJsonPatch(c *gin.Context) {
	original := &model.Update{
		Person: model.Person{
			ID:    "123",
			Email: "emaill",
			Name:  "name",
		},
		DateTime: time.Now(),
		Operations: []model.JsonPatchOperation{
			{
				Op:    "replace",
				Path:  "/name",
				Value: "Jane",
			},
			{
				Op:   "remove",
				Path: "/height",
			},
		},
	}

	originalBytes, err := json.Marshal(original)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	/* 	request, err := ioutil.ReadAll(c.Request.Body)
	   	if err != nil {
	   		fmt.Println(err)
	   		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	   		return
	   	}

	   	patchedJson, err := jsonpatch.MergePatch(originalBytes, request)
	   	if err != nil {
	   		fmt.Println(err)
	   		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	   		return
	   	} */

	c.Data(http.StatusOK, "application/json", originalBytes)
}

func CreateDataSourceHandler() func(c *gin.Context) {
	service := service.InitService()
	return func(c *gin.Context) {
		logrus.Infof("Creating data source")
		bytes, err := c.GetRawData()

		if err != nil {
			logrus.Errorf("Unable to get bytes from request.")

			c.JSON(http.StatusBadRequest, err.Error())
		} else {
			err := service.StoreUpdate(c.Request.Context(), bytes)
			if err == nil {
				c.JSON(http.StatusOK, gin.H{"status": "ok"})
			} else {
				c.JSON(http.StatusInternalServerError, err.Error())
			}
		}
	}
}
