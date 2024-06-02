package viz

import (
	"chroma-viz/hub"
	"chroma-viz/library/pages"
	"chroma-viz/library/props"
	"chroma-viz/library/templates"
	"fmt"
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
	chromaHub, err := hub.NewDataBase(numTemplates)
	if err != nil {
		t.Error(err)
	}

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

	geos := []string{templates.GEO_RECT, templates.GEO_CIRCLE, templates.GEO_TEXT}
	props := []string{props.RECT_PROP, props.CIRCLE_PROP, props.TEXT_PROP}

	for j := 1; j < numGeometries; j++ {
		geoIndex := rand.Int() % len(geos)

		geo := templates.NewGeometry(
			j,
			props[geoIndex],
			props[geoIndex],
			geos[geoIndex],
			rand.Int()%2000,
			rand.Int()%2000,
			0,
		)

		color := fmt.Sprintf("%f %f %f %f", rand.Float64(), rand.Float64(), rand.Float64(), rand.Float64())

		switch geos[geoIndex] {
		case templates.GEO_RECT:
			rect := templates.NewRectangle(
				*geo,
				rand.Int()%1000,
				rand.Int()%1000,
				rand.Int()%10,
				color,
			)
			err = chromaHub.AddRectangle(tempID, *rect)

		case templates.GEO_CIRCLE:
			circle := templates.NewCircle(
				*geo,
				rand.Int()%200,
				rand.Int()%200,
				rand.Int()%10,
				rand.Int()%200,
				color,
			)
			err = chromaHub.AddCircle(tempID, *circle)
		case templates.GEO_TEXT:
			text := templates.NewText(*geo, "some text", color)
			err = chromaHub.AddText(tempID, *text)
		}

		if j%3 == 0 {
			var tempFrame *templates.Keyframe
			if j%2 == 0 {
				tempFrame = templates.NewKeyFrame(0, j, "rel_x", false, false)
			} else {
				tempFrame = templates.NewKeyFrame(0, j, "rel_y", false, false)
			}

			startFrame := templates.NewSetFrame(*tempFrame, rand.Int()%2000)
			chromaHub.AddSetFrame(tempID, *startFrame)

			tempFrame.FrameNum = 1

			endFrame := templates.NewUserFrame(*tempFrame)
			chromaHub.AddUserFrame(tempID, *endFrame)

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
