package viz

import (
	"chroma-viz/hub"
	"log"
	"math/rand"
	"strconv"
	"testing"
	"time"

	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
)

var numTemplates = 100
var numPages = 100_000
var numGeometries = 100

func TestGui(t *testing.T) {
	defer CloseViz()
	log.Printf(
		"Testing with %d Templates, %d Pages and %d Geometries\n",
		numTemplates, numPages, numGeometries,
	)

	importHook = importRandomPages
	chromaHub := hub.NewDataBase()

    start := time.Now()
	geo := []string{"rect", "text", "circle", "image"}

	for i := 1; i <= numTemplates; i++ {
		chromaHub.AddTemplate(i, "", "", "")
		chromaHub.Templates[i].Title = "Template " + strconv.Itoa(i)
		numGeo := 0

		for j := 0; j < numGeometries; j++ {
			geoIndex := rand.Int() % len(geo)
			chromaHub.AddGeometry(i, numGeo, geo[geoIndex])
			numGeo++
		}
	}

	end := time.Now()
	elapsed := end.Sub(start)
	log.Printf("Built Graphics Hub in %s\n", elapsed)

	go hub.StartHub(chromaHub, 9000, -1)

    time.Sleep(time.Second)
	InitialiseViz("./conf.json")

	app, err := gtk.ApplicationNew("app.chroma.viz", glib.APPLICATION_FLAGS_NONE)
	if err != nil {
		log.Print(err)
	}

	app.Connect("activate", VizGui)
	app.Run([]string{})
}

func importRandomPages(tempTree *TempTree, showTree *ShowTree) {

	temps := make([]int, 0, len(tempTree.Temps.Temps))

	for _, temp := range tempTree.Temps.Temps {
		if temp == nil {
			continue
		}

		temps = append(temps, temp.TempID)
	}

	if len(temps) == 0 {
		log.Fatal("Error: Graphics Hub is empty")
	}

	start := time.Now()
	for i := 0; i < numPages; i++ {
		index := rand.Int() % len(temps)
		template := tempTree.Temps.Temps[temps[index]]
		page := showTree.show.AddPage(template.Title, template)
		showTree.ImportPage(page)
	}

	t := time.Now()
	elapsed := t.Sub(start)
	log.Printf("Built Show in %s\n", elapsed)
}
