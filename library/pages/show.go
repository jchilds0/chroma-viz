package pages

import (
	"chroma-viz/library/templates"
	"encoding/json"
	"fmt"
	"os"
)

type Show struct {
	NumPages int
	Pages    map[int]*Page
}

func NewShow() *Show {
	show := &Show{NumPages: 1}
	show.Pages = make(map[int]*Page)
	return show
}

func (show *Show) AddPage(page *Page) {
	page.PageNum = show.NumPages
	show.Pages[page.PageNum] = page
	show.NumPages++
}

func (show *Show) TempToPage(title string, temp *templates.Template) (page *Page, err error) {
	show.NumPages++

	//page = NewPage(show.NumPages, temp.TempID, temp.Layer, temp.NumGeo, title)
	page = NewPageFromTemplate(temp)
	if _, ok := show.Pages[page.PageNum]; ok {
		err = fmt.Errorf("Page %d already exists", page.PageNum)
		return
	}

	show.Pages[page.PageNum] = page
	return
}

func (show *Show) ImportShow(filename string) error {
	buf, err := os.ReadFile(filename)
	if err != nil {
		return err
	}

	err = json.Unmarshal(buf, show)
	if err != nil {
		return err
	}

	return nil
}

func (show *Show) ExportShow(filename string) (err error) {
	file, err := os.Create(filename)
	if err != nil {
		return
	}
	defer file.Close()

	buf, err := json.Marshal(show)
	if err != nil {
		return
	}

	_, err = file.Write(buf)
	return
}
