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
var conf *library.Config

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
	defer closeViz()

	var err error
	conf, err = library.ImportConfig("viz/conf.json")
	if err != nil {
		log.Fatal(err)
	}

	conn.hub = library.NewConnection("Hub", conf.HubAddr, conf.HubPort)
	conn.hub.Connect()

	for _, c := range conf.Connections {
		if c.Type == "engine" {
			conn.eng = append(conn.eng, library.NewConnection(c.Name, c.Address, c.Port))
		} else if c.Type == "preview" {
			conn.prev = append(conn.prev, library.NewConnection(c.Name, c.Address, c.Port))
		} else {
			log.Printf("Unknown connection type %s", c.Type)
		}
	}

	app, err := gtk.ApplicationNew("app.chroma.viz", glib.APPLICATION_FLAGS_NONE)
	if err != nil {
		log.Print(err)
	}

	app.Connect("activate", VizGui)
	app.Run([]string{})
}

func closeViz() {
	for _, c := range conn.eng {
		if c == nil {
			continue
		}

		if c.IsConnected() {
			c.CloseConn()
			log.Printf("Closed %s\n", c.Name)
		}
	}

	for _, c := range conn.prev {
		if c == nil {
			continue
		}

		if c.IsConnected() {
			c.CloseConn()
			log.Printf("Closed %s\n", c.Name)
		}
	}
}
