package editor

import (
	"chroma-viz/library/props"
	"chroma-viz/library/shows"
	"chroma-viz/library/tcp"
	"log"

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

func NewEditor(sendEngine, sendPreview func(tcp.Animator, int)) *Editor {
	var err error
	editor := &Editor{}

	editor.Box, err = gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0)
	if err != nil {
		log.Fatalf("Error creating editor (%s)", err)
	}

	editor.actions, err = gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 0)
	if err != nil {
		log.Fatalf("Error creating editor (%s)", err)
	}

	editor.Box.PackStart(editor.actions, false, false, 10)
	return editor
}

func (editor *Editor) AddAction(label string, start bool, action func()) {
	button, err := gtk.ButtonNewWithLabel(label)
	if err != nil {
		log.Fatalf("Error creating editor (%s)", err)
	}

	button.Connect("clicked", action)

	if start {
		editor.actions.PackStart(button, false, false, 10)
	} else {
		editor.actions.PackEnd(button, false, false, 10)
	}
}

func (editor *Editor) PageEditor() {
	var err error
	editor.tabs, err = gtk.NotebookNew()
	if err != nil {
		log.Fatalf("Error creating editor (%s)", err)
	}

	tab, err := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0)
	if err != nil {
		log.Fatalf("Error creating editor (%s)", err)
	}

	tabLabel, err := gtk.LabelNew("Select A Page")
	if err != nil {
		log.Fatalf("Error creating editor (%s)", err)
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
				log.Printf("Error creating prop editor %d", i)
			}
		}
	}
}

func (editor *Editor) PropertyEditor() {
	var err error
	editor.tabs, err = gtk.NotebookNew()
	if err != nil {
		log.Fatal(err)
	}

	tab, err := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0)
	if err != nil {
		log.Fatalf("Error creating editor (%s)", err)
	}

	tabLabel, err := gtk.LabelNew("Select A Page")
	if err != nil {
		log.Fatalf("Error creating editor (%s)", err)
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
				log.Printf("Error creating prop editor %d", i)
			}
		}
	}
}

func (edit *Editor) UpdateProps() {
	for _, item := range edit.pairs {
		if item.prop == nil {
			continue
		}

		item.prop.UpdateProp(item.editor)
	}
}

func (editor *Editor) SetPage(page *shows.Page) {
	num_pages := editor.tabs.GetNPages()
	for i := 0; i < num_pages; i++ {
		editor.tabs.RemovePage(0)
	}

	editor.Page = page
	editor.pairs = make([]Pairing, 0, 10)
	propCount := make([]int, props.NUM_PROPS)
	for _, prop := range editor.Page.GetPropMap() {
		if prop == nil {
			log.Print("Editor recieved nil prop")
			continue
		}

		typed := prop.PropType
		label, err := gtk.LabelNew(prop.Name)
		if err != nil {
			log.Printf("Error setting page (%s)", err)
			return
		}

		// pair up with prop editor
		var propEdit *props.PropertyEditor
		if propCount[typed] == len(editor.propEdit[typed]) {
			// we ran out of editors, add a new one
			propEdit, err = props.NewPropertyEditor(typed)
			if err != nil {
				log.Printf("Error creating prop editor %d", typed)
			}

			editor.propEdit[typed] = append(editor.propEdit[typed], propEdit)
		} else {
			propEdit = editor.propEdit[typed][propCount[typed]]
		}

		propEdit.UpdateEditor(prop)
		editor.tabs.AppendPage(propEdit.Box, label)
		propCount[typed]++
		editor.pairs = append(editor.pairs, Pairing{prop: prop, editor: propEdit})
	}
}

func (editor *Editor) SetProperty(prop *props.Property) {
	num_pages := editor.tabs.GetNPages()
	for i := 0; i < num_pages; i++ {
		editor.tabs.RemovePage(0)
	}

	geoLabel, err := gtk.LabelNew("Geometry")
	if err != nil {
		log.Printf("Error setting property (%s)", err)
		return
	}

	editor.pairs = nil
	propEdit := editor.propEdit[prop.PropType][0]

	visibleBox, err := propEdit.CreateVisibleEditor()
	if err != nil {
		log.Printf("Error setting property (%s)", err)
		return
	}

	visibleLabel, err := gtk.LabelNew("Visible")
	if err != nil {
		log.Printf("Error setting property (%s)", err)
		return
	}

	propEdit.UpdateEditorAllProp(prop)
	editor.pairs = []Pairing{{prop: prop, editor: propEdit}}

	editor.tabs.AppendPage(propEdit.Box, geoLabel)
	editor.tabs.AppendPage(visibleBox, visibleLabel)

	//editor.Box.PackStart(editor.propBox, true, true, 0)
}
