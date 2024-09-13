package hub

import (
	"chroma-viz/library/templates"
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (hub *DataBase) StartRestAPI(port int) {
	router := gin.Default()

	router.GET("/templates", hub.templatesGET)
	router.POST("/templates", hub.templatesPOST)

	router.GET("/template/:id", hub.templateGET)
	router.POST("/template", hub.templatePOST)

	router.GET("/tempIDs", hub.tempidsGET)

	router.GET("/assets", hub.assetsGET)

	router.GET("/asset/:id", hub.assetGET)

	router.Run("localhost:" + strconv.Itoa(port))
}

func (hub *DataBase) templatesGET(c *gin.Context) {
	temps, err := hub.GetTemplates()
	if err != nil {
		Logger("Error get templates: %s", err)
		c.Status(http.StatusInternalServerError)
		return
	}

	var tempJSON struct {
		NumTemplates int
		Templates    []*templates.Template
	}

	tempJSON.NumTemplates = len(temps)
	tempJSON.Templates = temps

	c.IndentedJSON(http.StatusOK, tempJSON)
}

func (hub *DataBase) templatesPOST(c *gin.Context) {
	jsonData, err := io.ReadAll(c.Request.Body)
	if err != nil {
		Logger("Error put templates: %s", err)
		c.Status(http.StatusInternalServerError)
		return
	}

	var temps []*templates.Template
	err = json.Unmarshal(jsonData, temps)
	if err != nil {
		Logger("Error put templates: %s", err)
		c.Status(http.StatusInternalServerError)
		return
	}

	for _, temp := range temps {
		if temp == nil {
			continue
		}

		err = hub.ImportTemplate(*temp)
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

	c.IndentedJSON(http.StatusOK, temp)
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

	err = hub.ImportTemplate(temp)
	if err != nil {
		Logger("Error put template: %s", err)
	}

	c.Status(http.StatusOK)
}

func (hub *DataBase) tempidsGET(c *gin.Context) {
	ids, err := hub.TempIDs()
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	c.IndentedJSON(http.StatusOK, ids)
}

func (hub *DataBase) assetsGET(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, hub.Assets)
}

func (hub *DataBase) assetGET(c *gin.Context) {
	_, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}
}
