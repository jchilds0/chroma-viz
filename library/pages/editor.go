package pages

import (
	"chroma-viz/library/geometry"
	"chroma-viz/library/templates"
	"fmt"
	"log"

	"github.com/gotk3/gotk3/gtk"
)

type Editor struct {
	Box         *gtk.Box
	CurrentPage *Page
	actions     *gtk.Box

	tabCount    int
	tab         []*gtk.Frame
	tabContents map[int]gtk.IWidget
	notebook    *gtk.Notebook

	Rect   []*geometry.RectangleEditor
	Circle []*geometry.CircleEditor
	Clock  []*geometry.ClockEditor
	Image  []*geometry.ImageEditor
	Poly   []*geometry.PolygonEditor
	Text   []*geometry.TextEditor
	List   []*geometry.ListEditor
}

func NewEditor() (editor *Editor, err error) {
	editor = &Editor{
		tab:         make([]*gtk.Frame, 0, 128),
		tabContents: make(map[int]gtk.IWidget, 128),
	}

	editor.Box, err = gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0)
	if err != nil {
		return
	}

	editor.actions, err = gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 0)
	if err != nil {
		return
	}

	editor.Box.PackStart(editor.actions, false, false, 10)

	editor.notebook, err = gtk.NotebookNew()
	if err != nil {
		return
	}

	editor.notebook.SetScrollable(true)
	editor.notebook.Connect("focus-tab", editor.activateTab)

	editor.Box.PackStart(editor.notebook, true, true, 0)

	tab, _ := gtk.FrameNew("")
	editor.appendTab("Select A Page", tab)

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
	edit.CurrentPage.lock.Lock()
	defer edit.CurrentPage.lock.Unlock()

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

func (edit *Editor) activateTab() {
	pos := edit.notebook.GetCurrentPage()
	if pos >= len(edit.tab) || pos >= len(edit.tabContents) {
		log.Println("Tab", pos, "out of range")
		return
	}

	tab := edit.tab[pos]
	tab.SetVisible(true)

	tabChild, err := tab.GetChild()
	if err == nil {
		tab.Remove(tabChild)
	}

	contents, ok := edit.tabContents[pos]
	if !ok {
		log.Println("Missing tab contents", pos)
		return
	}

	tab.Add(contents)
}

func (edit *Editor) appendTab(label string, widget gtk.IWidget) (err error) {
	index := edit.tabCount

	if index < len(edit.tab) {
		tab := edit.tab[index]
		if tab == nil {
			return fmt.Errorf("Missing tab %d", index)
		}

		edit.notebook.SetTabLabelText(tab, label)
	} else {
		tab, err := gtk.FrameNew("")
		if err != nil {
			return err
		}

		tab.Add(widget)
		tab.SetVisible(true)

		edit.tab = append(edit.tab, tab)

		tabLabel, err := gtk.LabelNew(label)
		if err != nil {
			return err
		}

		edit.notebook.AppendPage(tab, tabLabel)
	}

	edit.tabContents[index] = widget
	edit.tabCount++
	return
}

// Load page 'page' into the editor
func (edit *Editor) SetPage(page *Page) (err error) {
	clear(edit.tabContents)
	edit.tabCount = 0
	edit.CurrentPage = page

	edit.Rect = updateEditor(edit, edit.Rect, edit.CurrentPage.Rect, geometry.NewRectangleEditor)
	edit.Circle = updateEditor(edit, edit.Circle, edit.CurrentPage.Circle, geometry.NewCircleEditor)
	edit.Clock = updateEditor(edit, edit.Clock, edit.CurrentPage.Clock, geometry.NewClockEditor)
	edit.Image = updateEditor(edit, edit.Image, edit.CurrentPage.Image, geometry.NewImageEditor)
	edit.Poly = updateEditor(edit, edit.Poly, edit.CurrentPage.Poly, geometry.NewPolygonEditor)
	edit.Text = updateEditor(edit, edit.Text, edit.CurrentPage.Text, geometry.NewTextEditor)
	edit.List = updateEditor(edit, edit.List, edit.CurrentPage.List, geometry.NewListEditor)

	edit.notebook.SetCurrentPage(0)
	edit.activateTab()
	return
}

func updateEditor[T geometry.Editor[S], S geometry.Geometer[T]](
	edit *Editor, editors []T, geos []S, init func() (T, error)) []T {
	diff := len(geos) - len(editors)
	if diff > 0 {
		for range diff {
			edit, err := init()
			if err != nil {
				log.Print(err)
				continue
			}

			editors = append(editors, edit)
		}
	}

	for i := range geos {
		if isNil(geos[i]) {
			continue
		}

		err := editors[i].UpdateEditor(geos[i])
		if err != nil {
			log.Println("Updating editor:", err)
			continue
		}

		edit.appendTab(geos[i].GetName(), editors[i].GetBox())
	}

	return editors
}

func (edit *Editor) UpdateTemplate(newTemp *templates.Template) {
	edit.CurrentPage = NewPageFromTemplate(newTemp)
	edit.UpdateProps()
	edit.SetPage(edit.CurrentPage)
}
