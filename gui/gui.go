package gui

import (
	"log"

	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
)

var conn map[string]*Connection

func InitConnections(){
    conn = make(map[string]*Connection)
}

func AddConnection(name string, ip string, port int) {
    conn[name] = NewConnection(ip, port)
}

// TODO: seperate preview setup from gui setup
func SetupMainGui(app *gtk.Application) {

    win, err := gtk.ApplicationWindowNew(app)
    if err != nil {
        log.Fatal(err)
    }

    win.SetDefaultSize(800, 600)
    win.SetTitle("Chroma Viz")

    editView := NewEditor()
    showView := NewShow(editView)
    tempView := NewTempList(showView)

    showView.ImportShow(tempView, "/home/josh/Documents/projects/chroma-viz/shows/testing.show")

    box, err := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0)
    if err != nil {
        log.Fatal(err)
    }

    win.Add(box)

    /* Menu layout */
    builder, err := gtk.BuilderNew()
    if err := builder.AddFromFile("/home/josh/Documents/projects/chroma-viz/gtk/menus.ui"); err != nil {
        log.Fatal(err)
    }

    menu, err := builder.GetObject("menubar")
    if err != nil {
        log.Fatal(err)
    }

    app.SetMenubar(menu.(*glib.MenuModel))

    importAction := glib.SimpleActionNew("import_show", nil)
    importAction.Connect("activate", func() { guiImportShow(win, showView, tempView) })
    app.AddAction(importAction)

    exportAction := glib.SimpleActionNew("export_show", nil)
    exportAction.Connect("activate", func() { guiExportShow(win, showView) })
    app.AddAction(exportAction)

    /* Body layout */
    bodyBox, err := gtk.PanedNew(gtk.ORIENTATION_HORIZONTAL)
    if err != nil {
        log.Fatal(err)
    }

    box.PackStart(bodyBox, true, true, 0)

    leftBox, err := gtk.PanedNew(gtk.ORIENTATION_VERTICAL)
    if err != nil {
        log.Fatal(err)
    }

    rightBox, err := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0)
    if err != nil {
        log.Fatal(err)
    }

    bodyBox.Pack1(leftBox, true, true)
    bodyBox.Pack2(rightBox, true, true)

    /* left */
    leftBox.SetHExpand(true)

    /* templates */
    templates, err := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0)
    if err != nil {
        log.Print(err)
    }

    leftBox.Pack1(templates, true, true)

    header1, err := gtk.HeaderBarNew()
    if err != nil {
        log.Fatal(err)
    }

    header1.SetTitle("Templates")
    templates.PackStart(header1, false, false, 0)

    scroll1, err := gtk.ScrolledWindowNew(nil, nil)
    if err != nil {
        log.Fatal(err)
    }

    templates.PackStart(scroll1, true, true, 0)
    scroll1.Add(tempView)

    /* show */
    shows, err := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0)
    if err != nil {
        log.Fatal(err)
    }

    leftBox.Pack2(shows, true, true)

    header2, err := gtk.HeaderBarNew()
    if err != nil {
        log.Fatal(err)
    }

    header2.SetTitle("Show")
    shows.PackStart(header2, false, false, 0)
    scroll2, err := gtk.ScrolledWindowNew(nil, nil)
    if err != nil {
        log.Fatal(err)
    }

    shows.PackStart(scroll2, true, true, 0)
    scroll2.Add(showView)

    /* right */
    preview := setup_preview_window()
    rightBox.PackStart(editView.Box(), true, true, 0)
    rightBox.PackEnd(preview, false, false, 0)

    /* Lower Bar layout */
    lowerBox, err := gtk.ActionBarNew()
    if err != nil {
        log.Fatal(err)
    }

    box.PackEnd(lowerBox, false, false, 0)

    button, err := gtk.ButtonNew()
    if err != nil {
        log.Fatal(err)
    }

    lowerBox.PackEnd(button)
    button.SetLabel("Exit")
    button.Connect("clicked", func() { gtk.MainQuit() })

    for name, render := range conn {
        eng := NewEngineWidget(name, render)
        lowerBox.PackStart(eng)
    }

    win.ShowAll()
}

func guiImportShow(win *gtk.ApplicationWindow, show *ShowTree, temp *TempTree) {
    dialog, err := gtk.FileChooserDialogNewWith2Buttons(
        "Import Show", win, gtk.FILE_CHOOSER_ACTION_OPEN, 
        "_Cancel", gtk.RESPONSE_CANCEL, "_Open", gtk.RESPONSE_ACCEPT)

    if err != nil {
        log.Print(err)
    }

    res := dialog.Run()

    if res == gtk.RESPONSE_ACCEPT {
        filename := dialog.GetFilename()
        show.ImportShow(temp, filename)
    }
    dialog.Destroy()
}

func guiExportShow(win *gtk.ApplicationWindow, show *ShowTree) {
    dialog, err := gtk.FileChooserDialogNewWith2Buttons(
        "Save Show", win, gtk.FILE_CHOOSER_ACTION_SAVE, 
        "_Cancel", gtk.RESPONSE_CANCEL, "_Save", gtk.RESPONSE_ACCEPT)

    if err != nil {
        log.Print(err)
    }

    dialog.SetCurrentName(".show")
    res := dialog.Run()
    if res == gtk.RESPONSE_ACCEPT {
        filename := dialog.GetFilename()
        show.ExportShow(filename)
    }

    dialog.Destroy()
}
