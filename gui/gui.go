package gui

import (
	"github.com/gotk3/gotk3/gtk"
)

var conn map[string]*Connection

// TODO: seperate preview setup from gui setup
func LaunchGui(mainConn map[string]*Connection) {
    conn = mainConn 
    gtk.Init(nil)

    win, _ := gtk.WindowNew(gtk.WINDOW_TOPLEVEL)
    win.SetTitle("Chroma Viz")
    win.Connect("destroy", func() { 
        gtk.MainQuit() 
    })

    editView := NewEditor()
    showView := NewShow(editView)
    tempView := NewTempList(showView)

    showView.ImportShow(tempView, "/home/josh/Documents/projects/chroma-viz/shows/testing.show")

    box, _ := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0)

    /* Menu layout */
    menuBar, _ := gtk.MenuBarNew()
    box.PackStart(menuBar, false, false, 0)

    fileMenu, _ := gtk.MenuItemNewWithMnemonic("File")
    menuBar.Append(fileMenu)
    fileSubMenu, _ := gtk.MenuNew()
    fileMenu.SetSubmenu(fileSubMenu)

    newShow, _ := gtk.MenuItemNewWithLabel("New Show")
    fileSubMenu.Append(newShow)

    openShow, _ := gtk.MenuItemNewWithLabel("Open Show")
    fileSubMenu.Append(openShow)

    openShow.Connect("activate", func() {
        dialog, _ := gtk.FileChooserDialogNewWith2Buttons(
            "Import Show", win, gtk.FILE_CHOOSER_ACTION_OPEN, 
            "_Cancel", gtk.RESPONSE_CANCEL, "_Open", gtk.RESPONSE_ACCEPT)

        res := dialog.Run()

        if res == gtk.RESPONSE_ACCEPT {
            filename := dialog.GetFilename()
            showView.ImportShow(tempView, filename)
        }
        dialog.Destroy()
    })

    saveShow, _ := gtk.MenuItemNewWithLabel("Save Show")
    fileSubMenu.Append(saveShow)

    saveShow.Connect("activate", func() {
        dialog, _ := gtk.FileChooserDialogNewWith2Buttons(
            "Save Show", win, gtk.FILE_CHOOSER_ACTION_SAVE, 
            "_Cancel", gtk.RESPONSE_CANCEL, "_Save", gtk.RESPONSE_ACCEPT)

        dialog.SetCurrentName(".show")
        res := dialog.Run()
        if res == gtk.RESPONSE_ACCEPT {
            filename := dialog.GetFilename()
            showView.ExportShow(filename)
        }

        dialog.Destroy()
    })

    editMenu, _ := gtk.MenuItemNewWithMnemonic("Edit")
    menuBar.Append(editMenu)
    editSubMenu, _ := gtk.MenuNew()
    editMenu.SetSubmenu(editSubMenu)

    /* Body layout */
    bodyBox, _ := gtk.PanedNew(gtk.ORIENTATION_HORIZONTAL)
    box.PackStart(bodyBox, true, true, 0)


    leftBox, _ := gtk.PanedNew(gtk.ORIENTATION_VERTICAL)
    rightBox, _ := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0)
    bodyBox.Pack1(leftBox, true, true)
    bodyBox.Pack2(rightBox, true, true)

    /* left */
    leftBox.SetHExpand(true)

    /* templates */
    templates, _ := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0)
    leftBox.Pack1(templates, true, true)

    header1, _ := gtk.HeaderBarNew()
    header1.SetTitle("Templates")
    templates.PackStart(header1, false, false, 0)
    scroll1, _ := gtk.ScrolledWindowNew(nil, nil)
    templates.PackStart(scroll1, true, true, 0)
    scroll1.Add(tempView)

    /* show */
    shows, _ := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0)
    leftBox.Pack2(shows, true, true)

    header2, _ := gtk.HeaderBarNew()
    header2.SetTitle("Show")
    shows.PackStart(header2, false, false, 0)
    scroll2, _ := gtk.ScrolledWindowNew(nil, nil)
    shows.PackStart(scroll2, true, true, 0)
    scroll2.Add(showView)

    /* right */
    preview := setup_preview_window()
    rightBox.PackStart(editView.Box(), true, true, 0)
    rightBox.PackEnd(preview, false, false, 0)

    /* Lower Bar layout */
    lowerBox, _ := gtk.ActionBarNew()
    box.PackEnd(lowerBox, false, false, 0)

    button, _ := gtk.ButtonNew()
    lowerBox.PackEnd(button)

    button.SetLabel("Exit")
    button.Connect("clicked", func() { gtk.MainQuit() })

    for name, render := range conn {
        eng := NewEngineWidget(name, render)
        lowerBox.PackStart(eng)
    }

    win.Add(box)
    win.SetDefaultSize(800, 600)
    win.ShowAll()

    gtk.Main()
}

