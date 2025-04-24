package hub

import (
	"chroma-viz/library/templates"
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (hub *DataBase) Router() *gin.Engine {
	router := gin.Default()

	router.GET("/templates", hub.templatesGET)
	router.POST("/templates", hub.templatesPOST)

	router.GET("/template/list", hub.tempidsGET)
	router.GET("/template/:id", hub.templateGET)
	router.POST("/template", hub.templatePOST)

	router.GET("/assets", hub.assetsGET)
	router.POST("/assets", hub.assetsPOST)

	router.GET("/asset/:id", hub.assetGET)
	router.POST("/asset", hub.assetPOST)

	router.POST("/clean", hub.cleanPOST)
	router.POST("/generate", hub.generatePOST)

	return router
}

func (hub *DataBase) templatesGET(c *gin.Context) {
	temps, err := hub.GetTemplates()
	if err != nil {
		Logger("Error get templates: %s", err)
		c.Status(http.StatusInternalServerError)
		return
	}

	archive := Templates{
		NumTemplates: len(temps),
		Templates:    temps,
	}

	c.JSON(http.StatusOK, archive)
}

func (hub *DataBase) templatesPOST(c *gin.Context) {
	jsonData, err := io.ReadAll(c.Request.Body)
	if err != nil {
		Logger("Error put templates: %s", err)
		c.Status(http.StatusInternalServerError)
		return
	}

	var archive Templates
	err = json.Unmarshal(jsonData, &archive)
	if err != nil {
		Logger("Error put templates: %s", err)
		c.Status(http.StatusInternalServerError)
		return
	}

	for _, temp := range archive.Templates {
		if temp == nil {
			continue
		}

		err = hub.ImportTemplate(temp)
		if err != nil {
			Logger("Error put templates: %s", err)
		}
	}

	c.Status(http.StatusOK)
}

func (hub *DataBase) templateGET(c *gin.Context) {
	tempID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		Logger("Error get template: %s", err)
		c.Status(http.StatusInternalServerError)
		return
	}

	temp, err := hub.GetTemplate(tempID)
	if err != nil {
		Logger("Error get template: %s", err)
		c.Status(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, temp)
}

func (hub *DataBase) templatePOST(c *gin.Context) {
	jsonData, err := io.ReadAll(c.Request.Body)
	if err != nil {
		Logger("Error put template: %s", err)
		c.Status(http.StatusInternalServerError)
		return
	}

	var temp templates.Template
	err = json.Unmarshal(jsonData, &temp)
	if err != nil {
		Logger("Error put template: %s", err)
		c.Status(http.StatusInternalServerError)
		return
	}

	err = hub.ImportTemplate(&temp)
	if err != nil {
		Logger("Error put template: %s", err)
	}

	c.Status(http.StatusOK)
}

func (hub *DataBase) tempidsGET(c *gin.Context) {
	ids, err := hub.TempIDs()
	if err != nil {
		Logger("Error get tempids: %s", err)
		c.Status(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, ids)
}

func (hub *DataBase) assetsGET(c *gin.Context) {
	assets, err := hub.GetAssets()
	if err != nil {
		Logger("Error get assets: %s", err)
		c.Status(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, assets)
}

func (hub *DataBase) assetsPOST(c *gin.Context) {
	dataJSON, err := io.ReadAll(c.Request.Body)
	if err != nil {
		Logger("Error post assets: %s", err)
		c.Status(http.StatusInternalServerError)
		return
	}

	var assets []Asset
	err = json.Unmarshal(dataJSON, &assets)
	if err != nil {
		Logger("Error post assets: %s", err)
		c.Status(http.StatusInternalServerError)
		return
	}

	for _, asset := range assets {
		err := hub.ImportAsset(asset)
		if err != nil {
			Logger("Error post assets: %s", err)
		}
	}

	c.Status(http.StatusOK)
}

func (hub *DataBase) assetGET(c *gin.Context) {
	assetID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		Logger("Error get asset: %s", err)
		c.Status(http.StatusInternalServerError)
		return
	}

	asset, err := hub.GetAsset(assetID)
	if err != nil {
		Logger("Error get asset: %s", err)
		c.Status(http.StatusInternalServerError)
		return
	}

	c.Data(http.StatusOK, "image", asset.Image)
}

func (hub *DataBase) assetPOST(c *gin.Context) {
	dataJSON, err := io.ReadAll(c.Request.Body)
	if err != nil {
		Logger("Error post asset: %s", err)
		c.Status(http.StatusInternalServerError)
		return
	}

	var asset Asset
	err = json.Unmarshal(dataJSON, &asset)
	if err != nil {
		Logger("Error post asset: %s", err)
		c.Status(http.StatusInternalServerError)
		return
	}

	err = hub.ImportAsset(asset)
	if err != nil {
		Logger("Error post asset: %s", err)
		c.Status(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusOK)
}

func (hub *DataBase) cleanPOST(c *gin.Context) {
	hub.CleanDB()
	c.Status(http.StatusOK)
}

func (hub *DataBase) generatePOST(c *gin.Context) {
	numTemp, numGeo := 100, 1000
	for i := 1; i <= numTemp; i++ {
		err := hub.randomTemplate(int64(i), numGeo)
		if err != nil {
			Logger("Error generating hub: %s", err)
		}
	}

	c.Status(http.StatusOK)
}
