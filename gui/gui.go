package gui

import (
	"fmt"
	"log"
	"os/exec"
	"strconv"

	"github.com/gotk3/gotk3/gtk"
)

// TODO: seperate preview setup from gui setup

func LaunchGui(conn map[string]*Connection) {
    gtk.Init(nil)

    win, _ := gtk.WindowNew(gtk.WINDOW_TOPLEVEL)
    win.SetTitle("Chroma Viz")
    win.Connect("destroy", func() { 
        gtk.MainQuit() 
    })

    editView := NewEditor(conn)
    showView := NewShow(editView, conn["Preview"])
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

    //box.PackStart(menuBox, false, false, 0)
 
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
    preview := setup_preview_window(conn["Preview"])
    rightBox.PackStart(editView.Box(), true, true, 0)
    rightBox.PackEnd(preview, false, false, 0)

    /* Lower Bar layout */
    lowerBox, _ := gtk.ActionBarNew()
    button, _ := gtk.ButtonNew()
    button.SetLabel("Exit")
    button.Connect("clicked", func() { gtk.MainQuit() })
    lowerBox.PackEnd(button)

    for name, render := range conn {
        eng := NewEngineWidget(name, render)
        lowerBox.PackStart(eng)
        box.PackEnd(lowerBox, false, false, 0)
    }
    
    win.Add(box)
    win.SetDefaultSize(800, 600)
    win.ShowAll()
    gtk.Main()
}

func setup_preview_window(conn *Connection) *gtk.Frame {
    soc, _ := gtk.SocketNew()
    window, _ := gtk.FrameNew("")
    window.Add(soc)

    window.Connect("draw", func(window *gtk.Frame) {
        width := window.GetAllocatedWidth()
        height := width * 9 / 16
        window.SetSizeRequest(-1, height)
    })

    var xid uint
    soc.SetVisible(true)
    soc.Connect("realize", func(soc *gtk.Socket) {
        xid = soc.GetId()
        prev := exec.Command("/home/josh/Documents/projects/chroma-engine/build/chroma-engine", "-wid", strconv.Itoa(int(xid)))
        log.Print(prev.String())

        if err := prev.Start(); err != nil {
            log.Fatal(err)
        }

        fmt.Printf("xid: %d\n", xid)
    })

    soc.Connect("plug-added", func() { log.Printf("Plug inserted: %d", xid) })
    window.SetVisible(true)

    return window 
}
