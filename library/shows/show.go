package shows

import (
	"chroma-viz/library/templates"
	"encoding/json"
	"log"
	"os"
)

type Show struct {
	NumPages int
	Pages    map[int]*Page
}

func NewShow() *Show {
	show := &Show{}
	show.Pages = make(map[int]*Page)
	return show
}

func (show *Show) SetPage(pageNum int, title string, temp *templates.Template) {
	page := NewPage(pageNum, temp.TempID, temp.Layer, temp.NumGeo, title)
    page.PropMap = temp.GetPropMap()

	if _, ok := show.Pages[pageNum]; ok {
		log.Printf("Page %d already exists", pageNum)
		return
	}

	show.Pages[pageNum] = page
}

func (show *Show) AddPage(title string, temp *templates.Template) *Page {
	show.NumPages++
	show.SetPage(show.NumPages, title, temp)

	return show.Pages[show.NumPages]
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

func (show *Show) ExportShow(filename string) {
	file, err := os.Create(filename)
	if err != nil {
		log.Printf("Error exporting show (%s)", err)
		return
	}
	defer file.Close()

	buf, err := json.Marshal(show)
	if err != nil {
		log.Printf("Error exporting show (%s)", err)
		return
	}

	_, err = file.Write(buf)
	if err != nil {
		log.Printf("Error exporting show (%s)", err)
		return
	}
}
