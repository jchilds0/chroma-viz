package geometry

import (
	"chroma-viz/library/attribute"
	"chroma-viz/library/util"
	"fmt"
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
	GEO_LIST   = "List"
	GEO_CLOCK  = "Clock"
)

const (
	ATTR_REL_X        = "rel_x"
	ATTR_REL_Y        = "rel_y"
	ATTR_WIDTH        = "width"
	ATTR_HEIGHT       = "height"
	ATTR_X_LOWER      = "x_lower"
	ATTR_X_UPPER      = "x_upper"
	ATTR_Y_LOWER      = "y_lower"
	ATTR_Y_UPPER      = "y_upper"
	ATTR_INNER_RADIUS = "inner_radius"
	ATTR_OUTER_RADIUS = "outer_radius"
	ATTR_START_ANGLE  = "start_angle"
	ATTR_END_ANGLE    = "end_angle"
	ATTR_MASK         = "mask"
	ATTR_ROUND        = "rounding"
	ATTR_COLOR        = "color"
	ATTR_COLOR_R      = "red"
	ATTR_COLOR_G      = "green"
	ATTR_COLOR_B      = "blue"
	ATTR_COLOR_A      = "alpha"
	ATTR_PARENT       = "parent"
	ATTR_STRING       = "string"
	ATTR_SCALE        = "scale"
	ATTR_POINT        = "point"
	ATTR_NUM_POINTS   = "num_points"
)

var Attrs = map[string]string{
	ATTR_REL_X:        "Rel X",
	ATTR_REL_Y:        "Rel Y",
	ATTR_WIDTH:        "Width",
	ATTR_HEIGHT:       "Height",
	ATTR_X_LOWER:      "X Lower",
	ATTR_X_UPPER:      "X Upper",
	ATTR_Y_LOWER:      "Y Lower",
	ATTR_Y_UPPER:      "Y Upper",
	ATTR_INNER_RADIUS: "Inner Radius",
	ATTR_OUTER_RADIUS: "Outer Radius",
	ATTR_START_ANGLE:  "Start Angle",
	ATTR_END_ANGLE:    "End Angle",
	ATTR_MASK:         "Mask",
	ATTR_ROUND:        "Rounding",
	ATTR_COLOR:        "Color",
	ATTR_COLOR_R:      "Red",
	ATTR_COLOR_G:      "Green",
	ATTR_COLOR_B:      "Blue",
	ATTR_COLOR_A:      "Alpha",
	ATTR_PARENT:       "Parent",
	ATTR_STRING:       "String",
	ATTR_SCALE:        "Scale",
}

func UpdateAttrList(model *gtk.ListStore, geoType string) (err error) {
	model.Clear()

	attrs := []string{
		ATTR_REL_X, ATTR_REL_Y, ATTR_X_LOWER,
		ATTR_X_UPPER, ATTR_Y_LOWER, ATTR_Y_UPPER,
	}

	switch geoType {
	case GEO_RECT:
		attrs = append(attrs, ATTR_WIDTH, ATTR_HEIGHT)
		attrs = append(attrs, ATTR_COLOR_R, ATTR_COLOR_G, ATTR_COLOR_B, ATTR_COLOR_A)

	case GEO_CIRCLE:
		attrs = append(attrs, ATTR_INNER_RADIUS, ATTR_OUTER_RADIUS)
		attrs = append(attrs, ATTR_START_ANGLE, ATTR_END_ANGLE)
		attrs = append(attrs, ATTR_COLOR_R, ATTR_COLOR_G, ATTR_COLOR_B, ATTR_COLOR_A)

	case GEO_TEXT:
		attrs = append(attrs, ATTR_WIDTH, ATTR_HEIGHT)
		attrs = append(attrs, ATTR_SCALE)
		attrs = append(attrs, ATTR_COLOR_R, ATTR_COLOR_G, ATTR_COLOR_B, ATTR_COLOR_A)

	case GEO_IMAGE:
		attrs = append(attrs, ATTR_SCALE)
	case GEO_POLY:
		attrs = append(attrs, ATTR_COLOR_R, ATTR_COLOR_G, ATTR_COLOR_B, ATTR_COLOR_A)

	case GEO_LIST:
	case GEO_CLOCK:

	default:
		return fmt.Errorf("Unknown geometry type %s", geoType)
	}

	for _, name := range attrs {
		iter := model.Append()

		model.SetValue(iter, 0, name)
		model.SetValue(iter, 1, Attrs[name])
	}

	return
}

type Geometer[S any] interface {
	UpdateGeometry(S) error
	GetName() string
	GetGeometry() *Geometry
	GetGeometryID() int
	Encode(*strings.Builder)
}

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
	util.EngineAddKeyValue(b, "geo_num", g.GeometryID)
	util.EngineAddKeyValue(b, g.RelX.Name, g.RelX.Value)
	util.EngineAddKeyValue(b, g.RelY.Name, g.RelY.Value)
	util.EngineAddKeyValue(b, g.Mask.Name, g.Mask.Value)
}

func (g Geometry) GetName() string {
	return g.Name
}

func (g Geometry) GetGeometryID() int {
	return g.GeometryID
}

func (g Geometry) GetGeometry() *Geometry {
	return &g
}

type Editor[S any] interface {
	UpdateEditor(S) error
	GetBox() *gtk.ScrolledWindow
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
