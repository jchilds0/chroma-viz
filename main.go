package main

import (
	"chroma-viz/gui"
	"log"

	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
)

func main() {
    gui.InitConnections()
    gui.AddConnection("Engine", "127.0.0.1", 6800)
    gui.AddConnection("Preview", "127.0.0.1", 6100)

    app, err := gtk.ApplicationNew("app.chroma.viz", glib.APPLICATION_FLAGS_NONE)
    if err != nil {
        log.Print(err)
    }

    app.Connect("activate", gui.SetupMainGui)
    app.Run([]string{})
}

