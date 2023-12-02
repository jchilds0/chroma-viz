package gui

import (
	"fmt"
	"log"

	"github.com/gotk3/gotk3/gtk"
)

type Editor struct {
    box       *gtk.Box
    edit      *gtk.Box
    header    *gtk.HeaderBar
    page      *Page
    conn      *Connection 
}

func NewEditor(conn *Connection) *Editor {
    editor := &Editor{}
    editor.box, _ = gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0)
    editor.conn = conn

    editor.header, _ = gtk.HeaderBarNew()
    editor.header.SetTitle("Editor")
    editor.box.PackStart(editor.header, false, false, 0)

    actions, _ := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 0)
    editor.box.PackStart(actions, false, false, 10)

    take1, _ := gtk.ButtonNewWithLabel("Take On")
    take2, _ := gtk.ButtonNewWithLabel("Continue")
    take3, _ := gtk.ButtonNewWithLabel("Take Off")
    actions.PackStart(take1, false, false, 10)
    actions.PackStart(take2, false, false, 0)
    actions.PackStart(take3, false, false, 10)

    take1.Connect("clicked", func() { 
        if editor.page == nil { 
            fmt.Println("No page selected")
            return
        }

        editor.conn.SendPage(editor.page, ANIMATE_ON)
        editor.conn.Read()
    })

    take2.Connect("clicked", func() { 
        if editor.page == nil { 
            fmt.Println("No page selected")
            return
        }

        editor.conn.SendPage(editor.page, CONTINUE)
    })

    take3.Connect("clicked", func() { 
        if editor.page == nil { 
            fmt.Println("No page selected")
            return
        }

        editor.conn.SendPage(editor.page, ANIMATE_OFF)
    })
    
    editor.edit, _ = gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0)
    editor.box.PackStart(editor.edit, true, true, 0)

    return editor
}

func (edit *Editor) SetPage(page *Page) {
    edit.page = page

    edit.edit.Destroy()
    edit.edit, _ = gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 10)
    edit.edit.SetVisible(true)
    edit.box.PackStart(edit.edit, true, true, 0)

    label, _ := gtk.LabelNew(page.title)
    label.SetVisible(true)
    edit.edit.PackStart(label, false, false, 0)

    for _, key := range page.propList {
        prop, ok := page.props[key]

        if !ok {
            log.Printf("Unknown prop %s", key)
            return
        }

        edit.edit.PackStart(prop.Editor(key), false, false, 0)
    }

}

func (edit *Editor) Box() *gtk.Box {
    return edit.box
}

func IntEditor(name string, value *int, max int) *gtk.Box {
    padding := 10
    box, _ := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 0)
    box.SetVisible(true)

    label, _ := gtk.LabelNew(name)
    label.SetVisible(true)
    box.PackStart(label, false, false, uint(padding))

    spin, _ := gtk.SpinButtonNewWithRange(0, float64(max), 1)
    spin.SetVisible(true)
    spin.SetValue(float64(*value))
    spin.Connect("value-changed", func(spin *gtk.SpinButton) { *value = spin.GetValueAsInt() } )
    box.PackStart(spin, false, false, 0)

    return box
}
