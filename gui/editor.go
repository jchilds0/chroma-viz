package gui

import (
	"fmt"

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
    editor.box.PackStart(actions, false, false, 0)

    take1, _ := gtk.ButtonNewWithLabel("Take On")
    take2, _ := gtk.ButtonNewWithLabel("Continue")
    take3, _ := gtk.ButtonNewWithLabel("Take Off")
    actions.PackStart(take1, false, false, 0)
    actions.PackStart(take2, false, false, 0)
    actions.PackStart(take3, false, false, 0)

    take1.Connect("clicked", func() { 
        if editor.page == nil { 
            fmt.Println("No page selected")
            return
        }

        editor.conn.SendPage(editor.page.templateID, ANIMATE_ON)
    })

    take2.Connect("clicked", func() { 
        if editor.page == nil { 
            fmt.Println("No page selected")
            return
        }

        editor.conn.SendPage(editor.page.templateID, CONTINUE)
    })

    take3.Connect("clicked", func() { 
        if editor.page == nil { 
            fmt.Println("No page selected")
            return
        }

        editor.conn.SendPage(editor.page.templateID, ANIMATE_OFF)
    })
    
    editor.edit, _ = gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0)
    editor.box.PackStart(editor.edit, true, true, 0)

    return editor
}

func (edit *Editor) SetPage(page *Page) {
    edit.header.SetTitle("Editor - " + page.title)
    edit.page = page

    edit.edit, _ = gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0)
    edit.edit.Add(IntEditor("x Pos: ", 100))
    edit.edit.PackStart(IntEditor("y Pos: ", 100), false, false, 0)

    edit.box.QueueDraw()
}

func (edit *Editor) Box() *gtk.Box {
    return edit.box
}

func IntEditor(name string, max int) *gtk.Box {
    box, _ := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 0)

    label, _ := gtk.LabelNew(name)
    box.PackStart(label, false, false, 0)

    spin, _ := gtk.SpinButtonNewWithRange(0, float64(max), 1)
    box.PackStart(spin, false, false, 0)

    return box
}
