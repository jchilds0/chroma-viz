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
    box.PackStart(menuLayout(), false, false, 0)
    box.PackStart(bodyLayout(conn), true, true, 0)
    box.PackEnd(lowerBarLayout(conn), false, false, 0)
    
    win.Add(box)
    win.SetDefaultSize(800, 600)
    win.ShowAll()
    gtk.Main()
}

func menuLayout() *gtk.HeaderBar {
    box, _ := gtk.HeaderBarNew()
    menuBar, _ := gtk.MenuBarNew()
    box.PackStart(menuBar)

    fileMenu, _ := gtk.MenuItemNewWithMnemonic("File")
    menuBar.Append(fileMenu)
    fileSubMenu, _ := gtk.MenuNew()
    fileMenu.SetSubmenu(fileSubMenu)

    editMenu, _ := gtk.MenuItemNewWithMnemonic("Edit")
    menuBar.Append(editMenu)
    editSubMenu, _ := gtk.MenuNew()
    editMenu.SetSubmenu(editSubMenu)
    return box
}

func bodyLayout(conn *Connection) *gtk.Box {
    box, _ := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 0)
    editView := NewEditor(conn)
    showView := NewShow(editView)
    tempView := NewTempList(showView)

    left := leftBox(showView, tempView)
    right := rightBox(editView)

    box.PackStart(left, true, true, 0)
    box.PackStart(right, true, true, 0)

    return box
}

func leftBox(showView *ShowTree, tempView *TempTree) *gtk.Box {
    box, _ := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0)
    box.SetHExpand(true)

    header1, _ := gtk.HeaderBarNew()
    header1.SetTitle("Templates")
    box.PackStart(header1, false, false, 0)

    scroll1, _ := gtk.ScrolledWindowNew(nil, nil)
    box.PackStart(scroll1, true, true, 0)
    scroll1.Add(tempView)

    header2, _ := gtk.HeaderBarNew()
    header2.SetTitle("Show")
    box.PackStart(header2, false, false, 0)

    scroll2, _ := gtk.ScrolledWindowNew(nil, nil)
    box.PackStart(scroll2, true, true, 0)
    scroll2.Add(showView)

    return box
}

func rightBox(editView *Editor) *gtk.Box {
    box, _ := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0)
    box.PackStart(editView.Box(), true, true, 0)

    return box
}

func lowerBarLayout(conn *Connection) *gtk.ActionBar {
    box, _ := gtk.ActionBarNew()

    button, _ := gtk.ButtonNew()
    button.SetLabel("Exit")
    button.Connect("clicked", func() { gtk.MainQuit() })
    box.PackEnd(button)

    eng1 := NewEngineWidget(conn)
    box.PackStart(eng1)

    return box 
}

