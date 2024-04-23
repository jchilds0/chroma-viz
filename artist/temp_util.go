package artist

import (
	"chroma-viz/library/attribute"
	"chroma-viz/library/gtk_utils"
	"chroma-viz/library/pages"
	"chroma-viz/library/templates"
	"log"
	"strconv"

	"github.com/gotk3/gotk3/gtk"
)

func artistPageToTemplate(page pages.Page, tempView *TempTree, tempID, title, layer string) (temp *templates.Template, err error) {
	template := page.CreateTemplate()
	template.Title = title
	template.TempID, err = strconv.ParseInt(tempID, 10, 64)
	if err != nil {
		return
	}

	template.Layer, err = strconv.Atoi(layer)
	if err != nil {
		return
	}

	template.Keyframe = tempView.keyframes()

	// update parent
	model := tempView.geoModel.ToTreeModel()
	if iter, ok := model.GetIterFirst(); ok {
		updateParentGeometry(template, model, iter, 0)
	}

	return
}

func updateParentGeometry(template *templates.Template, model *gtk.TreeModel, iter *gtk.TreeIter, parentID int) {
	ok := true
	for ok {
		geoID, err := gtk_utils.ModelGetValue[int](model, iter, GEO_NUM)
		if err != nil {
			log.Print("Error getting prop id")
			ok = model.IterNext(iter)
			continue
		}

		var geom *templates.Geometry
		for _, geo := range template.Geometry {
			geom = geo.Geom()

			if geom.GeoID != geoID {
				continue
			}

			geom.Parent = parentID
			break
		}

		var childIter gtk.TreeIter
		if ok = model.IterChildren(iter, &childIter); ok {
			updateParentGeometry(template, model, &childIter, geoID)
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
		tempView.AddGeoRow(newRow, geo.Name, geo_name[geo.PropType], id)
		geometryToTreeView(page, tempView, newRow, id)
	}
}
