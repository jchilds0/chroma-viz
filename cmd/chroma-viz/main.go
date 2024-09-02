package main

import (
	"chroma-viz/library"
	"chroma-viz/library/pages"
	"chroma-viz/library/templates"
	"flag"
	"log"
	"math/rand"
	"net"
	"os"
	"runtime/pprof"
	"time"

	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
)

var profile = flag.String("profile", "", "write profile to file")
var configPath = flag.String("c", "./viz/conf.json", "config json")
var importRandom = flag.Int("t", 0, "import random pages")
var conf *library.Config
var numTemplates = 100

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
	conf, err = library.ImportConfig(*configPath)
	if err != nil {
		log.Fatal(err)
	}

	if *importRandom != 0 {
		importHook = importRandomPages
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

func importRandomPages(hub net.Conn, tempTree *TempTree, showTree *ShowTree) {
	start := time.Now()
	showTree.treeView.SetModel(nil)

	for i := 0; i < *importRandom; i++ {
		index := (rand.Int() % numTemplates) + 1
		template, err := templates.GetTemplate(hub, index)
		if err != nil {
			log.Fatal(err)
			return
		}

		page := pages.NewPageFromTemplate(&template)
		page.PageNum = showTree.show.NumPages
		showTree.ImportPage(page)
	}

	end := time.Now()
	elapsed := end.Sub(start)
	log.Printf("Built Show in %s\n", elapsed)

	showTree.treeView.SetModel(showTree.treeList)
}
