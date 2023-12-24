package main

import (
	"chroma-viz/artist"
	"chroma-viz/viz"
	"flag"
	"log"
	"os"
	"runtime/pprof"

	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
)

var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")
var mode = flag.String("mode", "", "chroma mode (artist | viz)")

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

    var app *gtk.Application
    var err error
    if *mode == "artist" {
        artist.InitConnections()
        defer artist.CloseConn()

        artist.AddConnection("Preview", "127.0.0.1", 6100)

        app, err = gtk.ApplicationNew("app.chroma.artist", glib.APPLICATION_FLAGS_NONE)
        if err != nil {
            log.Print(err)
        }

        app.Connect("activate", artist.ArtistGui)
    } else if *mode == "viz" {
        viz.InitConnections()
        defer viz.CloseConn()

        viz.AddConnection("Engine", "127.0.0.1", 6800)
        viz.AddConnection("Preview", "127.0.0.1", 6100)

        app, err = gtk.ApplicationNew("app.chroma.viz", glib.APPLICATION_FLAGS_NONE)
        if err != nil {
            log.Print(err)
        }

        app.Connect("activate", viz.VizGui)
    } else {
        flag.PrintDefaults()
        return
    }

    app.Run([]string{})
}

