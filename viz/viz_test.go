package viz

import (
	"chroma-viz/hub"
	"chroma-viz/library/pages"
	"chroma-viz/library/props"
	"chroma-viz/library/templates"
	"log"
	"math/rand"
	"net"
	"os"
	"runtime/pprof"
	"strconv"
	"testing"
	"time"

	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
)

var numTemplates = 100
var numPages = 1_000
var numGeometries = 100

func TestGui(t *testing.T) {
	defer CloseViz()

	createHub := true

	f, err := os.Create("../perf/viz_test.prof")
	if err != nil {
		log.Fatal(err)
	}

	pprof.StartCPUProfile(f)
	defer pprof.StopCPUProfile()
	log.Printf(
		"Testing with %d Templates, %d Pages and %d Geometries\n",
		numTemplates, numPages, numGeometries,
	)

	importHook = importRandomPages
	chromaHub := hub.NewDataBase(numTemplates)

	if createHub {
		err = chromaHub.CleanDB()
		if err != nil {
			log.Print(err)
		}

		log.Printf("Cleaned out graphics hub")

		start := time.Now()

		var i int64
		for i = 1; i <= int64(numTemplates); i++ {
			randomTemplate(chromaHub, i)
		}

		end := time.Now()
		elapsed := end.Sub(start)
		log.Printf("Built Graphics Hub in %s\n", elapsed)
	}

	go chromaHub.StartHub(9000)

	time.Sleep(time.Second)
	InitialiseViz("./conf.json")

	app, err := gtk.ApplicationNew("app.chroma.viz", glib.APPLICATION_FLAGS_NONE)
	if err != nil {
		log.Print(err)
	}

	app.Connect("activate", VizGui)
	app.Run([]string{})
}

func randomTemplate(chromaHub *hub.DataBase, tempID int64) {
	err := chromaHub.AddTemplate(tempID, "Template "+strconv.FormatInt(tempID, 10), 0)
	if err != nil {
		log.Fatalf("Error adding template (%s)", err)
	}

	geos := []int{props.RECT_PROP, props.TEXT_PROP, props.CIRCLE_PROP}

	for j := 0; j < numGeometries; j++ {
		geoIndex := rand.Int() % len(geos)
		prop := geos[geoIndex]

		geo := templates.NewGeometry(
			props.PropType(prop),
			prop,
			rand.Int()%2000,
			rand.Int()%2000,
			byte(rand.Int()%255),
			byte(rand.Int()%255),
			byte(rand.Int()%255),
			byte(rand.Int()%255),
			0,
		)

		switch prop {
		case props.RECT_PROP:
			rect := templates.NewRectangle(
				*geo,
				rand.Int()%1000,
				rand.Int()%1000,
				rand.Int()%10,
			)
			err = chromaHub.AddRectangle(tempID, *rect)

		case props.CIRCLE_PROP:
			circle := templates.NewCircle(
				*geo,
				rand.Int()%1000,
				rand.Int()%1000,
				rand.Int()%1000,
				rand.Int()%1000,
			)
			err = chromaHub.AddCircle(tempID, *circle)
		case props.TEXT_PROP:
			text := templates.NewText(*geo, "some text")
			err = chromaHub.AddText(tempID, *text)
		}

		if err != nil {
			log.Fatalf("Error adding attributes (%s)", err)
		}
	}
}

func importRandomPages(hub net.Conn, tempTree *TempTree, showTree *ShowTree) {
	start := time.Now()
	showTree.treeView.SetModel(nil)

	for i := 0; i < numPages; i++ {
		index := (rand.Int() % numTemplates) + 1
		page, err := pages.GetPage(hub, index)
		if err != nil {
			log.Print(err)
			continue
		}

		page.PageNum = showTree.show.NumPages
		showTree.ImportPage(page)
	}

	end := time.Now()
	elapsed := end.Sub(start)
	log.Printf("Built Show in %s\n", elapsed)

	showTree.treeView.SetModel(showTree.treeList)
}
