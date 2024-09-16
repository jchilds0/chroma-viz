package templates

import (
	"chroma-viz/library/geometry"
	"log"

	"github.com/gotk3/gotk3/gtk"
)

type Editor struct {
	Box     *gtk.Box
	tabs    *gtk.Notebook
	actions *gtk.Box

	CurrentFrameID int
	CurrentGeoID   int
	CurrentKeyID   int
	geometryTab    *gtk.Box
	visibleTab     *gtk.Box
	keyframeTab    *gtk.Box

	SetFrameEdit  *SetFrameEditor
	BindFrameEdit *BindFrameEditor

	Rect   *geometry.RectangleEditor
	Circle *geometry.CircleEditor
	Clock  *geometry.ClockEditor
	Image  *geometry.ImageEditor
	Poly   *geometry.PolygonEditor
	Text   *geometry.TextEditor
	List   *geometry.ListEditor
}

func NewEditor(frameModel *gtk.ListStore, geoModel *gtk.ListStore) (editor *Editor, err error) {
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

	editor.Box.PackStart(editor.tabs, true, true, 0)

	var label *gtk.Label
	{
		editor.geometryTab, err = gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0)
		if err != nil {
			return
		}

		label, err = gtk.LabelNew("Geometry")
		if err != nil {
			return
		}

		editor.tabs.AppendPage(editor.geometryTab, label)
	}

	{
		editor.visibleTab, err = gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0)
		if err != nil {
			return
		}

		label, err = gtk.LabelNew("Visible")
		if err != nil {
			return
		}

		editor.tabs.AppendPage(editor.visibleTab, label)
	}

	{
		editor.keyframeTab, err = gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0)
		if err != nil {
			return
		}

		label, err = gtk.LabelNew("Keyframe")
		if err != nil {
			return
		}

		editor.tabs.AppendPage(editor.keyframeTab, label)
	}

	editor.SetFrameEdit, err = NewSetFrameEditor()
	if err != nil {
		return
	}

	editor.BindFrameEdit, err = NewBindFrameEditor(frameModel, geoModel)
	if err != nil {
		return
	}

	editor.Rect, err = geometry.NewRectangleEditor()
	if err != nil {
		return
	}

	editor.Circle, err = geometry.NewCircleEditor()
	if err != nil {
		return
	}

	editor.Clock, err = geometry.NewClockEditor()
	if err != nil {
		return
	}

	editor.Image, err = geometry.NewImageEditor()
	if err != nil {
		return
	}

	editor.Poly, err = geometry.NewPolygonEditor()
	if err != nil {
		return
	}

	editor.Text, err = geometry.NewTextEditor()
	if err != nil {
		return
	}

	editor.List, err = geometry.NewListEditor()
	if err != nil {
		return
	}

	editor.Clear()

	return
}

func (edit *Editor) Clear() {
	clearBox(edit.geometryTab)
	clearBox(edit.keyframeTab)
	clearBox(edit.visibleTab)
}

func clearBox(box *gtk.Box) {
	children := box.GetChildren()
	if children == nil {
		return
	}

	children.Foreach(func(child any) {
		box.Remove(child.(*gtk.Widget))
	})
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
	updateGeometry(temp.Rectangle, edit.Rect, edit.CurrentGeoID)
	updateGeometry(temp.Circle, edit.Circle, edit.CurrentGeoID)
	updateGeometry(temp.Clock, edit.Clock, edit.CurrentGeoID)
	updateGeometry(temp.Image, edit.Image, edit.CurrentGeoID)
	updateGeometry(temp.Polygon, edit.Poly, edit.CurrentGeoID)
	updateGeometry(temp.Text, edit.Text, edit.CurrentGeoID)
	updateGeometry(temp.List, edit.List, edit.CurrentGeoID)

}

func updateGeometry[T geometry.Geometer[S], S any](geos []T, edit S, geoID int) {
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
	clearBox(edit.geometryTab)
	clearBox(edit.visibleTab)

	updateEditor(edit, edit.Rect, temp.Rectangle)
	updateEditor(edit, edit.Circle, temp.Circle)
	updateEditor(edit, edit.Clock, temp.Clock)
	updateEditor(edit, edit.Image, temp.Image)
	updateEditor(edit, edit.Poly, temp.Polygon)
	updateEditor(edit, edit.Text, temp.Text)
	updateEditor(edit, edit.List, temp.List)

	return
}

func updateEditor[T geometry.Editor[S], S geometry.Geometer[T]](edit *Editor, editor T, geos []S) {
	for _, geo := range geos {
		if geo.GetGeometryID() != edit.CurrentGeoID {
			continue
		}

		editor.UpdateEditor(geo)
		edit.geometryTab.PackStart(editor.GetBox(), true, true, 0)
		//edit.tabs.AppendPage(editor.GetVisibleBox(), visibleLabel)

		edit.tabs.SetCurrentPage(0)
	}
}

func (edit *Editor) ClearFrame() {
	clearBox(edit.keyframeTab)
}

func (edit *Editor) SetFrame(frame SetFrame) {
	edit.SetFrameEdit.UpdateEditor(frame)
	edit.keyframeTab.PackStart(edit.SetFrameEdit.Scroll, true, true, 0)
	edit.tabs.SetCurrentPage(2)
}

func (edit *Editor) BindFrame(frame BindFrame) {
	err := edit.BindFrameEdit.UpdateEditor(frame)
	if err != nil {
		log.Printf("Error updating bind frame: %s", err)
	}

	edit.keyframeTab.PackStart(edit.BindFrameEdit.Scroll, true, true, 0)
	edit.tabs.SetCurrentPage(2)
}
