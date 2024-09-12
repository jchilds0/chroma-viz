package hub

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (hub *DataBase) getTemplates(c *gin.Context) {
	temps, err := hub.GetTemplates()
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	c.IndentedJSON(http.StatusOK, temps)
}

func (hub *DataBase) getTemplate(c *gin.Context) {
	tempID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	temp, err := hub.GetTemplate(tempID)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	c.IndentedJSON(http.StatusOK, temp)
}

func (hub *DataBase) getTemplateIDs(c *gin.Context) {
	ids, err := hub.TempIDs()
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	c.IndentedJSON(http.StatusOK, ids)
}

func (hub *DataBase) getAssets(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, hub.Assets)
}

func (hub *DataBase) getAsset(c *gin.Context) {
	_, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}
}
