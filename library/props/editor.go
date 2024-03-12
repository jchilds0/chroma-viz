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
    the Page is sent to the editor, which calls UpdateEditor. 
    This updates the PropertyEditor with the data from the 
    Property.

    Similarly UpdateProp is used to send the updated data from 
    the PropertyEditor's back to the Property.
    
    The Properties are built up from a collection of Attributes,
    and in a similar way, PropertyEditor's are built up from a 
    collection of AttributeEditors.

*/

type PropertyEditor struct {
    PropType    int
    Box         *gtk.Box
    editor      map[string]attribute.Editor
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
        propEdit.editor["x"] = attribute.NewIntEditor("x", 0, float64(width))
        propEdit.editor["y"] = attribute.NewIntEditor("y", 0, float64(height))
        propEdit.editor["width"] = attribute.NewIntEditor("Width", 0, float64(width))
        propEdit.editor["height"] = attribute.NewIntEditor("Height", 0, float64(height))
        propEdit.editor["color"] = attribute.NewColorEditor("Color")

        propEdit.Box.PackStart(propEdit.editor["x"].Box(), false, false, padding)
        propEdit.Box.PackStart(propEdit.editor["y"].Box(), false, false, padding)
        propEdit.Box.PackStart(propEdit.editor["width"].Box(), false, false, padding)
        propEdit.Box.PackStart(propEdit.editor["height"].Box(), false, false, padding)
        propEdit.Box.PackStart(propEdit.editor["color"].Box(), false, false, padding)
    
    case TEXT_PROP:
        propEdit.editor["x"] = attribute.NewIntEditor("Center x", 0, float64(width))
        propEdit.editor["y"] = attribute.NewIntEditor("Center y", 0, float64(height))
        propEdit.editor["string"] = attribute.NewStringEditor("Text")
        propEdit.editor["color"] = attribute.NewColorEditor("Color")

        propEdit.Box.PackStart(propEdit.editor["x"].Box(), false, false, padding)
        propEdit.Box.PackStart(propEdit.editor["y"].Box(), false, false, padding)
        propEdit.Box.PackStart(propEdit.editor["string"].Box(), false, false, padding)
        propEdit.Box.PackStart(propEdit.editor["color"].Box(), false, false, padding)

    case CIRCLE_PROP:
        propEdit.editor["x"] = attribute.NewIntEditor("Center x", 0, float64(width))
        propEdit.editor["y"] = attribute.NewIntEditor("Center y", 0, float64(height))
        propEdit.editor["inner_radius"] = attribute.NewIntEditor("Inner Radius", 0, float64(width))
        propEdit.editor["outer_radius"] = attribute.NewIntEditor("Outer Radius", 0, float64(width))
        propEdit.editor["start_angle"] = attribute.NewIntEditor("Start Angle", 0, 360)
        propEdit.editor["end_angle"] = attribute.NewIntEditor("End Angle", 0, 360)
        propEdit.editor["color"] = attribute.NewColorEditor("Color")

        propEdit.Box.PackStart(propEdit.editor["x"].Box(), false, false, padding)
        propEdit.Box.PackStart(propEdit.editor["y"].Box(), false, false, padding)
        propEdit.Box.PackStart(propEdit.editor["inner_radius"].Box(), false, false, padding)
        propEdit.Box.PackStart(propEdit.editor["outer_radius"].Box(), false, false, padding)
        propEdit.Box.PackStart(propEdit.editor["start_angle"].Box(), false, false, padding)
        propEdit.Box.PackStart(propEdit.editor["end_angle"].Box(), false, false, padding)
        propEdit.Box.PackStart(propEdit.editor["color"].Box(), false, false, padding)

    case GRAPH_PROP:
        propEdit.editor["x"] = attribute.NewIntEditor("Center x", 0, float64(width))
        propEdit.editor["y"] = attribute.NewIntEditor("Center y", 0, float64(height))
        propEdit.editor["node"] = attribute.NewListEditor("Graph", []string{"x Pos", "y Pos"})
        propEdit.editor["color"] = attribute.NewColorEditor("Color")

        propEdit.Box.PackStart(propEdit.editor["x"].Box(), false, false, padding)
        propEdit.Box.PackStart(propEdit.editor["y"].Box(), false, false, padding)
        propEdit.Box.PackStart(propEdit.editor["node"].Box(), false, false, padding)
        propEdit.Box.PackStart(propEdit.editor["color"].Box(), false, false, padding)

    case TICKER_PROP:
        propEdit.editor["x"] = attribute.NewIntEditor("Center x", 0, float64(width))
        propEdit.editor["y"] = attribute.NewIntEditor("Center y", 0, float64(height))
        propEdit.editor["color"] = attribute.NewColorEditor("Color")
        propEdit.editor["text"] = attribute.NewListEditor("Ticker", []string{"Text"})
    
        propEdit.Box.PackStart(propEdit.editor["x"].Box(), false, false, padding)
        propEdit.Box.PackStart(propEdit.editor["y"].Box(), false, false, padding)
        propEdit.Box.PackStart(propEdit.editor["text"].Box(), false, false, padding)

    case CLOCK_PROP:
        propEdit.editor["x"] = attribute.NewIntEditor("Center x", 0, float64(width))
        propEdit.editor["y"] = attribute.NewIntEditor("Center y", 0, float64(height))
        propEdit.editor["clock"] = attribute.NewClockEditor("Time")
        propEdit.editor["color"] = attribute.NewColorEditor("Color")

        propEdit.Box.PackStart(propEdit.editor["x"].Box(), false, false, padding)
        propEdit.Box.PackStart(propEdit.editor["y"].Box(), false, false, padding)
        propEdit.Box.PackStart(propEdit.editor["clock"].Box(), false, false, padding)
        propEdit.Box.PackStart(propEdit.editor["color"].Box(), false, false, padding)
    case IMAGE_PROP:
        propEdit.editor["x"] = attribute.NewIntEditor("Center x", 0, float64(width))
        propEdit.editor["y"] = attribute.NewIntEditor("Center y", 0, float64(height))
        propEdit.editor["string"] = attribute.NewStringEditor("Image")
        propEdit.editor["scale"] = attribute.NewFloatEditor("Scale", 0.01, 10, 0.01)

        propEdit.Box.PackStart(propEdit.editor["x"].Box(), false, false, padding)
        propEdit.Box.PackStart(propEdit.editor["y"].Box(), false, false, padding)
        propEdit.Box.PackStart(propEdit.editor["string"].Box(), false, false, padding)
        propEdit.Box.PackStart(propEdit.editor["scale"].Box(), false, false, padding)
    default:
        return nil, fmt.Errorf("Unknown Prop %d", typed)
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

