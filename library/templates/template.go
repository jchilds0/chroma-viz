package templates

import (
	"chroma-viz/library/props"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/gotk3/gotk3/gtk"
)

/*

   Templates form the basis of pages. Each page corresponds to
   one page which specifies the shape of the page, usually with
   a number of properties that can be edited by the user.

*/

type Template struct {
    Title       string
    TempID      int
    NumGeo      int
    Layer       int
    Geometry    map[int]*props.Property
    AnimateOn   string
    Continue    string
    AnimateOff  string
}

func NewTemplate(title string, id int, layer int, num_geo int) *Template {
    temp := &Template{Title: title, TempID: id, Layer: layer, NumGeo: num_geo}

    temp.Geometry = make(map[int]*props.Property, num_geo)
    return temp
}

func (temp *Template) TemplateToListRow() *gtk.ListBoxRow {
    row1, err := gtk.ListBoxRowNew()
    if err != nil {
        log.Fatalf("Error converting template to row (%s)", err)
    }

    row1.Add(TextToBuffer(temp.Title))
    return row1
}

func (temp *Template) AddGeometry(name string, geo_id, typed int, visible map[string]bool) *props.Property {
    temp.Geometry[geo_id] = props.NewProperty(typed, name, visible, func(){})
    temp.NumGeo++

    return temp.Geometry[geo_id]
}

func TextToBuffer(text string) *gtk.TextView {
    text1, err := gtk.TextViewNew()
    if err != nil {
        log.Fatalf("Error creating text buffer (%s)", err)
    }

    buffer, err := text1.GetBuffer()
    if err != nil {
        log.Fatalf("Error creating text buffer (%s)", err)
    }

    buffer.SetText(text)
    return text1
}

// T -> {'id': num, 'num_geo': num, 'layer': num, 'geometry': [G]} | T, T 
func (temp *Template) Encode() string {
    first := true 
    templates := ""
    for geo_id, prop := range temp.Geometry {
        if first {
            templates = prop.Encode(geo_id)
            first = false 
            continue
        }

        templates = fmt.Sprintf("%s,%s", templates, prop.Encode(geo_id))
    }

    return fmt.Sprintf("{'id': %d, 'num_geo': %d, 'name': '%s', 'layer': %d, " + 
        "'anim_on': '%s', 'anim_cont': '%s', 'anim_off': '%s', 'geometry': [%s]}", 
        temp.TempID, len(temp.Geometry), temp.Title, temp.Layer, 
        temp.AnimateOn, temp.Continue, temp.AnimateOff, templates)
}

func ExportTemplate(temp *Template, filename string) error {
    file, err := os.Create(filename)
    if err != nil {
        return err
    }
    defer file.Close()

    buf, err := json.Marshal(temp)
    if err != nil {
        return err
    }

    _, err = file.Write(buf)
    if err != nil {
        return err
    }

    return nil
}

func (temp *Template) GetTemplateID() int {
    return temp.TempID
}

func (temp *Template) GetLayer() int {
    return temp.Layer 
}

func (temp *Template) GetPropMap() map[int]*props.Property {
    return temp.Geometry
}
