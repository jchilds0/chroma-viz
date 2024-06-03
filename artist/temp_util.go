package artist

import (
	"chroma-viz/library/attribute"
	"chroma-viz/library/pages"
	"chroma-viz/library/templates"
	"chroma-viz/library/util"
	"log"
	"strconv"

	"github.com/gotk3/gotk3/gtk"
)

func artistPageToTemplate(page pages.Page, tempView *TempTree, tempID, title, layer string) (temp *templates.Template, err error) {
	temp = page.CreateTemplate()
	temp.Title = title
	temp.TempID, err = strconv.ParseInt(tempID, 10, 64)
	if err != nil {
		return
	}

	temp.Layer, err = strconv.Atoi(layer)
	if err != nil {
		return
	}

	tempView.keyframes(temp)
	return
}

func updateParentGeometry(page *pages.Page, model *gtk.TreeModel, iter *gtk.TreeIter, parentID int) {
	ok := true
	for ok {
		geoID, err := util.ModelGetValue[int](model, iter, GEO_NUM)
		if err != nil {
			log.Print("Error getting prop id")
			ok = model.IterNext(iter)
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

		ok = model.IterNext(iter)
	}
}

func geometryToTreeView(page *pages.Page, tempView *TempTree, iter *gtk.TreeIter, propID int) {
	for id, geo := range page.PropMap {
		parentAttr := geo.Attr["parent"]
		if parentAttr == nil {
			log.Print("Error getting parent attr")
			continue
		}

		parent := parentAttr.(*attribute.IntAttribute)

		if parent.Value != propID {
			continue
		}

		newRow := tempView.geoModel.Append(iter)
		tempView.AddGeoRow(newRow, geo.Name, geoNames[geo.PropType], id)

		if id == propID {
			continue
		}

		geometryToTreeView(page, tempView, newRow, id)
	}
}
