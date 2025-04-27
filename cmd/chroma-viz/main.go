package main

import (
	"chroma-viz/library"
	"chroma-viz/library/hub"
	"chroma-viz/library/pages"
	"chroma-viz/library/templates"
	"flag"
	"fmt"
	"log"
	"math/rand"
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

	for _, c := range conf.Connections {
		if c.Type == "engine" {
			conn.eng = append(conn.eng, library.NewConnection(c.Name, c.Address, c.Port))
		} else if c.Type == "preview" {
			conn.prev = append(conn.prev, library.NewConnection(c.Name, c.Address, c.Port))
		} else {
			log.Printf("Unknown connection type %s", c.Type)
		}
	}

	app, err := gtk.ApplicationNew(conf.Name+".chroma.viz", glib.APPLICATION_FLAGS_NONE)
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

		c.CloseConn()
		log.Printf("Closed %s\n", c.Name)
	}

	for _, c := range conn.prev {
		if c == nil {
			continue
		}

		c.CloseConn()
		log.Printf("Closed %s\n", c.Name)
	}
}

func importRandomPages(c hub.Client, tempTree *TempTree, showTree ShowTree) {
	start := time.Now()

	show := showTree.(*MediaSequencer)
	show.treeView.SetModel(nil)

	for _ = range *importRandom {
		index := (rand.Int() % numTemplates) + 1
		path := fmt.Sprintf("/template/%d", index)

		var template templates.Template
		err := conf.ChromaHub.GetJSON(path, &template)
		if err != nil {
			log.Fatal(err)
			return
		}

		err = template.Init(true)
		if err != nil {
			log.Print(err)
			return
		}

		page := pages.NewPage(&template)
		showTree.WritePage(*page)
	}

	end := time.Now()
	elapsed := end.Sub(start)
	log.Printf("Built Show in %s\n", elapsed)

	show.treeView.SetModel(show.treeList)
}
