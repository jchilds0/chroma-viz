package main

import (
	"chroma-viz/library/hub"
	"chroma-viz/library/props"
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

	geos := []string{templates.GEO_RECT, templates.GEO_CIRCLE, templates.GEO_TEXT}
	props := []string{props.RECT_PROP, props.CIRCLE_PROP, props.TEXT_PROP}

	for j := 1; j < numGeo; j++ {
		geoIndex := rand.Int() % len(geos)

		geo := templates.Geometry{
			GeoNum:   j,
			Name:     props[geoIndex],
			GeoType:  geos[geoIndex],
			PropType: props[geoIndex],
			RelX:     rand.Int() % 2000,
			RelY:     rand.Int() % 2000,
		}

		color := fmt.Sprintf("%f %f %f %f", rand.Float64(), rand.Float64(), rand.Float64(), rand.Float64())

		switch geos[geoIndex] {
		case templates.GEO_RECT:
			rect := templates.NewRectangle(
				geo,
				rand.Int()%1000,
				rand.Int()%1000,
				rand.Int()%10,
				color,
			)
			err = chromaHub.AddRectangle(tempID, *rect)

		case templates.GEO_CIRCLE:
			circle := templates.NewCircle(
				geo,
				rand.Int()%200,
				rand.Int()%200,
				rand.Int()%10,
				rand.Int()%200,
				color,
			)
			err = chromaHub.AddCircle(tempID, *circle)
		case templates.GEO_TEXT:
			text := templates.NewText(geo, "some text", color, 1.0)
			err = chromaHub.AddText(tempID, *text)
		}

		if j%3 == 0 {
			var tempFrame *templates.Keyframe
			if j%2 == 0 {
				tempFrame = templates.NewKeyFrame(0, j, "rel_x", false)
			} else {
				tempFrame = templates.NewKeyFrame(0, j, "rel_y", false)
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
