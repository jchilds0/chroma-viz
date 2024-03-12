package editor

import (
	"chroma-viz/library/props"
	"chroma-viz/library/shows"
	"log"

	"github.com/gotk3/gotk3/gtk"
)

type Pairing struct {
    prop      *props.Property
    editor    *props.PropertyEditor
}

type Editor struct {
    Box       *gtk.Box
    tabs      *gtk.Notebook
    header    *gtk.HeaderBar
    actions   *gtk.Box
    propBox   *gtk.Box
    Page      *shows.Page
    pairs     []Pairing
    propEdit  [][]*props.PropertyEditor
}

func NewEditor(sendEngine, sendPreview func(*shows.Page, int)) *Editor {
    var err error
    editor := &Editor{
    }

    editor.Box, err = gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0)
    if err != nil { 
        log.Fatalf("Error creating editor (%s)", err) 
    }

    editor.header, err = gtk.HeaderBarNew()
    if err != nil { 
        log.Fatalf("Error creating editor (%s)", err) 
    }

    editor.header.SetTitle("Editor")
    editor.Box.PackStart(editor.header, false, false, 0)

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
    for i := 0; i < num_pages; i++  {
        editor.tabs.RemovePage(0)
    }

    editor.Page = page
    editor.pairs = make([]Pairing, 0, 10)
    propCount := make([]int, props.NUM_PROPS)
    for _, prop := range editor.Page.PropMap {
        if prop == nil {
            log.Print("Editor recieved nil prop")
            continue
        }

        typed := prop.PropType
        label, err := gtk.LabelNew(prop.Name)
        if err != nil { 
            log.Fatalf("Error setting page (%s)", err) 
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
    if editor.propBox != nil {
        editor.Box.Remove(editor.propBox)
    }

    editor.pairs = nil
    propEdit := editor.propEdit[prop.PropType][0]

    propEdit.UpdateEditor(prop)
    editor.propBox = propEdit.Box
    editor.pairs = []Pairing{{prop: prop, editor: propEdit}}

    editor.Box.PackStart(editor.propBox, true, true, 0)
}

