package gui

import (
	"fmt"

	"github.com/gotk3/gotk3/gtk"
)

type Editor struct {
    box       *gtk.Box
    tabs      *gtk.Notebook
    header    *gtk.HeaderBar
    page      *Page
}

func NewEditor() *Editor {
    editor := &Editor{}
    editor.box, _ = gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0)

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

        conn["Engine"].SendPage(editor.page, ANIMATE_ON)
        conn["Engine"].Read()
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

        conn["Engine"].SendPage(editor.page, CONTINUE)
    })

    take3.Connect("clicked", func() { 
        if editor.page == nil { 
            fmt.Println("No page selected")
            return
        }

        conn["Engine"].SendPage(editor.page, ANIMATE_OFF)
    })

    editor.tabs, _ = gtk.NotebookNew()
    tab, _ := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0)
    tabLabel, _ := gtk.LabelNew("Select A Page")

    editor.tabs.AppendPage(tab, tabLabel)

    editor.box.PackStart(editor.tabs, true, true, 0)

    return editor
}

func (edit *Editor) SetPage(page *Page) {
    num_pages := edit.tabs.GetNPages()
    for i := 0; i < num_pages; i++  {
        edit.tabs.RemovePage(0)
    }

    edit.page = page

    for name, key := range edit.page.propMap {
        label, _ := gtk.LabelNew(name)
        edit.tabs.AppendPage(key.Tab(), label)
    }
}

func (edit *Editor) Box() *gtk.Box {
    return edit.box
}

