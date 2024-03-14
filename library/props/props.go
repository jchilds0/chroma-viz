package props

import (
	"chroma-viz/library/attribute"
	"encoding/json"
	"fmt"
	"log"
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

func PropType(prop int) string {
    switch prop {
    case RECT_PROP:
        return "rect"
    case TEXT_PROP:
        return "text"
    case CIRCLE_PROP:
        return "circle"
    case GRAPH_PROP:
        return "graph"
    case TICKER_PROP:
        return "ticker"
    case CLOCK_PROP:
        return "clock"
    case IMAGE_PROP:
        return "image"
    default:
        log.Printf("Unknown prop type %d", prop)
        return ""
    }
}

func GeoType(prop int) string {
    switch prop {
    case RECT_PROP:
        return "rect"
    case TEXT_PROP:
        return "text"
    case CIRCLE_PROP:
        return "circle"
    case GRAPH_PROP:
        return "graph"
    case TICKER_PROP:
        return "text"
    case CLOCK_PROP:
        return "text"
    case IMAGE_PROP:
        return "image"
    default:
        log.Printf("Unknown geo type %d", prop)
        return ""
    }
}

var PropToString map[int]string = map[int]string{
    RECT_PROP: "rect",
    TEXT_PROP: "text",
    CIRCLE_PROP: "circle",
    GRAPH_PROP: "graph",
    TICKER_PROP: "ticker",
    CLOCK_PROP: "clock",
    IMAGE_PROP: "image",
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

*/

type Property struct {
    Name      string
    PropType  int
    Visible   map[string]bool
    Attr      map[string]attribute.Attribute
}

func NewProperty(typed int, name string, visible map[string]bool, cont func()) *Property {
    prop := &Property{Name: name, PropType: typed, Visible: visible}

    if visible == nil {
        prop.Visible = make(map[string]bool)
    }

    prop.Attr = make(map[string]attribute.Attribute, 10)

    switch (typed) {
    case RECT_PROP:
        prop.Attr["x"] = attribute.NewIntAttribute("rel_x")
        prop.Attr["y"] = attribute.NewIntAttribute("rel_y")
        prop.Attr["width"] = attribute.NewIntAttribute("width")
        prop.Attr["height"] = attribute.NewIntAttribute("height")
        prop.Attr["color"] = attribute.NewColorAttribute("color")

    case TEXT_PROP:
        prop.Attr["x"] = attribute.NewIntAttribute("rel_x")
        prop.Attr["y"] = attribute.NewIntAttribute("rel_y")
        prop.Attr["string"] = attribute.NewStringAttribute("string")
        prop.Attr["color"] = attribute.NewColorAttribute("color")

    case CIRCLE_PROP:
        prop.Attr["x"] = attribute.NewIntAttribute("rel_x")
        prop.Attr["y"] = attribute.NewIntAttribute("rel_y")
        prop.Attr["inner_radius"] = attribute.NewIntAttribute("inner_radius")
        prop.Attr["outer_radius"] = attribute.NewIntAttribute("outer_radius")
        prop.Attr["start_angle"] = attribute.NewIntAttribute("start_angle")
        prop.Attr["end_angle"] = attribute.NewIntAttribute("end_angle")
        prop.Attr["color"] = attribute.NewColorAttribute("color")

    case GRAPH_PROP:
        prop.Attr["x"] = attribute.NewIntAttribute("rel_x")
        prop.Attr["y"] = attribute.NewIntAttribute("rel_y")
        prop.Attr["node"] = attribute.NewListAttribute("graph_node", 2, false)
        prop.Attr["color"] = attribute.NewColorAttribute("color")

    case TICKER_PROP:
        prop.Attr["x"] = attribute.NewIntAttribute("rel_x")
        prop.Attr["y"] = attribute.NewIntAttribute("rel_y")
        prop.Attr["text"] = attribute.NewListAttribute("string", 1, true)
        prop.Attr["color"] = attribute.NewColorAttribute("color")

    case CLOCK_PROP:
        prop.Attr["x"] = attribute.NewIntAttribute("rel_x")
        prop.Attr["y"] = attribute.NewIntAttribute("rel_y")
        prop.Attr["clock"] = attribute.NewClockAttribute("string", cont)
        prop.Attr["color"] = attribute.NewColorAttribute("color")

    case IMAGE_PROP:
        prop.Attr["x"] = attribute.NewIntAttribute("rel_x")
        prop.Attr["y"] = attribute.NewIntAttribute("rel_y")
        prop.Attr["scale"] = attribute.NewFloatAttribute("scale")
        prop.Attr["string"] = attribute.NewStringAttribute("string")

    default:
        log.Printf("Unknown Prop %d", typed)
        return nil
    }

    return prop
}

type PropertyJSON struct {
    Name      string
    PropType  int
    Visible   map[string]bool
    Attr      map[string]attribute.AttributeJSON
}

func (prop *Property) UnmarshalJSON(b []byte) error {
    var tempProp PropertyJSON

    err := json.Unmarshal(b, &tempProp)
    if err != nil {
        return err
    }

    prop.Name = tempProp.Name
    prop.PropType = tempProp.PropType
    prop.Visible = tempProp.Visible
    prop.Attr = make(map[string]attribute.Attribute, 10)

    for name, attrJSON := range tempProp.Attr {
        prop.Attr[name] = attrJSON.Attr()
    }

    return nil
}

/*
    Convert a Property to a string to be sent to Chroma Engine
*/
func (prop *Property) String() (s string) {
    for name, attr := range prop.Attr {
        if !prop.Visible[name] {
            continue
        }

        s += attr.String()
    }
    return
}

// G -> {'id': 123, 'name': 'abc', 'prop_type': 'abc', 'geo_type': 'abc', 'visible': [...], 'attr': [A]} | G, G
func (prop *Property) Encode(geo_id int) string {
    first := true 
    attrs := ""

    for _, attr := range prop.Attr {
        if first {
            attrs = attr.Encode()
            first = false 
            continue
        }

        attrs = fmt.Sprintf("%s,%s", attrs, attr.Encode())
    }

    visible := ""
    first = true
    for name, vis := range prop.Visible {
        if first {
            visible = fmt.Sprintf("'%s': '%v'", name, vis)
            first = false
            continue
        }

        visible += fmt.Sprintf(",'%s': '%v'", name, vis)
    }

    return fmt.Sprintf("{'id': %d, 'name': '%s', 'prop_type': '%s', 'geo_type': '%s', 'visible': [%s], 'attr': [%s]}", 
        geo_id, prop.Name, PropType(prop.PropType), GeoType(prop.PropType), visible, attrs)
}

/*
    Update Property with the data in PropertyEditor
*/
func (prop *Property) UpdateProp(propEdit *PropertyEditor) {
    editors := propEdit.editor

    for name, attr := range prop.Attr {
        if _, ok := editors[name]; !ok{
            continue
        }

        err := attr.Update(editors[name])
        if err != nil {
            log.Print(err)
        }
    }
}
