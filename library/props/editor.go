package props

import (
	"chroma-viz/library/attribute"
	"fmt"
	"log"

	"github.com/gotk3/gotk3/gtk"
)

/*

    Creating PropertyEditor's has a non negligible cost due to 
    the gtk c calls, so to speed up initialisation we use a 
    limited number of PropertyEditor's (enough to show a Page).

    PropertyEditor's consist of gtk ui elements which are used
    to edit Properties. When a user selects a Page in the Show,
    the Page is sent to the editor, which calls UpdateEditor 
    (or UpdateEditorAllProp). This updates the PropertyEditor 
    with the data from the Property.

    The changes made to the PropertyEditor are synced back to 
    the Property using UpdateProp (see props.go).

    The Properties are built up from a collection of Attributes,
    and in a similar way, PropertyEditor's are built up from a 
    collection of AttributeEditors.

*/

type PropertyEditor struct {
    PropType    int
    Box         *gtk.Box
    editor      map[string]attribute.Editor
    visible     map[string]*gtk.CheckButton
}

func NewPropertyEditor(typed int) (propEdit *PropertyEditor, err error) {
    propEdit = &PropertyEditor{PropType: typed}
    width := 1920
    height := 1080 

    propEdit.editor = make(map[string]attribute.Editor, 10)
    propEdit.Box, err = gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0)
    if err != nil {
        return
    }

    propEdit.Box.SetVisible(true)

    switch (typed) {
    case RECT_PROP:
        propEdit.editor["x"] = attribute.NewIntEditor("x", -float64(width), float64(width))
        propEdit.editor["y"] = attribute.NewIntEditor("y", -float64(width), float64(height))
        propEdit.editor["width"] = attribute.NewIntEditor("Width", 0, float64(width))
        propEdit.editor["height"] = attribute.NewIntEditor("Height", 0, float64(height))
        propEdit.editor["rounding"] = attribute.NewIntEditor("Rounding", 0, float64(width))
        propEdit.editor["color"] = attribute.NewColorEditor("Color")

    case TEXT_PROP:
        propEdit.editor["x"] = attribute.NewIntEditor("Center x", 0, float64(width))
        propEdit.editor["y"] = attribute.NewIntEditor("Center y", 0, float64(height))
        propEdit.editor["string"] = attribute.NewStringEditor("Text")
        propEdit.editor["color"] = attribute.NewColorEditor("Color")

    case CIRCLE_PROP:
        propEdit.editor["x"] = attribute.NewIntEditor("Center x", 0, float64(width))
        propEdit.editor["y"] = attribute.NewIntEditor("Center y", 0, float64(height))
        propEdit.editor["inner_radius"] = attribute.NewIntEditor("Inner Radius", 0, float64(width))
        propEdit.editor["outer_radius"] = attribute.NewIntEditor("Outer Radius", 0, float64(width))
        propEdit.editor["start_angle"] = attribute.NewIntEditor("Start Angle", 0, 360)
        propEdit.editor["end_angle"] = attribute.NewIntEditor("End Angle", 0, 360)
        propEdit.editor["color"] = attribute.NewColorEditor("Color")

    case GRAPH_PROP:
        propEdit.editor["x"] = attribute.NewIntEditor("Center x", 0, float64(width))
        propEdit.editor["y"] = attribute.NewIntEditor("Center y", 0, float64(height))
        propEdit.editor["node"] = attribute.NewListEditor("Graph", []string{"x Pos", "y Pos"})
        propEdit.editor["color"] = attribute.NewColorEditor("Color")

    case TICKER_PROP:
        propEdit.editor["x"] = attribute.NewIntEditor("Center x", 0, float64(width))
        propEdit.editor["y"] = attribute.NewIntEditor("Center y", 0, float64(height))
        propEdit.editor["color"] = attribute.NewColorEditor("Color")
        propEdit.editor["text"] = attribute.NewListEditor("Ticker", []string{"Text"})
    
    case CLOCK_PROP:
        propEdit.editor["x"] = attribute.NewIntEditor("Center x", 0, float64(width))
        propEdit.editor["y"] = attribute.NewIntEditor("Center y", 0, float64(height))
        propEdit.editor["clock"] = attribute.NewClockEditor("Time")
        propEdit.editor["color"] = attribute.NewColorEditor("Color")

    case IMAGE_PROP:
        propEdit.editor["x"] = attribute.NewIntEditor("Center x", 0, float64(width))
        propEdit.editor["y"] = attribute.NewIntEditor("Center y", 0, float64(height))
        propEdit.editor["string"] = attribute.NewStringEditor("Image")
        propEdit.editor["scale"] = attribute.NewFloatEditor("Scale", 0.01, 10, 0.01)

    default:
        return nil, fmt.Errorf("Unknown Prop %d", typed)
    }

    propEdit.AddEditors()
    return
}

var propOrder = map[int][]string {
    RECT_PROP: { "x", "y", "width", "height", "rounding", "color" },
    TEXT_PROP: { "x", "y", "color", "string" },
    CIRCLE_PROP: { "x", "y", "inner_radius", "outer_radius", "start_angle", "end_angle", "color" },
    GRAPH_PROP: { "x", "y", "color", "node" },
    TICKER_PROP: { "x", "y", "color", "text" },
    CLOCK_PROP: { "x", "y", "color", "clock" },
    IMAGE_PROP: { "x", "y", "scale", "string" },
}

func (propEdit *PropertyEditor) AddEditors() {
    order := propOrder[propEdit.PropType]

    if len(order) == 0 {
        log.Printf("Prop order for %d has length 0", propEdit.PropType)
    }

    for _, name := range order {
        propEdit.Box.PackStart(propEdit.editor[name].Box(), false, false, padding)
    }
}

func (propEdit *PropertyEditor) CreateVisibleEditor() (box *gtk.Box, err error) {
    widthChars := 12 
    order := propOrder[propEdit.PropType]
    propEdit.visible = make(map[string]*gtk.CheckButton)

    box, err = gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0)
    if err != nil {
        return
    }

    box.SetVisible(true)
    
    row, err := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 0)
    if err != nil {
        return 
    }

    attrLabel, err := gtk.LabelNew("Attribute")
    if err != nil {
        return
    }

    visibleLabel, err := gtk.LabelNew("Visible")
    if err != nil {
        return
    }

    row.SetVisible(true)
    attrLabel.SetVisible(true)
    visibleLabel.SetVisible(true)

    attrLabel.SetWidthChars(widthChars)
    visibleLabel.SetWidthChars(widthChars)

    row.PackStart(attrLabel, false, false, padding)
    row.PackStart(visibleLabel, false, false, padding)
    box.PackStart(row, false, false, padding)

    for _, name := range order {
        attr := propEdit.editor[name]

        row, err = gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 0)
        if err != nil {
            return 
        }

        attrLabel, err = gtk.LabelNew(attr.Name())
        if err != nil {
            return
        }

        propEdit.visible[name], err = gtk.CheckButtonNew()
        if err != nil {
            return 
        }

        row.SetVisible(true)
        attrLabel.SetVisible(true)
        attrLabel.SetWidthChars(widthChars)

        propEdit.visible[name].SetVisible(true)
        propEdit.visible[name].SetMarginStart(40)

        row.PackStart(attrLabel, false, false, padding)
        row.PackStart(propEdit.visible[name], false, false, padding)
        box.PackStart(row, false, false, padding)
    }

    return
}

/*
    Update PropertyEditor with the data in Property
*/
func (propEdit *PropertyEditor) UpdateEditor(prop *Property) {
    for name, edit := range propEdit.editor {
        if _, ok := prop.Attr[name]; !ok {
            continue
        }

        edit.Box().SetVisible(prop.Visible[name])
        err := edit.Update(prop.Attr[name])
        if err != nil {
            log.Print(err)
        }
    }
}

func (propEdit *PropertyEditor) UpdateEditorAllProp(prop *Property) {
    for name, edit := range propEdit.editor {
        if _, ok := prop.Attr[name]; !ok {
            continue
        }

        edit.Box().SetVisible(true)
        if check := propEdit.visible[name]; check != nil {
            check.SetActive(prop.Visible[name])
        }

        err := edit.Update(prop.Attr[name])
        if err != nil {
            log.Print(err)
        }
    }
}
