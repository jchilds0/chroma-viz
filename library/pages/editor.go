package pages

import (
	"chroma-viz/library/geometry"
	"chroma-viz/library/templates"
	"log"

	"github.com/gotk3/gotk3/gtk"
)

type Editor struct {
	Box         *gtk.Box
	tabs        *gtk.Notebook
	actions     *gtk.Box
	CurrentPage *Page

	Rect   []*geometry.RectangleEditor
	Circle []*geometry.CircleEditor
	Clock  []*geometry.ClockEditor
	Image  []*geometry.ImageEditor
	Poly   []*geometry.PolygonEditor
	Text   []*geometry.TextEditor
	List   []*geometry.ListEditor
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

	editor.tabs.SetScrollable(true)

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

	numEditors := 10
	editor.Rect = initEditors(numEditors, geometry.NewRectangleEditor)
	editor.Circle = initEditors(numEditors, geometry.NewCircleEditor)
	editor.Clock = initEditors(numEditors, geometry.NewClockEditor)
	editor.Image = initEditors(numEditors, geometry.NewImageEditor)
	editor.Poly = initEditors(numEditors, geometry.NewPolygonEditor)
	editor.Text = initEditors(numEditors, geometry.NewTextEditor)
	editor.List = initEditors(numEditors, geometry.NewListEditor)

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

func initEditors[T any](numEditors int, init func() (T, error)) []T {
	editors := make([]T, numEditors)

	var err error
	for i := range numEditors {
		editors[i], err = init()
		if err != nil {
			log.Print(err)
		}
	}

	return editors
}

// Store the editor values in the properties
func (edit *Editor) UpdateProps() {
	updateGeometry(edit.CurrentPage.Rect, edit.Rect)
	updateGeometry(edit.CurrentPage.Circle, edit.Circle)
	updateGeometry(edit.CurrentPage.Clock, edit.Clock)
	updateGeometry(edit.CurrentPage.Image, edit.Image)
	updateGeometry(edit.CurrentPage.Poly, edit.Poly)
	updateGeometry(edit.CurrentPage.Text, edit.Text)
	updateGeometry(edit.CurrentPage.List, edit.List)
}

func updateGeometry[T geometry.Geometer[S], S any](geos []T, editors []S) {
	for i := range geos {
		err := geos[i].UpdateGeometry(editors[i])
		if err != nil {
			log.Print(err)
		}
	}
}

// Load page 'page' into the editor
func (edit *Editor) SetPage(page *Page) (err error) {
	num_pages := edit.tabs.GetNPages()
	for i := 0; i < num_pages; i++ {
		edit.tabs.RemovePage(0)
	}

	edit.CurrentPage = page

	updateEditor(edit, edit.Rect, edit.CurrentPage.Rect, geometry.NewRectangleEditor)
	updateEditor(edit, edit.Circle, edit.CurrentPage.Circle, geometry.NewCircleEditor)
	updateEditor(edit, edit.Clock, edit.CurrentPage.Clock, geometry.NewClockEditor)
	updateEditor(edit, edit.Image, edit.CurrentPage.Image, geometry.NewImageEditor)
	updateEditor(edit, edit.Poly, edit.CurrentPage.Poly, geometry.NewPolygonEditor)
	updateEditor(edit, edit.Text, edit.CurrentPage.Text, geometry.NewTextEditor)
	updateEditor(edit, edit.List, edit.CurrentPage.List, geometry.NewListEditor)

	return
}

func updateEditor[T geometry.Editor[S], S geometry.Geometer[T]](
	edit *Editor, editors []T, geos []S, init func() (T, error)) {
	diff := len(geos) - len(editors)
	if diff > 0 {
		for _ = range diff {
			edit, err := init()
			if err != nil {
				log.Print(err)
				continue
			}

			editors = append(editors, edit)
		}
	}

	var label *gtk.Label
	for i := range geos {
		if isNil(geos[i]) {
			continue
		}

		err := editors[i].UpdateEditor(geos[i])
		if err != nil {
			log.Print(err)
		}

		label, err = gtk.LabelNew(geos[i].GetName())
		if err != nil {
			return
		}

		edit.tabs.AppendPage(editors[i].GetBox(), label)
	}
}

func (edit *Editor) UpdateTemplate(newTemp *templates.Template) {
	edit.CurrentPage = NewPageFromTemplate(newTemp)
	edit.UpdateProps()
	edit.SetPage(edit.CurrentPage)
}
