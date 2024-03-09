package props

import (
	"chroma-viz/attribute"
	"fmt"
	"log"
	"strings"

	"github.com/gotk3/gotk3/gtk"
)

const padding = 10

const (
    END_OF_CONN = iota + 1
    END_OF_MESSAGE
    ANIMATE_ON
    CONTINUE
    ANIMATE_OFF
)

const (
    RECT_PROP = iota
    TEXT_PROP
    CIRCLE_PROP 
    GRAPH_PROP
    TICKER_PROP
    CLOCK_PROP
    IMAGE_PROP
    NUM_PROPS
)

var StringToProp map[string]int = map[string]int{
    "rect": RECT_PROP,
    "text": TEXT_PROP,
    "circle": CIRCLE_PROP,
    "graph": GRAPH_PROP,
    "ticker": TICKER_PROP,
    "clock": CLOCK_PROP,
    "image": IMAGE_PROP,
}

/*

    Templates are made up of a collection of Properties.

    Property encodes the data needed by chroma engine to display a geometry.
    Each Property has an associated PropertyEditor which generates the gtk
    ui elements needed to edit the corresponding Property. The app generates 
    a small set of PropertyEditor's (enough to show one template) due to the 
    cost of creating gtk ui elements greater than the objects to store the 
    data.

    The user creates a Page from a Template, which involves creating a 
    Property for each Property in the Template. When the user wants to edit 
    the Properties of a Page, the editor uses the Properties of the Page to 
    update PropertyEditor's with the corresponding type (UpdateEditor). 
    In turn the PropertyEditor's update the data stored in Properties on 
    change by the user (UpdateProp).

    Each Property is built up from Attributes, which are simple building 
    blocks like an integer field. We don't always want to show all 
    Attributes so the Property keeps track of the visible Attributes with
    a map, and updates the SetVisible of each gtk element accordingly

    For synchronizing of attributes, the key of the Attributes map in a 
    Property matches the key of the Editors map in a PropertyEditor.
    Each Attribute has a name attribute which is the identifier used when sending 
    the attribute to Chroma Engine. Each Editor also has a name, which is 
    string displayed to the user when editing the attribute.

    PropToString encodes the attributes in a property as a string to be 
    sent to Chroma Engine.

    EncodeProp and DecodeProp are used to import and export a show.

*/

type PropertyEditor interface {
    Box() *gtk.Box
    Editors() map[string]attribute.Editor
}

type Property struct {
    Name      string
    PropType  int
    visible   map[string]bool
    Attr      map[string]attribute.Attribute
}

func NewPropertyEditor(typed int, animate, cont func()) (PropertyEditor, error) {
    width := 1920
    height := 1080 

    switch (typed) {
    case RECT_PROP:
        return NewRectEditor(width, height, animate)
    case TEXT_PROP:
        return NewTextEditor(width, height, animate)
    case CIRCLE_PROP:
        return NewCircleEditor(width, height, animate)
    case GRAPH_PROP:
        return NewGraphEditor(width, height, animate)
    case TICKER_PROP:
        return NewTickerEditor(width, height, animate)
    case CLOCK_PROP:
        return NewClockEditor(width, height, animate, cont)
    case IMAGE_PROP:
        return NewImageEditor(width, height, animate)
    default:
        return nil, fmt.Errorf("Unknown Prop %d", typed)
    }
}

func NewProperty(typed int, name string, visible map[string]bool, cont func()) *Property {
    prop := &Property{Name: name, PropType: typed, visible: visible}

    prop.Attr = make(map[string]attribute.Attribute, 10)

    switch (typed) {
    case RECT_PROP:
        prop.Attr["x"] = attribute.NewIntAttribute("x", "rel_x")
        prop.Attr["y"] = attribute.NewIntAttribute("y", "rel_y")
        prop.Attr["width"] = attribute.NewIntAttribute("width", "width")
        prop.Attr["height"] = attribute.NewIntAttribute("height", "height")
        prop.Attr["color"] = attribute.NewColorAttribute("color", "color")

    case TEXT_PROP:
        prop.Attr["x"] = attribute.NewIntAttribute("x", "rel_x")
        prop.Attr["y"] = attribute.NewIntAttribute("y", "rel_y")
        prop.Attr["string"] = attribute.NewStringAttribute("string","string")
        prop.Attr["color"] = attribute.NewColorAttribute("color", "color")

    case CIRCLE_PROP:
        prop.Attr["x"] = attribute.NewIntAttribute("x", "rel_x")
        prop.Attr["y"] = attribute.NewIntAttribute("y", "rel_y")
        prop.Attr["inner_radius"] = attribute.NewIntAttribute("inner_radius", "inner_radius")
        prop.Attr["outer_radius"] = attribute.NewIntAttribute("outer_radius", "outer_radius")
        prop.Attr["start_angle"] = attribute.NewIntAttribute("start_angle", "start_angle")
        prop.Attr["end_angle"] = attribute.NewIntAttribute("end_angle", "end_angle")
        prop.Attr["color"] = attribute.NewColorAttribute("color", "color")

    case GRAPH_PROP:
        prop.Attr["x"] = attribute.NewIntAttribute("x", "rel_x")
        prop.Attr["y"] = attribute.NewIntAttribute("y", "rel_y")
        prop.Attr["node"] = attribute.NewListAttribute("node", "graph_node", 2, false)
        prop.Attr["color"] = attribute.NewColorAttribute("color", "color")

    case TICKER_PROP:
        prop.Attr["x"] = attribute.NewIntAttribute("x", "rel_x")
        prop.Attr["y"] = attribute.NewIntAttribute("y", "rel_y")
        prop.Attr["text"] = attribute.NewListAttribute("text", "string", 1, true)
        prop.Attr["color"] = attribute.NewColorAttribute("color", "color")

    case CLOCK_PROP:
        prop.Attr["x"] = attribute.NewIntAttribute("x", "rel_x")
        prop.Attr["y"] = attribute.NewIntAttribute("y", "rel_y")
        prop.Attr["clock"] = attribute.NewClockAttribute("string", "string", cont)
        prop.Attr["color"] = attribute.NewColorAttribute("color", "color")

    case IMAGE_PROP:
        prop.Attr["x"] = attribute.NewIntAttribute("x", "rel_x")
        prop.Attr["y"] = attribute.NewIntAttribute("y", "rel_y")
        prop.Attr["scale"] = attribute.NewFloatAttribute("scale", "scale")
        prop.Attr["string"] = attribute.NewStringAttribute("string", "string")

    default:
        log.Printf("Unknown Prop %d", typed)
        return nil
    }

    return prop
}

func (prop *Property) String() (s string) {
    for name, attr := range prop.Attr {
        if !prop.visible[name] {
            continue
        }

        s += attr.String()
    }
    return
}

func (prop *Property) Encode() (s string) {
    for _, attr := range prop.Attr {
        s += attr.Encode()
    }
    return
}

func (prop *Property) Decode(s string) {
    attrs := strings.Split(s, ";")

    for _, attr := range attrs[1:] {
        name := strings.Split(attr, " ")[0]
        if name == "" {
            continue
        }

        if !prop.visible[name] {
            continue
        }

        if prop.Attr[name] == nil {
            log.Printf("Error prop %s missing prop attr %s", prop.Name, name)
            continue
        }

        err := prop.Attr[name].Decode(attr)
        if err != nil {
            log.Printf("Error decoding prop %s in %s", name, prop.Name)
        }
    }
}

/*

    UpdateEditor and UpdateProp are used to synchronize data between
    the Properties and PropertyEditor's. 

    UpdateEditor sends the data contained in the Property to the 
    PropertyEditor. This is called by the Editor object when a 
    user selects a Page.

    UpdateProp sends the data contained in the PropertyEditor to the 
    Property. This is called before any animation action, since the 
    Property object is used to animate to Chroma Engine. As a side 
    effect, all current Editors send an update action to the preview
    when a value is changed. This has the effect of saving the current
    editor state on change, removing the need to have Editor call 
    UpdateProp.

*/

func UpdateEditor(propEdit PropertyEditor, prop *Property) {
    for name, edit := range propEdit.Editors() {
        edit.Box().SetVisible(prop.visible[name])
        err := edit.Update(prop.Attr[name])
        if err != nil {
            log.Print(err)
        }
    }
}

func UpdateProp(prop *Property, propEdit PropertyEditor) {
    editors := propEdit.Editors()

    for name, attr := range prop.Attr {
        err := attr.Update(editors[name])
        if err != nil {
            log.Print(err)
        }
    }
}
