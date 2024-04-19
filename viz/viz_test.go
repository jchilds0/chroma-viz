package viz

import (
	"chroma-viz/hub"
	"chroma-viz/library/attribute"
	"chroma-viz/library/pages"
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

var numTemplates = 10_000
var numPages = 10_000
var numGeometries = 100

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
	geo := []string{"rect", "text", "circle"}

	for i := 1; i <= numTemplates; i++ {
		chromaHub.AddTemplate(i, "", "", "")
		chromaHub.Templates[i].Title = "Template " + strconv.Itoa(i)

		for j := 0; j < numGeometries; j++ {
			geoIndex := rand.Int() % len(geo)
			chromaHub.AddGeometry(i, j, geo[geoIndex])
		}
	}

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

func importRandomPages(hub net.Conn, tempTree *TempTree, showTree *ShowTree) {
	start := time.Now()
	for i := 0; i < numPages; i++ {
		index := (rand.Int() % numTemplates) + 1
		page, err := pages.GetPage(hub, index)
		if err != nil {
			log.Print(err)
			continue
		}

		page.PageNum = showTree.show.NumPages
        for _, prop := range page.PropMap {
            for name, attr := range prop.Attr {
                if name == "parent" {
                    continue
                }

                prop.Visible[name] = true

                switch a := attr.(type) {
                case *attribute.IntAttribute:
                    if (name == "rel_x" || name == "rel_y") {
                        a.Value = rand.Int() % 2000
                    } else {
                        a.Value = rand.Int() % 200
                    }
                case *attribute.ColorAttribute:
                    a.Red = float64(rand.Int() % 255) / 255
                    a.Green = float64(rand.Int() % 255) / 255
                    a.Blue = float64(rand.Int() % 255) / 255
                    a.Alpha = 1.0
                case *attribute.StringAttribute:
                    a.Value = "some text"
                }
            }
        }
		showTree.ImportPage(page)
	}

	end := time.Now()
	elapsed := end.Sub(start)
	log.Printf("Built Show in %s\n", elapsed)
}

