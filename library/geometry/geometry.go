package geometry

import (
	"chroma-viz/library/attribute"
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
	GEO_RECT   = "rect"
	GEO_CIRCLE = "circle"
	GEO_TEXT   = "text"
	GEO_IMAGE  = "image"
	GEO_POLY   = "polygon"
)

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
	geo.RelX.Name = "rel_x"
	geo.RelY.Name = "rel_y"
	geo.Parent.Name = "parent"
	geo.Mask.Name = "mask"

	return geo
}

func (g *Geometry) UpdateGeometry(gEdit *GeometryEditor) (err error) {
	err = g.RelX.UpdateAttribute(&gEdit.RelX)
	if err != nil {
		return
	}

	err = g.RelY.UpdateAttribute(&gEdit.RelY)
	if err != nil {
		return
	}

	err = g.Mask.UpdateAttribute(&gEdit.Mask)
	return
}

func (g *Geometry) EncodeEngine(b strings.Builder) {

}

func (g *Geometry) GetName() string {
	return g.Name
}

func (g *Geometry) GetGeometryID() int {
	return g.GeometryID
}

type GeometryEditor struct {
	Scroll    *gtk.ScrolledWindow
	ScrollBox *gtk.Box

	RelX attribute.IntEditor
	RelY attribute.IntEditor
	Mask attribute.IntEditor
}

func NewGeometryEditor() *GeometryEditor {
	return nil
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
