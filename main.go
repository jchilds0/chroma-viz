package main

import (
	"chroma-viz/gui"
	"flag"
	"log"
	"os"
	"runtime/pprof"

	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
)
var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")

func main() {

    flag.Parse()
    if *cpuprofile != "" {
        f, err := os.Create(*cpuprofile)
        if err != nil {
            log.Fatal(err)
        }
        pprof.StartCPUProfile(f)
        defer pprof.StopCPUProfile()
    }

    gui.InitConnections()
    gui.AddConnection("Engine", "127.0.0.1", 6800)
    gui.AddConnection("Preview", "127.0.0.1", 6100)

    app, err := gtk.ApplicationNew("app.chroma.viz", glib.APPLICATION_FLAGS_NONE)
    if err != nil {
        log.Print(err)
    }

    gtk.Init(nil)

    app.Connect("activate", gui.MainGui)
    app.Run([]string{})

    gui.CloseConn()
}

