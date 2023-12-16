package gui

import (
	"fmt"
	"log"

	"github.com/gotk3/gotk3/gtk"
)

type Editor struct {
    box       *gtk.Box
    tabs      *gtk.Notebook
    header    *gtk.HeaderBar
    page      *Page
}

func NewEditor() *Editor {
    var err error
    editor := &Editor{}

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

    actions.PackStart(take1, false, false, 10)
    actions.PackStart(take2, false, false, 0)
    actions.PackStart(take3, false, false, 10)

    take1.Connect("clicked", func() { 
        if editor.page == nil { 
            fmt.Println("No page selected")
            return
        }

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

        conn["Engine"].sendPage <- CONTINUE 
    })

    take3.Connect("clicked", func() { 
        if editor.page == nil { 
            fmt.Println("No page selected")
            return
        }

        conn["Engine"].sendPage <- ANIMATE_OFF
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
    return editor
}

func (edit *Editor) SetPage(page *Page) {
    num_pages := edit.tabs.GetNPages()
    for i := 0; i < num_pages; i++  {
        edit.tabs.RemovePage(0)
    }

    edit.page = page

    for name, key := range edit.page.propMap {
        label, err := gtk.LabelNew(name)
        if err != nil { 
            log.Fatalf("Error setting page (%s)", err) 
        }

        edit.tabs.AppendPage(key.Tab(), label)
    }
}

func (edit *Editor) Box() *gtk.Box {
    return edit.box
}

