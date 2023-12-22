package main 

import (
	"chroma-viz/gui"
	"log"
	"math"
	"time"

	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
)

func TestGui() {
    gui.InitConnections()
    gui.AddConnection("Engine", "127.0.0.1", 6800)
    gui.AddConnection("Preview", "127.0.0.1", 6100)

    app, err := gtk.ApplicationNew("app.chroma.viz", glib.APPLICATION_FLAGS_NONE)
    if err != nil {
        log.Print(err)
    }

    gui.InitGui()
    loagPages()
    app.Connect("activate", gui.MainGui)
    app.Run([]string{})

    gui.CloseConn()
}

func loagPages() {
    var page *gui.Template
    num_temps := int(math.Pow(10, 4))
    num_props := 100
    num_pages := 100

    log.Printf("Testing with %d Templates, %d Properties, %d Pages\n", num_temps, num_props, num_pages)

    start := time.Now()
    for i := 1; i < num_temps; i++ {
        page = gui.TempView.AddTemplate("Template", i, gui.LOWER_FRAME)

        for j := 0; j < num_props; j++ {
            page.AddProp("Background", gui.RectProp)
            page.AddProp("Text", gui.TickerProp)
        }
    }

    t := time.Now()
    elapsed := t.Sub(start)
    log.Printf("Built Templates in %s\n", elapsed)

    start = time.Now()
    for i := 0; i < num_pages; i++ {
        gui.ShowView.NewShowPage(page)
    }

    t = time.Now()
    elapsed = t.Sub(start)
    log.Printf("Built Templates in %s\n", elapsed)
}
