package gui

import (
	"fmt"
	"log"
	"os/exec"
	"strconv"

	"github.com/gotk3/gotk3/gtk"
)

// TODO: seperate preview setup from gui setup

func LaunchGui(conn *Connection) {
    gtk.Init(nil)

    win, _ := gtk.WindowNew(gtk.WINDOW_TOPLEVEL)
    win.SetTitle("Chroma Viz")
    win.Connect("destroy", func() { 
        gtk.MainQuit() 
    })

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
    bodyBox, _ := gtk.PanedNew(gtk.ORIENTATION_HORIZONTAL)
    box.PackStart(bodyBox, true, true, 0)

    editView := NewEditor(conn)
    showView := NewShow(editView)
    tempView := NewTempList(showView)

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
    preview := setup_preview_window(conn)
    rightBox.PackStart(editView.Box(), true, true, 0)
    rightBox.PackEnd(preview, false, false, 0)

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

func setup_preview_window(conn *Connection) *gtk.Frame {
    var prev *exec.Cmd
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
        prev = exec.Command("/home/josh/Documents/projects/chroma-engine/build/chroma-engine", "-wid", strconv.Itoa(int(xid)))
        log.Print(prev.String())

        if err := prev.Start(); err != nil {
            log.Fatal(err)
        }

        fmt.Printf("xid: %d\n", xid)
    })

    // window.Connect("destroy", func() { 
    //     soc.Destroy()
    //     if prev != nil && prev.Process != nil {
    //         time.Sleep(5 * time.Second)
    //         prev.Process.Kill() 
    //     }
    // })

    soc.Connect("plug-added", func() { log.Printf("Plug inserted: %d", xid) })
    window.SetVisible(true)

    return window 
}
