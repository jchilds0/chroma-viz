package gui

import (
	"github.com/gotk3/gotk3/gtk"
)

func LaunchGui() {
    gtk.Init(nil)

    win, _ := gtk.WindowNew(gtk.WINDOW_TOPLEVEL)
    win.SetTitle("Chroma Viz")
    win.Connect("destroy", func() { gtk.MainQuit() })

    box, _ := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 10)
    box.PackStart(menuLayout(), false, false, 0)
    box.PackStart(bodyLayout(), true, true, 0)
    box.PackEnd(lowerBarLayout(), false, false, 0)
    
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

func bodyLayout() *gtk.Box {
    box, _ := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 0)
    left := leftBox()
    right := rightBox()

    box.PackStart(left, true, true, 0)
    box.PackStart(right, true, true, 0)
    box.SetVAlign(gtk.ALIGN_CENTER)

    return box
}

func leftBox() *gtk.Box {
    box, _ := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0)
    templates := templatesBox()
    show := showBox()

    box.PackStart(templates, true, true, 0)
    box.PackStart(show, true, true, 0)

    return box
}

func templatesBox() *gtk.Box {
    box, _ := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0)
    header, _ := gtk.HeaderBarNew()
    header.SetTitle("Templates")

    templates, _ := gtk.ListBoxNew()

    box.PackStart(header, false, false, 0)
    box.PackStart(templates, true, true, 0)

    return box
}

func showBox() *gtk.Box {
    box, _ := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0)
    header, _ := gtk.HeaderBarNew()
    header.SetTitle("Show")

    show, _ := gtk.ListBoxNew()

    box.PackStart(header, false, false, 0)
    box.PackStart(show, true, true, 0)

    return box
}

func rightBox() *gtk.Box {
    box, _ := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0)

    return box
}

func lowerBarLayout() *gtk.HeaderBar {
    box, _ := gtk.HeaderBarNew()

    button, _ := gtk.ButtonNew()
    button.SetLabel("Exit")
    button.Connect("clicked", func() { gtk.MainQuit() })
    return box 
}

