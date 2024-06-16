package main

import (
	"chroma-viz/library/attribute"
	"chroma-viz/library/pages"
	"chroma-viz/library/util"
	"log"

	"github.com/gotk3/gotk3/gtk"
)

func updateParentGeometry(page *pages.Page, model *gtk.TreeModel, iter *gtk.TreeIter, parentID int) {
	nextIterExists := true

	for nextIterExists {
		geoID, err := util.ModelGetValue[int](model, iter, GEO_NUM)
		if err != nil {
			log.Print("Error getting prop id")
			nextIterExists = model.IterNext(iter)
			continue
		}

		attr := page.PropMap[geoID].Attr["parent"]
		if attr == nil {
			log.Printf("Missing parent attr")
			return
		}

		intAttr, ok := attr.(*attribute.IntAttribute)
		if !ok {
			log.Printf("Missing int attr")
			return
		}

		intAttr.Value = parentID

		var childIter gtk.TreeIter
		if ok = model.IterChildren(iter, &childIter); ok {
			updateParentGeometry(page, model, &childIter, geoID)
		}

		nextIterExists = model.IterNext(iter)
	}
}
