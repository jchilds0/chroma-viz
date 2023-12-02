package gui

import (
	"github.com/gotk3/gotk3/gtk"
)

func LaunchGui(conn *Connection) {
    gtk.Init(nil)

    win, _ := gtk.WindowNew(gtk.WINDOW_TOPLEVEL)
    win.SetTitle("Chroma Viz")
    win.Connect("destroy", func() { gtk.MainQuit() })

    box, _ := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0)

    /* Menu layout */
    menuBox, _ := gtk.HeaderBarNew()
    menuBar, _ := gtk.MenuBarNew()
    menuBox.PackStart(menuBar)

    fileMenu, _ := gtk.MenuItemNewWithMnemonic("File")
    menuBar.Append(fileMenu)
    fileSubMenu, _ := gtk.MenuNew()
    fileMenu.SetSubmenu(fileSubMenu)

    editMenu, _ := gtk.MenuItemNewWithMnemonic("Edit")
    menuBar.Append(editMenu)
    editSubMenu, _ := gtk.MenuNew()
    editMenu.SetSubmenu(editSubMenu)
    box.PackStart(menuBox, false, false, 0)
 
    /* Body layout */
    bodyBox, _ := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 0)
    editView := NewEditor(conn)
    showView := NewShow(editView)
    tempView := NewTempList(showView)

    /* left */
    leftBox, _ := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0)
    leftBox.SetHExpand(true)

    header1, _ := gtk.HeaderBarNew()
    header1.SetTitle("Templates")
    leftBox.PackStart(header1, false, false, 0)

    scroll1, _ := gtk.ScrolledWindowNew(nil, nil)
    leftBox.PackStart(scroll1, true, true, 0)
    scroll1.Add(tempView)

    header2, _ := gtk.HeaderBarNew()
    header2.SetTitle("Show")
    leftBox.PackStart(header2, false, false, 0)

    scroll2, _ := gtk.ScrolledWindowNew(nil, nil)
    leftBox.PackStart(scroll2, true, true, 0)
    scroll2.Add(showView)

    /* right */
    rightBox, _ := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0)
    rightBox.PackStart(editView.Box(), true, true, 0)

    bodyBox.PackStart(leftBox, true, true, 0)
    bodyBox.PackStart(rightBox, true, true, 0)

    box.PackStart(bodyBox, true, true, 0)

    /* Lower Bar layout */
    lowerBox, _ := gtk.ActionBarNew()

    button, _ := gtk.ButtonNew()
    button.SetLabel("Exit")
    button.Connect("clicked", func() { gtk.MainQuit() })
    lowerBox.PackEnd(button)

    eng1 := NewEngineWidget(conn)
    lowerBox.PackStart(eng1)
    box.PackEnd(lowerBox, false, false, 0)
    
    win.Add(box)
    win.SetDefaultSize(800, 600)
    win.ShowAll()
    gtk.Main()
}

