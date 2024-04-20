package viz

import (
	"chroma-viz/hub"
	"chroma-viz/library/attribute"
	"chroma-viz/library/pages"
	"chroma-viz/library/props"
	"fmt"
	"log"
	"math/rand"
	"net"
	"os"
	"reflect"
	"runtime/pprof"
	"strconv"
	"testing"
	"time"

	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
)

var numTemplates = 100
var numPages = 10_000
var numGeometries = 10

func TestGui(t *testing.T) {
	defer CloseViz()

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

	start := time.Now()

	// var i int64
	// for i = 1; i <= int64(numTemplates); i++ {
	// 	randomTemplate(chromaHub, i)
	// }

	end := time.Now()
	elapsed := end.Sub(start)
	log.Printf("Built Graphics Hub in %s\n", elapsed)

	go hub.StartHub(chromaHub, 9000)

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

	geos := []*props.Property{
		props.NewProperty(props.RECT_PROP, "rect", true, nil),
		props.NewProperty(props.TEXT_PROP, "text", true, nil),
		props.NewProperty(props.CIRCLE_PROP, "circle", true, nil),
	}

	for j := 0; j < numGeometries; j++ {
		geoIndex := rand.Int() % len(geos)
		prop := geos[geoIndex]

		geo_id, err := chromaHub.AddGeometry(tempID, prop.Name, props.PropType(prop.PropType))
		if err != nil {
			log.Fatalf("Error adding geometry (%s)", err)
		}

		for name, attr := range prop.Attr {
			var value string
			switch attr.(type) {
			case *attribute.IntAttribute:
				if name == "rel_x" || name == "rel_y" {
					value = strconv.Itoa(rand.Int() % 2000)
				} else if name != "parent" {
					value = strconv.Itoa(rand.Int() % 200)
				}

			case *attribute.ColorAttribute:
				value = fmt.Sprintf("%f %f %f %f",
					float64(rand.Int()%255)/255,
					float64(rand.Int()%255)/255,
					float64(rand.Int()%255)/255,
					1.0,
				)

			case *attribute.StringAttribute:
				value = "some text"
			}

			_, err := chromaHub.AddAttribute(geo_id, name, value, reflect.TypeOf(attr).String()[1:], true)
			if err != nil {
				log.Fatalf("Error adding attribute (%s)", err)
			}
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
