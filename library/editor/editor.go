package editor

import (
	"chroma-viz/library/pages"
	"chroma-viz/library/props"
	"chroma-viz/library/tcp"
	"fmt"

	"github.com/gotk3/gotk3/gtk"
)

type Pairing struct {
	prop   *props.Property
	editor *props.PropertyEditor
}

type Editor struct {
	Box      *gtk.Box
	tabs     *gtk.Notebook
	header   *gtk.HeaderBar
	actions  *gtk.Box
	propBox  *gtk.Box
	Page     tcp.Animator
	pairs    []Pairing
	propEdit [][]*props.PropertyEditor
}

func NewEditor(sendEngine, sendPreview func(tcp.Animator, int)) (editor *Editor, err error) {
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
	editor.propEdit = make([][]*props.PropertyEditor, props.NUM_PROPS)

	for i := range editor.propEdit {
		num := 10
		editor.propEdit[i] = make([]*props.PropertyEditor, num)

		for j := 0; j < num; j++ {
			editor.propEdit[i][j], err = props.NewPropertyEditor(i)
			if err != nil {
				return
			}
		}
	}

	return
}

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
	editor.propEdit = make([][]*props.PropertyEditor, props.NUM_PROPS)

	for i := range editor.propEdit {
		num := 1
		editor.propEdit[i] = make([]*props.PropertyEditor, num)

		for j := 0; j < num; j++ {
			editor.propEdit[i][j], err = props.NewPropertyEditor(i)
			if err != nil {
				return
			}
		}
	}

	return
}

func (edit *Editor) UpdateProps() {
	for _, item := range edit.pairs {
		if item.prop == nil {
			continue
		}

		item.prop.UpdateProp(item.editor)
	}
}

func (editor *Editor) SetPage(page *pages.Page) (err error) {
	num_pages := editor.tabs.GetNPages()
	for i := 0; i < num_pages; i++ {
		editor.tabs.RemovePage(0)
	}

	editor.Page = page
	editor.pairs = make([]Pairing, 0, 10)
	propCount := make([]int, props.NUM_PROPS)

	var label *gtk.Label

	for _, prop := range editor.Page.GetPropMap() {
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

		if propCount[typed] == len(editor.propEdit[typed]) {
			// we ran out of editors, add a new one
			propEdit, err = props.NewPropertyEditor(typed)
			if err != nil {
				return
			}

			editor.propEdit[typed] = append(editor.propEdit[typed], propEdit)
		} else {
			propEdit = editor.propEdit[typed][propCount[typed]]
		}

		propEdit.UpdateEditor(prop)
		editor.tabs.AppendPage(propEdit.Scroll, label)
		propCount[typed]++
		editor.pairs = append(editor.pairs, Pairing{prop: prop, editor: propEdit})
	}

	return
}

func (editor *Editor) SetProperty(prop *props.Property) (err error) {
	num_pages := editor.tabs.GetNPages()
	for i := 0; i < num_pages; i++ {
		editor.tabs.RemovePage(0)
	}

	geoLabel, err := gtk.LabelNew("Geometry")
	if err != nil {
		return
	}

	editor.pairs = nil
	propEdit := editor.propEdit[prop.PropType]
	if propEdit == nil || len(propEdit) <= 0 {
		err = fmt.Errorf("Prop edit %s is nil", props.PropType(prop.PropType))
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
	editor.pairs = []Pairing{{prop: prop, editor: propEdit[0]}}

	editor.tabs.AppendPage(propEdit[0].Scroll, geoLabel)
	editor.tabs.AppendPage(visibleBox, visibleLabel)
	return
}
