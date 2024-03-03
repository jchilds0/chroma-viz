package editor

import (
	"chroma-viz/props"
	"chroma-viz/shows"
	"chroma-viz/tcp"
	"fmt"
	"log"

	"github.com/gotk3/gotk3/gtk"
)

type Pairing struct {
    prop      props.Property
    editor    props.PropertyEditor
}

type Editor struct {
    Box       *gtk.Box
    tabs      *gtk.Notebook
    header    *gtk.HeaderBar
    propBox   *gtk.Box
    Page      *shows.Page
    pairs     []Pairing
    propEdit  [][]props.PropertyEditor
    sendEngine func(*shows.Page, int)
    sendPreview func(*shows.Page, int)
}

func NewEditor(sendEngine, sendPreview func(*shows.Page, int)) *Editor {
    var err error
    editor := &Editor{
        sendEngine: sendEngine,
        sendPreview: sendPreview,
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

    return editor
}

func (editor *Editor) EnginePanel() {
    var err error

    actions, err := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 0)
    if err != nil { 
        log.Fatalf("Error creating editor (%s)", err) 
    }

    editor.Box.PackStart(actions, false, false, 10)

    take1, err := gtk.ButtonNewWithLabel("Take On")
    if err != nil { 
        log.Fatalf("Error creating editor (%s)", err) 
    }

    take2, err := gtk.ButtonNewWithLabel("Continue")
    if err != nil { 
        log.Fatalf("Error creating editor (%s)", err) 
    }

    take3, err := gtk.ButtonNewWithLabel("Take Off")
    if err != nil { 
        log.Fatalf("Error creating editor (%s)", err) 
    }

    save, err := gtk.ButtonNewWithLabel("Save")
    if err != nil { 
        log.Fatalf("Error creating editor (%s)", err) 
    }

    actions.PackStart(take1, false, false, 10)
    actions.PackStart(take2, false, false, 0)
    actions.PackStart(take3, false, false, 10)
    actions.PackEnd(save, false, false, 10)

    take1.Connect("clicked", func() { 
        if editor.Page == nil { 
            fmt.Println("No page selected")
            return
        }

        editor.UpdateProps(tcp.ANIMATE_ON)
        editor.sendEngine(editor.Page, tcp.ANIMATE_ON)
    })

    take2.Connect("clicked", func() { 
        if editor.Page == nil { 
            fmt.Println("No page selected")
            return
        }

        editor.UpdateProps(tcp.CONTINUE)
        editor.sendEngine(editor.Page, tcp.CONTINUE)
    })

    take3.Connect("clicked", func() { 
        if editor.Page == nil { 
            log.Printf("No page selected")
            return
        }

        editor.UpdateProps(props.ANIMATE_OFF)
        editor.sendEngine(editor.Page, tcp.ANIMATE_OFF)
    })

    save.Connect("clicked", func() {
        if editor.Page == nil {
            log.Printf("No page selected")
            return
        }

        editor.sendPreview(editor.Page, tcp.ANIMATE_ON)
    })
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
    editor.propEdit = make([][]props.PropertyEditor, props.NUM_PROPS)

    animate := func() { 
        editor.UpdateProps(tcp.ANIMATE_ON)
        editor.sendPreview(editor.Page, tcp.ANIMATE_ON)
    }

    cont := func() {
        editor.UpdateProps(tcp.CONTINUE)
        editor.sendPreview(editor.Page, tcp.CONTINUE)
        editor.sendEngine(editor.Page, tcp.CONTINUE)
    }
    
    for i := range editor.propEdit {
        num := 10
        editor.propEdit[i] = make([]props.PropertyEditor, num)

        for j := 0; j < num; j++ {
            editor.propEdit[i][j] = props.NewPropertyEditor(i, animate, cont)
        }
    }
}

func (editor *Editor) PropertyEditor() {
    // prop editors 
    editor.propEdit = make([][]props.PropertyEditor, props.NUM_PROPS)

    animate := func() { 
        editor.UpdateProps(tcp.ANIMATE_ON)
        editor.sendPreview(editor.Page, tcp.ANIMATE_ON)
    }

    cont := func() {
        editor.UpdateProps(tcp.CONTINUE)
        editor.sendPreview(editor.Page, tcp.CONTINUE)
        editor.sendEngine(editor.Page, tcp.CONTINUE)
    }
    
    for i := range editor.propEdit {
        num := 1
        editor.propEdit[i] = make([]props.PropertyEditor, num)

        for j := 0; j < num; j++ {
            editor.propEdit[i][j] = props.NewPropertyEditor(i, animate, cont)
        }
    }
}

func (edit *Editor) UpdateProps(action int) {
    for _, item := range edit.pairs {
        if item.prop == nil {
            continue
        }

        item.prop.Update(item.editor, action)
    }
}

func (editor *Editor) SetPage(page *shows.Page) {
    num_pages := editor.tabs.GetNPages()
    for i := 0; i < num_pages; i++  {
        editor.tabs.RemovePage(0)
    }

    animate := func() { 
        editor.UpdateProps(tcp.ANIMATE_ON)
        editor.sendPreview(editor.Page, tcp.ANIMATE_ON)
    }

    cont := func() {
        editor.UpdateProps(tcp.CONTINUE)
        editor.sendPreview(editor.Page, tcp.CONTINUE)
        editor.sendEngine(editor.Page, tcp.CONTINUE)
    }

    editor.Page = page
    editor.pairs = make([]Pairing, 0, 10)
    propCount := make([]int, props.NUM_PROPS)
    for _, prop := range editor.Page.PropMap {
        if prop == nil {
            log.Print("Editor recieved nil prop")
            continue
        }

        typed := prop.Type()
        label, err := gtk.LabelNew(prop.Name())
        if err != nil { 
            log.Fatalf("Error setting page (%s)", err) 
        }

        // pair up with prop editor
        var propEdit props.PropertyEditor
        if propCount[typed] == len(editor.propEdit[typed]) {
            // we ran out of editors, add a new one
            propEdit = props.NewPropertyEditor(typed, animate, cont)
            editor.propEdit[typed] = append(editor.propEdit[typed], propEdit)
        } else {
            propEdit = editor.propEdit[typed][propCount[typed]]
        }

        propEdit.Update(prop)
        editor.tabs.AppendPage(propEdit.Box(), label)
        propCount[typed]++
        editor.pairs = append(editor.pairs, Pairing{prop: prop, editor: propEdit})
    }
}

func (editor *Editor) SetProperty(prop props.Property) {
    if editor.propBox != nil {
        editor.Box.Remove(editor.propBox)
    }

    propType := prop.Type()
    propEdit := editor.propEdit[propType][0]
    propEdit.Update(prop)
    editor.propBox = propEdit.Box()
    editor.pairs = []Pairing{{prop: prop, editor: propEdit}}

    editor.Box.PackStart(editor.propBox, true, true, 0)
}

