package props

import (
	"chroma-viz/library/attribute"
	"chroma-viz/library/templates"
	"encoding/json"
	"log"
	"strconv"
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
	RECT_PROP   = "rect"
	TEXT_PROP   = "text"
	CIRCLE_PROP = "circle"
	TICKER_PROP = "ticker"
	CLOCK_PROP  = "clock"
	IMAGE_PROP  = "image"
)

var PropToGeo = map[string]string{
	RECT_PROP:   templates.GEO_RECT,
	TEXT_PROP:   templates.GEO_TEXT,
	CIRCLE_PROP: templates.GEO_CIRCLE,
	TICKER_PROP: templates.GEO_TEXT,
	CLOCK_PROP:  templates.GEO_TEXT,
	IMAGE_PROP:  templates.GEO_IMAGE,
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
	PropType string
	Visible  map[string]bool
	Attr     map[string]attribute.Attribute
	temp     bool
}

func NewProperty(typed, name string, isTemp bool, visible map[string]bool) *Property {
	prop := &Property{Name: name, PropType: typed, Visible: visible, temp: isTemp}

	if visible == nil {
		prop.Visible = make(map[string]bool)
	}

	prop.Attr = make(map[string]attribute.Attribute, 10)

	prop.Attr["rel_x"] = attribute.NewIntAttribute("rel_x")
	prop.Attr["rel_y"] = attribute.NewIntAttribute("rel_y")
	prop.Attr["parent"] = attribute.NewIntAttribute("parent")

	switch typed {
	case RECT_PROP:
		prop.Attr["width"] = attribute.NewIntAttribute("width")
		prop.Attr["height"] = attribute.NewIntAttribute("height")
		prop.Attr["rounding"] = attribute.NewIntAttribute("rounding")
		prop.Attr["color"] = attribute.NewColorAttribute("color")

	case TEXT_PROP:
		prop.Attr["string"] = attribute.NewStringAttribute("string")
		prop.Attr["color"] = attribute.NewColorAttribute("color")

	case CIRCLE_PROP:
		prop.Attr["inner_radius"] = attribute.NewIntAttribute("inner_radius")
		prop.Attr["outer_radius"] = attribute.NewIntAttribute("outer_radius")
		prop.Attr["start_angle"] = attribute.NewIntAttribute("start_angle")
		prop.Attr["end_angle"] = attribute.NewIntAttribute("end_angle")
		prop.Attr["color"] = attribute.NewColorAttribute("color")

	case TICKER_PROP:
		prop.Attr["string"], _ = attribute.NewListAttribute("string", 1, true)
		prop.Attr["color"] = attribute.NewColorAttribute("color")

	case CLOCK_PROP:
		prop.Attr["string"] = attribute.NewClockAttribute("string")
		prop.Attr["color"] = attribute.NewColorAttribute("color")

	case IMAGE_PROP:
		prop.Attr["scale"] = attribute.NewFloatAttribute("scale")
		prop.Attr["image_id"] = attribute.NewAssetAttribute("image_id")

	default:
		return nil
	}

	return prop
}

func NewPropertyFromGeometry(geo templates.Geometry, attrMap map[string]string) (prop *Property) {
	prop = NewProperty(geo.PropType, geo.Name, false, nil)

	for name, value := range attrMap {
		var err error

		if attr, ok := prop.Attr[name]; ok {
			err = attr.Decode(value)
		}

		if err != nil {
			log.Print(err)
		}

		prop.Visible[name] = true
	}

	return
}

func (prop *Property) CreateGeometry(temp *templates.Template, geoID int) {
	var relX, relY, parent int
	if attr, ok := prop.Attr["rel_x"]; ok {
		relX, _ = strconv.Atoi(attr.Encode())
	}

	if attr, ok := prop.Attr["rel_y"]; ok {
		relY, _ = strconv.Atoi(attr.Encode())
	}

	if attr, ok := prop.Attr["parent"]; ok {
		parent, _ = strconv.Atoi(attr.Encode())
	}

	geo := templates.NewGeometry(geoID, prop.Name, prop.PropType, PropToGeo[prop.PropType], relX, relY, parent)

	switch prop.PropType {
	case RECT_PROP:
		var width, height, rounding int
		color := "0 0 0 0"

		if attr, ok := prop.Attr["width"]; ok {
			width, _ = strconv.Atoi(attr.Encode())
		}

		if attr, ok := prop.Attr["height"]; ok {
			height, _ = strconv.Atoi(attr.Encode())
		}

		if attr, ok := prop.Attr["rounding"]; ok {
			rounding, _ = strconv.Atoi(attr.Encode())
		}

		if attr, ok := prop.Attr["color"]; ok {
			color = attr.Encode()
		}

		rect := templates.NewRectangle(*geo, width, height, rounding, color)
		temp.Rectangle = append(temp.Rectangle, *rect)
	case CIRCLE_PROP:
		var inner, outer, start, end int
		color := "0 0 0 0"

		if attr, ok := prop.Attr["inner_radius"]; ok {
			inner, _ = strconv.Atoi(attr.Encode())
		}

		if attr, ok := prop.Attr["outer_radius"]; ok {
			outer, _ = strconv.Atoi(attr.Encode())
		}

		if attr, ok := prop.Attr["start_angle"]; ok {
			start, _ = strconv.Atoi(attr.Encode())
		}

		if attr, ok := prop.Attr["end_angle"]; ok {
			end, _ = strconv.Atoi(attr.Encode())
		}

		if attr, ok := prop.Attr["color"]; ok {
			color = attr.Encode()
		}

		circle := templates.NewCircle(*geo, inner, outer, start, end, color)
		temp.Circle = append(temp.Circle, *circle)

	case TEXT_PROP, TICKER_PROP, CLOCK_PROP:
		var s string
		color := "0 0 0 0"

		if attr, ok := prop.Attr["string"]; ok {
			s = attr.Encode()
		}

		if attr, ok := prop.Attr["color"]; ok {
			color = attr.Encode()
		}

		text := templates.NewText(*geo, s, color)
		temp.Text = append(temp.Text, *text)

	case IMAGE_PROP:
		var image_id, s string

		if attr, ok := prop.Attr["image_id"]; ok {
			image_id = attr.Encode()
		}

		id, _ := strconv.Atoi(image_id)

		if attr, ok := prop.Attr["scale"]; ok {
			s = attr.Encode()
		}

		scale, _ := strconv.ParseFloat(s, 64)

		a := templates.NewAsset(*geo, "", "", id, scale)
		temp.Asset = append(temp.Asset, *a)

	default:
		log.Printf("Error creating geom %s: Not Implemented", prop.PropType)
	}

	return
}

type PropertyJSON struct {
	Name     string
	PropType string
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
