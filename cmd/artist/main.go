package main

import (
	"chroma-viz/library"
	"chroma-viz/library/hub"
	"flag"
	"log"
	"os"
	"runtime/pprof"

	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
)

var profile = flag.String("profile", "", "write profile to file")

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
	conn = make(map[string]*library.Connection)
	chromaHub, err = hub.NewDataBase(10)
	if err != nil {
		log.Fatal(err)
	}

	conf, err = library.ImportConfig("artist/conf.json")
	if err != nil {
		log.Fatal(err)
	}

	hubConn = library.NewConnection("Hub", conf.HubAddr, conf.HubPort)
	hubConn.Connect()

	for _, c := range conf.Connections {
		conn[c.Name] = library.NewConnection(c.Name, c.Address, c.Port)
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
		if c.IsConnected() {
			c.CloseConn()
			log.Printf("Closed %s\n", name)
		}
	}
}
