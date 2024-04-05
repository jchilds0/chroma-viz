package props

import (
	"chroma-viz/library/attribute"
	"encoding/json"
	"strconv"
	"strings"
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
	"rect":   RECT_PROP,
	"text":   TEXT_PROP,
	"circle": CIRCLE_PROP,
	"graph":  GRAPH_PROP,
	"ticker": TICKER_PROP,
	"clock":  CLOCK_PROP,
	"image":  IMAGE_PROP,
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
		return ""
	}
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
   the Properties of a Page,

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
	Name     string
	PropType int
	Visible  map[string]bool
	Attr     map[string]attribute.Attribute
	temp     bool
}

func NewProperty(typed int, name string, isTemp bool, visible map[string]bool) *Property {
	prop := &Property{Name: name, PropType: typed, Visible: visible, temp: isTemp}

	if visible == nil {
		prop.Visible = make(map[string]bool)
	}

	prop.Attr = make(map[string]attribute.Attribute, 10)

	switch typed {
	case RECT_PROP:
		prop.Attr["rel_x"] = attribute.NewIntAttribute("rel_x")
		prop.Attr["rel_y"] = attribute.NewIntAttribute("rel_y")
		prop.Attr["width"] = attribute.NewIntAttribute("width")
		prop.Attr["height"] = attribute.NewIntAttribute("height")
		prop.Attr["rounding"] = attribute.NewIntAttribute("rounding")
		prop.Attr["color"] = attribute.NewColorAttribute("color")

	case TEXT_PROP:
		prop.Attr["rel_x"] = attribute.NewIntAttribute("rel_x")
		prop.Attr["rel_y"] = attribute.NewIntAttribute("rel_y")
		prop.Attr["string"] = attribute.NewStringAttribute("string")
		prop.Attr["color"] = attribute.NewColorAttribute("color")

	case CIRCLE_PROP:
		prop.Attr["rel_x"] = attribute.NewIntAttribute("rel_x")
		prop.Attr["rel_y"] = attribute.NewIntAttribute("rel_y")
		prop.Attr["inner_radius"] = attribute.NewIntAttribute("inner_radius")
		prop.Attr["outer_radius"] = attribute.NewIntAttribute("outer_radius")
		prop.Attr["start_angle"] = attribute.NewIntAttribute("start_angle")
		prop.Attr["end_angle"] = attribute.NewIntAttribute("end_angle")
		prop.Attr["color"] = attribute.NewColorAttribute("color")

	case GRAPH_PROP:
		prop.Attr["rel_x"] = attribute.NewIntAttribute("rel_x")
		prop.Attr["rel_y"] = attribute.NewIntAttribute("rel_y")
		prop.Attr["graph_node"], _ = attribute.NewListAttribute("graph_node", 2, false)
		prop.Attr["color"] = attribute.NewColorAttribute("color")

	case TICKER_PROP:
		prop.Attr["rel_x"] = attribute.NewIntAttribute("rel_x")
		prop.Attr["rel_y"] = attribute.NewIntAttribute("rel_y")
		prop.Attr["string"], _ = attribute.NewListAttribute("string", 1, true)
		prop.Attr["color"] = attribute.NewColorAttribute("color")

	case CLOCK_PROP:
		prop.Attr["rel_x"] = attribute.NewIntAttribute("rel_x")
		prop.Attr["rel_y"] = attribute.NewIntAttribute("rel_y")
		prop.Attr["string"] = attribute.NewClockAttribute("string")
		prop.Attr["color"] = attribute.NewColorAttribute("color")

	case IMAGE_PROP:
		prop.Attr["rel_x"] = attribute.NewIntAttribute("rel_x")
		prop.Attr["rel_y"] = attribute.NewIntAttribute("rel_y")
		prop.Attr["scale"] = attribute.NewFloatAttribute("scale")
		prop.Attr["image_id"] = attribute.NewAssetAttribute("image_id")

	default:
		return nil
	}

	return prop
}

type PropertyJSON struct {
	Name     string
	PropType int
	Visible  map[string]bool
	Attr     map[string]attribute.AttributeJSON
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
		if !prop.Visible[name] && !prop.temp {
			continue
		}

		s += attr.String()
	}
	return
}

// G -> {'id': 123, 'name': 'abc', 'prop_type': 'abc', 'geo_type': 'abc', 'visible': [...], 'attr': [A]} | G, G
func (prop *Property) Encode(geo_id int) (s string, err error) {
	var b strings.Builder
	first := true
	b.WriteString("{")

	b.WriteString("'id': ")
	b.WriteString(strconv.Itoa(geo_id))
	b.WriteString(", ")

	b.WriteString("'name': '")
	b.WriteString(prop.Name)
	b.WriteString("', ")

	b.WriteString("'prop_type': '")
	b.WriteString(PropType(prop.PropType))
	b.WriteString("', ")

	b.WriteString("'geo_type': '")
	b.WriteString(GeoType(prop.PropType))
	b.WriteString("', ")

	b.WriteString("'visible': [")
	for name, vis := range prop.Visible {
		if !vis {
			continue
		}

		if !first {
			b.WriteString(",")
		}

		first = false
		b.WriteString("'")
		b.WriteString(name)
		b.WriteString("': 'true'")
	}
	b.WriteString("], ")

	first = true
	b.WriteString("'attr': [")
	for _, attr := range prop.Attr {
		encode := attr.Encode()
		if encode == "" {
			continue
		}

		if !first {
			b.WriteString(",")
		}
		first = false

		b.WriteString(encode)
	}

	b.WriteString("]}")
	s = b.String()

	return
}

/*
Update Property with the data in PropertyEditor
*/
func (prop *Property) UpdateProp(propEdit *PropertyEditor) (err error) {
	editors := propEdit.editor

	for name, attr := range prop.Attr {
		if _, ok := editors[name]; !ok {
			continue
		}

		if check := propEdit.visible[name]; check != nil {
			prop.Visible[name] = check.GetActive()
		}

		err := attr.Update(editors[name])
		if err != nil {
			return err
		}
	}

	return
}

/*
Used by artist to set temp on imported templates
*/
func (prop *Property) SetTemp(isTemp bool) {
	prop.temp = isTemp
}
