package gui

import (
	"chroma-viz/props"
	"fmt"
	"log"

	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
)

type Pairing struct {
    prop      props.Property
    editor    props.PropertyEditor
}

type Editor struct {
    box       *gtk.Box
    tabs      *gtk.Notebook
    header    *gtk.HeaderBar
    page      *Page
    animate   func()
    cont      func()
    pairs     []Pairing
    propEdit  [][]props.PropertyEditor
}

func NewEditor() *Editor {
    var err error
    editor := &Editor{}

    glib.SignalNew("animate-on")
    glib.SignalNew("continue")
    glib.SignalNew("animate-off")

    editor.box, err = gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0)
    if err != nil { 
        log.Fatalf("Error creating editor (%s)", err) 
    }

    editor.header, err = gtk.HeaderBarNew()
    if err != nil { 
        log.Fatalf("Error creating editor (%s)", err) 
    }

    editor.header.SetTitle("Editor")
    editor.box.PackStart(editor.header, false, false, 0)

    actions, err := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 0)
    if err != nil { 
        log.Fatalf("Error creating editor (%s)", err) 
    }

    editor.box.PackStart(actions, false, false, 10)

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
        if editor.page == nil { 
            fmt.Println("No page selected")
            return
        }

        editor.UpdateProps(ANIMATE_ON)
        conn["Engine"].setPage <- editor.page
        conn["Engine"].sendPage <- ANIMATE_ON
    })

    take2.Connect("clicked", func() { 
        if editor.page == nil { 
            fmt.Println("No page selected")
            return
        }

        _, ok := conn["Engine"]
        if ok == false {
            fmt.Println("Engine not found")
            return
        }

        editor.UpdateProps(CONTINUE)
        conn["Engine"].setPage <- editor.page
        conn["Engine"].sendPage <- CONTINUE 
    })

    take3.Connect("clicked", func() { 
        if editor.page == nil { 
            log.Printf("No page selected")
            return
        }

        editor.UpdateProps(ANIMATE_OFF)
        conn["Engine"].setPage <- editor.page
        conn["Engine"].sendPage <- ANIMATE_OFF
    })

    save.Connect("clicked", func() {
        if editor.page == nil {
            log.Printf("No page selected")
            return
        }
    })

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
    editor.box.PackStart(editor.tabs, true, true, 0)

    // prop editors 
    editor.propEdit = make([][]props.PropertyEditor, props.NUM_PROPS)

    editor.animate = func() { 
        editor.UpdateProps(ANIMATE_ON)
        conn["Preview"].setPage <- editor.page
        conn["Preview"].sendPage <- ANIMATE_ON
    }

    editor.cont = func() {
        editor.UpdateProps(CONTINUE)
        conn["Engine"].sendPage <- CONTINUE
        conn["Preview"].sendPage <- CONTINUE
    }
    
    for i := range editor.propEdit {
        num := 10
        editor.propEdit[i] = make([]props.PropertyEditor, num)

        for j := 0; j < num; j++ {
            editor.propEdit[i][j] = props.NewPropertyEditor(i, editor.animate, editor.cont)
        }
    }

    return editor
}

func (edit *Editor) UpdateProps(action int) {
    for _, item := range edit.pairs {
        if item.prop == nil {
            continue
        }

        item.prop.Update(item.editor, action)
    }
}

func (edit *Editor) SetPage(page *Page) {
    num_pages := edit.tabs.GetNPages()
    for i := 0; i < num_pages; i++  {
        edit.tabs.RemovePage(0)
    }

    edit.page = page
    edit.pairs = make([]Pairing, 0, 10)
    propCount := make([]int, props.NUM_PROPS)
    for _, prop := range edit.page.propMap {
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
        if propCount[typed] == len(edit.propEdit[typed]) {
            // we ran out of editors, add a new one
            propEdit = props.NewPropertyEditor(typed, edit.animate, edit.cont)
            edit.propEdit[typed] = append(edit.propEdit[typed], propEdit)
        } else {
            propEdit = edit.propEdit[typed][propCount[typed]]
        }

        propEdit.Update(prop)
        edit.tabs.AppendPage(propEdit.Box(), label)
        propCount[typed]++
        edit.pairs = append(edit.pairs, Pairing{prop: prop, editor: propEdit})
    }
}

func (edit *Editor) Box() *gtk.Box {
    return edit.box
}

