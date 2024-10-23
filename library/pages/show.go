package pages

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type LocalShow struct {
	NumPages int
	Pages    map[int]*Page
	router   *gin.Engine
}

func NewLocalShow(mediaPort int) *LocalShow {
	show := &LocalShow{NumPages: 1}
	show.Pages = make(map[int]*Page)

	show.router = gin.Default()
	show.router.GET("/page/:id", show.pageGET)
	show.router.GET("/page", show.pagesGET)
	show.router.POST("/page", show.pagePOST)

	go show.router.Run(fmt.Sprintf("localhost:%d", mediaPort))

	return show
}

func (show *LocalShow) AddPage(page *Page) error {
	page.PageNum = show.NumPages
	show.Pages[page.PageNum] = page
	show.NumPages++
	return nil
}

func (show *LocalShow) GetPage(pageNum int) (*Page, bool) {
	page, ok := show.Pages[pageNum]
	return page, ok
}

func (show *LocalShow) GetPages() (map[int]PageData, error) {
	pages := make(map[int]PageData, len(show.Pages))
	for id, p := range show.Pages {
		pages[id] = PageData{
			PageNum: p.PageNum,
			Title:   p.Title,
			TempID:  p.TemplateID,
			Layer:   p.Layer,
		}
	}

	return pages, nil
}

func (show *LocalShow) DeletePage(pageNum int) {
	delete(show.Pages, pageNum)
}

func (show *LocalShow) Clear() {
	clear(show.Pages)
	show.NumPages = 1
}

type PageData struct {
	PageNum int
	TempID  int
	Layer   int
	Title   string
}

func (show *LocalShow) pagesGET(c *gin.Context) {
	pages, _ := show.GetPages()
	c.JSON(http.StatusOK, pages)
}

func (show *LocalShow) pageGET(c *gin.Context) {
	pageNum := c.GetInt("id")
	page, ok := show.GetPage(pageNum)
	if !ok {
		c.Status(http.StatusNotFound)
		return
	}

	c.JSON(http.StatusOK, page)
}

func (show *LocalShow) pagePOST(c *gin.Context) {
	jsonData, err := io.ReadAll(c.Request.Body)
	if err != nil {
		log.Println("Error put page:", err)
		c.Status(http.StatusInternalServerError)
		return
	}

	var page Page
	err = json.Unmarshal(jsonData, &page)
	if err != nil {
		log.Println("Error put page:", err)
		c.Status(http.StatusInternalServerError)
		return
	}

	show.AddPage(&page)
	c.JSON(http.StatusOK, page.PageNum)
}
