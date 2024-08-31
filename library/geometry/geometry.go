package geometry

import (
	"chroma-viz/library/attribute"
	"chroma-viz/library/parser"
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
	GEO_RECT   = "Rectangle"
	GEO_CIRCLE = "Circle"
	GEO_TEXT   = "Text"
	GEO_IMAGE  = "Image"
	GEO_POLY   = "Polygon"
	GEO_TICKER = "Ticker"
	GEO_CLOCK  = "Clock"
)

const (
	ATTR_REL_X        = "rel_x"
	ATTR_REL_Y        = "rel_y"
	ATTR_WIDTH        = "width"
	ATTR_HEIGHT       = "height"
	ATTR_INNER_RADIUS = "inner_radius"
	ATTR_OUTER_RADIUS = "outer_radius"
	ATTR_START_ANGLE  = "start_angle"
	ATTR_END_ANGLE    = "end_angle"
	ATTR_MASK         = "mask"
	ATTR_ROUND        = "rounding"
	ATTR_COLOR        = "color"
	ATTR_PARENT       = "parent"
	ATTR_STRING       = "string"
	ATTR_SCALE        = "scale"
)

var Attrs = map[string]string{
	ATTR_REL_X:        "Rel X",
	ATTR_REL_Y:        "Rel Y",
	ATTR_WIDTH:        "Width",
	ATTR_HEIGHT:       "Height",
	ATTR_INNER_RADIUS: "Inner Radius",
	ATTR_OUTER_RADIUS: "Outer Radius",
	ATTR_START_ANGLE:  "Start Angle",
	ATTR_END_ANGLE:    "End Angle",
	ATTR_MASK:         "Mask",
	ATTR_ROUND:        "Rounding",
	ATTR_COLOR:        "Color",
	ATTR_PARENT:       "Parent",
	ATTR_STRING:       "String",
	ATTR_SCALE:        "Scale",
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

type Geometry struct {
	GeometryID int
	Name       string
	GeoType    string

	RelX   attribute.IntAttribute
	RelY   attribute.IntAttribute
	Parent attribute.IntAttribute
	Mask   attribute.IntAttribute
}

func NewGeometry(geoNum int, name, geoType string) Geometry {
	geo := Geometry{
		GeometryID: geoNum,
		Name:       name,
		GeoType:    geoType,
	}

	// used by chroma engine for attribute identifier
	geo.RelX.Name = ATTR_REL_X
	geo.RelY.Name = ATTR_REL_Y
	geo.Parent.Name = ATTR_PARENT
	geo.Mask.Name = ATTR_MASK

	return geo
}

func (g *Geometry) UpdateGeometry(gEdit *GeometryEditor) (err error) {
	err = g.RelX.UpdateAttribute(gEdit.RelX)
	if err != nil {
		return
	}

	err = g.RelY.UpdateAttribute(gEdit.RelY)
	if err != nil {
		return
	}

	err = g.Mask.UpdateAttribute(gEdit.Mask)
	return
}

func (g *Geometry) Encode(b *strings.Builder) {
	parser.EngineAddKeyValue(b, "geo_num", g.GeometryID)
	parser.EngineAddKeyValue(b, g.RelX.Name, g.RelX.Value)
	parser.EngineAddKeyValue(b, g.RelY.Name, g.RelY.Value)
	parser.EngineAddKeyValue(b, g.Mask.Name, g.Mask.Value)
}

func (g *Geometry) GetName() string {
	return g.Name
}

func (g *Geometry) GetGeometryID() int {
	return g.GeometryID
}

func (g *Geometry) GetGeometry() *Geometry {
	return g
}

type GeometryEditor struct {
	Scroll    *gtk.ScrolledWindow
	ScrollBox *gtk.Box

	RelX *attribute.IntEditor
	RelY *attribute.IntEditor
	Mask *attribute.IntEditor
}

func NewGeometryEditor() (geoEdit *GeometryEditor, err error) {
	geoEdit = &GeometryEditor{}
	geoEdit.Scroll, err = gtk.ScrolledWindowNew(nil, nil)
	if err != nil {
		return
	}

	geoEdit.Scroll.SetVisible(true)
	geoEdit.Scroll.SetVExpand(true)

	geoEdit.ScrollBox, err = gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0)
	if err != nil {
		return
	}

	geoEdit.ScrollBox.SetVisible(true)
	geoEdit.Scroll.Add(geoEdit.ScrollBox)

	geoEdit.RelX, err = attribute.NewIntEditor(Attrs[ATTR_REL_X])
	if err != nil {
		return
	}

	geoEdit.RelY, err = attribute.NewIntEditor(Attrs[ATTR_REL_Y])
	if err != nil {
		return
	}

	geoEdit.Mask, err = attribute.NewIntEditor(Attrs[ATTR_MASK])
	if err != nil {
		return
	}

	geoEdit.ScrollBox.PackStart(geoEdit.RelX.Box, false, false, padding)
	geoEdit.ScrollBox.PackStart(geoEdit.RelY.Box, false, false, padding)
	geoEdit.ScrollBox.PackStart(geoEdit.Mask.Box, false, false, padding)

	return
}

func (gEdit *GeometryEditor) UpdateEditor(g *Geometry) (err error) {
	err = gEdit.RelX.UpdateEditor(&g.RelX)
	if err != nil {
		return
	}

	err = gEdit.RelY.UpdateEditor(&g.RelY)
	if err != nil {
		return
	}

	err = gEdit.Mask.UpdateEditor(&g.Mask)
	return
}

func (gEdit *GeometryEditor) GetBox() *gtk.ScrolledWindow {
	return gEdit.Scroll
}

func (gEdit *GeometryEditor) GetVisibleBox() *gtk.ScrolledWindow {
	return nil
}
