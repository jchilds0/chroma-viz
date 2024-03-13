package templates

import (
	"chroma-viz/library/props"
	"fmt"
	"log"

	"github.com/gotk3/gotk3/gtk"
)

/*

   Templates form the basis of pages. Each page corresponds to
   one page which specifies the shape of the page, usually with
   a number of properties that can be edited by the user.

   Prop is a simple struct used for storing data as we parse the
   templates sent by Chroma Hub on startup. This is then used to
   generate a Page.

*/

type Prop struct {
    Name        string
    Type        int 
    Visible     map[string]bool
    Value       map[string]string
}

func NewProp(name string, typed int) *Prop {
    p := &Prop{Name: name, Type: typed}
    p.Visible = make(map[string]bool)
    p.Value = make(map[string]string)

    return p
}

func (p *Prop) AddAttr(name, value string, visible bool) {
    p.Visible[name] = visible
    p.Value[name] = value
}

type Template struct {
    Title       string
    TempID      int
    NumProps    int
    Layer       int
    Geometry    map[int]*props.Property
}

func NewTemplate(title string, id int, layer int, num_geo int) *Template {
    temp := &Template{Title: title, TempID: id, Layer: layer}

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

func (temp *Template) AddProp(name string, geo_id, typed int) *props.Property {
    temp.Geometry[geo_id] = props.NewProperty(typed, name, nil, func(){})
    temp.NumProps++

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
        // TODO: Pull visible attrs from editor
        prop.Visible = make(map[string]bool)

        if first {
            templates = prop.Encode(geo_id)
            first = false 
            continue
        }

        templates = fmt.Sprintf("%s,%s", templates, prop.Encode(geo_id))
    }

    return fmt.Sprintf("{'id': %d, 'num_geo': %d, 'layer': %d, 'geometry': [%s]}", 
        temp.TempID, len(temp.Geometry), temp.Layer, templates)
}

