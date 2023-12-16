package gui

import (
	"log"
	"os/exec"
	"strconv"

	"github.com/gotk3/gotk3/gtk"
)

const engDir = "/home/josh/Documents/projects/chroma-engine/build/chroma-engine"

func setup_preview_window() *gtk.Frame {
    soc, err := gtk.SocketNew()
    if err != nil {
        log.Fatalf("Error setting up preview window (%s)", err)
    }

    window, err := gtk.FrameNew("")
    if err != nil {
        log.Fatalf("Error setting up preview window (%s)", err)
    }

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
        prev := exec.Command(engDir, "-wid", strconv.Itoa(int(xid)))
        log.Print(prev.String())

        if err := prev.Start(); err != nil {
            log.Fatal(err)
        }
    })

    soc.Connect("plug-added", func() { log.Printf("Plug inserted: %d", xid) })
    window.SetVisible(true)

    return window 
}
