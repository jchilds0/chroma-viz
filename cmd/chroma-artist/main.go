package main

import (
	"chroma-viz/library"
	"flag"
	"log"
	"os"
	"runtime/pprof"

	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
)

var profile = flag.String("profile", "", "write profile to file")
var configPath = flag.String("c", "artist/conf.json", "config json")

func main() {
	flag.Parse()
	if *profile != "" {
		f, err := os.Create(*profile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}
	defer closeConn()

	var err error
	conf, err = library.ImportConfig(*configPath)
	if err != nil {
		log.Fatal(err)
	}

	app, err := gtk.ApplicationNew("app.chroma.artist", glib.APPLICATION_FLAGS_NONE)
	if err != nil {
		log.Print(err)
	}

	app.Connect("activate", ArtistGui)
	app.Run([]string{})
}

func closeConn() {
	for name, c := range conn {
		c.CloseConn()
		log.Printf("Closed %s\n", name)
	}
}
