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

func (g *Geometry) UpdateGeometry(gEdit *GeometryEditor) (err error) {
	return
}

func (g *Geometry) EncodeEngine(b strings.Builder) {

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
	return
}
