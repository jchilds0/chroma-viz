package templates

import (
	"chroma-viz/library/geometry"
	"log"

	"github.com/gotk3/gotk3/gtk"
)

type Editor struct {
	Box          *gtk.Box
	tabs         *gtk.Notebook
	actions      *gtk.Box
	CurrentGeoID int

	Rect   *geometry.RectangleEditor
	Circle *geometry.CircleEditor
	Clock  *geometry.ClockEditor
	Image  *geometry.ImageEditor
	Poly   *geometry.PolygonEditor
	Text   *geometry.TextEditor
	List   *geometry.ListEditor
}

func NewEditor() (editor *Editor, err error) {
	editor = &Editor{}

	editor.Box, err = gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0)
	if err != nil {
		return
	}

	editor.actions, err = gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 0)
	if err != nil {
		return
	}

	editor.Box.PackStart(editor.actions, false, false, 10)

	editor.tabs, err = gtk.NotebookNew()
	if err != nil {
		return
	}

	tab, err := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0)
	if err != nil {
		return
	}

	tabLabel, err := gtk.LabelNew("Select A Page")
	if err != nil {
		return
	}

	editor.tabs.AppendPage(tab, tabLabel)
	editor.Box.PackStart(editor.tabs, true, true, 0)

	editor.Rect, err = geometry.NewRectangleEditor()
	if err != nil {
		log.Print(err)
	}

	editor.Circle, err = geometry.NewCircleEditor()
	if err != nil {
		log.Print(err)
	}

	editor.Clock, err = geometry.NewClockEditor()
	if err != nil {
		log.Print(err)
	}

	editor.Image, err = geometry.NewImageEditor()
	if err != nil {
		log.Print(err)
	}

	editor.Poly, err = geometry.NewPolygonEditor()
	if err != nil {
		log.Print(err)
	}

	editor.Text, err = geometry.NewTextEditor()
	if err != nil {
		log.Print(err)
	}

	editor.List, err = geometry.NewListEditor()
	if err != nil {
		log.Print(err)
	}

	return
}

/*
 * Add a button to the editor box with label 'label'
 * which triggers action 'action' when clicked.
 */
func (editor *Editor) AddAction(label string, start bool, action func()) (err error) {
	button, err := gtk.ButtonNewWithLabel(label)
	if err != nil {
		return
	}

	button.Connect("clicked", action)

	if start {
		editor.actions.PackStart(button, false, false, 10)
	} else {
		editor.actions.PackEnd(button, false, false, 10)
	}

	return
}

// Move values from editors to geometries
func (edit *Editor) UpdateGeometry(temp *Template) {
	updateGeometry[*geometry.Rectangle, *geometry.RectangleEditor](
		temp.Rectangle, edit.Rect, edit.CurrentGeoID,
	)

	updateGeometry[*geometry.Circle, *geometry.CircleEditor](
		temp.Circle, edit.Circle, edit.CurrentGeoID,
	)

	updateGeometry[*geometry.Clock, *geometry.ClockEditor](
		temp.Clock, edit.Clock, edit.CurrentGeoID,
	)

	updateGeometry[*geometry.Image, *geometry.ImageEditor](
		temp.Image, edit.Image, edit.CurrentGeoID,
	)

	updateGeometry[*geometry.Polygon, *geometry.PolygonEditor](
		temp.Polygon, edit.Poly, edit.CurrentGeoID,
	)

	updateGeometry[*geometry.Text, *geometry.TextEditor](
		temp.Text, edit.Text, edit.CurrentGeoID,
	)

	updateGeometry[*geometry.List, *geometry.ListEditor](
		temp.List, edit.List, edit.CurrentGeoID,
	)

}

type geometer[S any] interface {
	UpdateGeometry(S) error
	GetName() string
	GetGeometryID() int
	GetGeometry() *geometry.Geometry
}

func updateGeometry[T geometer[S], S any](geos []T, edit S, geoID int) {
	for _, geo := range geos {
		if geo.GetGeometryID() != geoID {
			continue
		}

		err := geo.UpdateGeometry(edit)
		if err != nil {
			log.Print(err)
		}
	}

	return
}

// Load a geometry into the editor
func (edit *Editor) UpdateEditor(temp *Template) (err error) {
	num_pages := edit.tabs.GetNPages()
	for i := 0; i < num_pages; i++ {
		edit.tabs.RemovePage(0)
	}

	updateEditor[*geometry.RectangleEditor, *geometry.Rectangle](
		edit, edit.Rect, temp.Rectangle,
	)

	updateEditor[*geometry.CircleEditor, *geometry.Circle](
		edit, edit.Circle, temp.Circle,
	)

	updateEditor[*geometry.ClockEditor, *geometry.Clock](
		edit, edit.Clock, temp.Clock,
	)

	updateEditor[*geometry.ImageEditor, *geometry.Image](
		edit, edit.Image, temp.Image,
	)

	updateEditor[*geometry.PolygonEditor, *geometry.Polygon](
		edit, edit.Poly, temp.Polygon,
	)

	updateEditor[*geometry.TextEditor, *geometry.Text](
		edit, edit.Text, temp.Text,
	)

	updateEditor[*geometry.ListEditor, *geometry.List](
		edit, edit.List, temp.List,
	)

	return
}

type editor[S any] interface {
	UpdateEditor(S) error
	GetBox() *gtk.ScrolledWindow
	GetVisibleBox() *gtk.ScrolledWindow
}

func updateEditor[T editor[S], S geometer[T]](edit *Editor, editor T, geos []S) {
	geoLabel, err := gtk.LabelNew("Geometry")
	if err != nil {
		return
	}

	//visibleLabel, err := gtk.LabelNew("Visible")
	if err != nil {
		return
	}

	for _, geo := range geos {
		if geo.GetGeometryID() != edit.CurrentGeoID {
			continue
		}

		editor.UpdateEditor(geo)
		edit.tabs.AppendPage(editor.GetBox(), geoLabel)
		//edit.tabs.AppendPage(editor.GetVisibleBox(), visibleLabel)
	}
}
