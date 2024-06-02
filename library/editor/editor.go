package editor

import (
	"chroma-viz/library/pages"
	"chroma-viz/library/props"
	"chroma-viz/library/tcp"
	"fmt"

	"github.com/gotk3/gotk3/gtk"
)

var propList = []string{
	props.RECT_PROP, props.TEXT_PROP, props.CIRCLE_PROP,
	props.TICKER_PROP, props.CLOCK_PROP, props.IMAGE_PROP,
}

type Pairing struct {
	prop   *props.Property
	editor *props.PropertyEditor
}

type Editor struct {
	Box            *gtk.Box
	tabs           *gtk.Notebook
	actions        *gtk.Box
	CurrentPage    tcp.Animator
	propEditPairs  []Pairing
	propertyEditor map[string][]*props.PropertyEditor
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

/*
 * Init the editor to load pages, i.e. a
 * tab for each property in the current page
 */
func (editor *Editor) PageEditor() (err error) {
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

	// prop editors
	editor.propertyEditor = make(map[string][]*props.PropertyEditor, len(propList))
	initNumPropEdits := 10

	for _, p := range propList {
		editor.propertyEditor[p] = make([]*props.PropertyEditor, initNumPropEdits)

		for j := 0; j < initNumPropEdits; j++ {
			editor.propertyEditor[p][j], err = props.NewPropertyEditor(p)

			if err != nil {
				return
			}
		}
	}

	return
}

/*
 * Init the editor to load properties, i.e.
 * a tab with the property attributes and
 * a tab with the visibility editor
 */
func (editor *Editor) PropertyEditor() (err error) {
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

	// prop editors
	editor.propertyEditor = make(map[string][]*props.PropertyEditor, len(propList))
	initNumPropEdits := 10

	for _, p := range propList {
		editor.propertyEditor[p] = make([]*props.PropertyEditor, initNumPropEdits)

		for j := 0; j < initNumPropEdits; j++ {
			editor.propertyEditor[p][j], err = props.NewPropertyEditor(p)

			if err != nil {
				return
			}
		}
	}

	return
}

/*
 * Store the editor values in the properties
 */
func (edit *Editor) UpdateProps() {
	for _, item := range edit.propEditPairs {
		if item.prop == nil {
			continue
		}

		item.prop.UpdateProp(item.editor)
	}
}

/*
 * Load page 'page' into the editor
 */
func (editor *Editor) SetPage(page *pages.Page) (err error) {
	num_pages := editor.tabs.GetNPages()
	for i := 0; i < num_pages; i++ {
		editor.tabs.RemovePage(0)
	}

	editor.CurrentPage = page
	editor.propEditPairs = make([]Pairing, 0, 10)
	propCount := make(map[string]int, len(propList))

	var label *gtk.Label
	for _, prop := range editor.CurrentPage.GetPropMap() {
		if prop == nil {
			continue
		}

		typed := prop.PropType
		label, err = gtk.LabelNew(prop.Name)
		if err != nil {
			return
		}

		// pair up with prop editor
		var propEdit *props.PropertyEditor
		if propCount[typed] == len(editor.propertyEditor[typed]) {
			// we ran out of editors, add a new one
			propEdit, err = props.NewPropertyEditor(typed)
			if err != nil {
				return
			}

			editor.propertyEditor[typed] = append(editor.propertyEditor[typed], propEdit)
		} else {
			propEdit = editor.propertyEditor[typed][propCount[typed]]
		}

		propEdit.UpdateEditor(prop)
		editor.tabs.AppendPage(propEdit.Scroll, label)
		editor.propEditPairs = append(editor.propEditPairs, Pairing{prop: prop, editor: propEdit})

		propCount[typed]++
	}

	return
}

/*
 * Load property 'prop' into the editor
 */
func (editor *Editor) SetProperty(prop *props.Property) (err error) {
	num_pages := editor.tabs.GetNPages()
	for i := 0; i < num_pages; i++ {
		editor.tabs.RemovePage(0)
	}

	geoLabel, err := gtk.LabelNew("Geometry")
	if err != nil {
		return
	}

	editor.propEditPairs = nil
	propEdit := editor.propertyEditor[prop.PropType]
	if propEdit == nil || len(propEdit) <= 0 {
		err = fmt.Errorf("Prop edit %s is nil", prop.PropType)
		return
	}

	visibleBox, err := propEdit[0].CreateVisibleEditor()
	if err != nil {
		return
	}

	visibleLabel, err := gtk.LabelNew("Visible")
	if err != nil {
		return
	}

	propEdit[0].UpdateEditorAllProp(prop)
	editor.propEditPairs = []Pairing{{prop: prop, editor: propEdit[0]}}

	editor.tabs.AppendPage(propEdit[0].Scroll, geoLabel)
	editor.tabs.AppendPage(visibleBox, visibleLabel)
	return
}
