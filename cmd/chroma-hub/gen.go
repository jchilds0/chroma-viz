package main

import (
	"chroma-viz/library/geometry"
	"chroma-viz/library/hub"
	"chroma-viz/library/templates"
	"fmt"
	"log"
	"math/rand"
	"strconv"
)

func randomTemplate(chromaHub *hub.DataBase, tempID int64, numGeo int) {
	err := chromaHub.AddTemplate(tempID, "Template "+strconv.FormatInt(tempID, 10), 0)
	if err != nil {
		log.Fatalf("Error adding template (%s)", err)
	}

	geos := []string{geometry.GEO_RECT, geometry.GEO_CIRCLE, geometry.GEO_TEXT}

	for j := 1; j < numGeo; j++ {
		geoIndex := rand.Int() % len(geos)

		geo := geometry.Geometry{
			GeometryID: j,
			Name:       geos[geoIndex],
			GeoType:    geos[geoIndex],
		}

		geo.RelX.Value = rand.Int() % 2000
		geo.RelY.Value = rand.Int() % 2000

		color := fmt.Sprintf("%f %f %f %f", rand.Float64(), rand.Float64(), rand.Float64(), rand.Float64())

		switch geos[geoIndex] {
		case geometry.GEO_RECT:
			rect := geometry.NewRectangle(geo)
			rect.Width.Value = rand.Int() % 1000
			rect.Height.Value = rand.Int() % 1000
			rect.Rounding.Value = rand.Int() % 10
			rect.Color.FromString(color)
			err = chromaHub.AddRectangle(tempID, *rect)

		case geometry.GEO_CIRCLE:
			circle := geometry.NewCircle(geo)
			circle.InnerRadius.Value = rand.Int() % 200
			circle.OuterRadius.Value = rand.Int() % 200
			circle.StartAngle.Value = rand.Int() % 10
			circle.EndAngle.Value = rand.Int() % 200
			circle.Color.FromString(color)
			err = chromaHub.AddCircle(tempID, *circle)

		case geometry.GEO_TEXT:
			text := geometry.NewText(geo)
			text.String.Value = "some text"
			text.Scale.Value = 1.0
			text.Color.FromString(color)
			err = chromaHub.AddText(tempID, *text)

		}

		if j%3 == 0 {
			var tempFrame *templates.Keyframe
			if j%2 == 0 {
				tempFrame = templates.NewKeyFrame(1, j, "rel_x", false)
			} else {
				tempFrame = templates.NewKeyFrame(1, j, "rel_y", false)
			}

			startFrame := templates.NewSetFrame(*tempFrame, rand.Int()%2000)
			chromaHub.AddSetFrame(tempID, *startFrame)

			tempFrame.FrameNum = 2

			endFrame := templates.NewUserFrame(*tempFrame)
			chromaHub.AddUserFrame(tempID, *endFrame)

		}

		if err != nil {
			log.Fatalf("Error adding attributes (%s)", err)
		}
	}
}
