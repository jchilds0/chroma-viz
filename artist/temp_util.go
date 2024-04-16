package artist

import (
	"chroma-viz/library/attribute"
	"chroma-viz/library/gtk_utils"
	"chroma-viz/library/props"
	"chroma-viz/library/templates"
	"log"

	"github.com/gotk3/gotk3/gtk"
)

func updateParentGeometry(model *gtk.TreeModel, iter *gtk.TreeIter, parentID int) {
	ok := true
	for ok {
		propID, err := gtk_utils.ModelGetValue[int](model, iter, GEO_NUM)
		if err != nil {
			log.Print("Error getting prop id")
			ok = model.IterNext(iter)
			continue
		}

		prop := template.Geometry[propID]
		if prop == nil {
			log.Print("Error getting prop")
			ok = model.IterNext(iter)
			continue
		}

		parentAttr := prop.Attr["parent"]
		if parentAttr == nil {
			log.Print("Error getting parent attr")
			ok = model.IterNext(iter)
			continue
		}

		parentAttr.(*attribute.IntAttribute).Value = parentID

		var childIter gtk.TreeIter
		if ok = model.IterChildren(iter, &childIter); ok {
			updateParentGeometry(model, &childIter, propID)
		}

		ok = model.IterNext(iter)
	}
}

/*
Chroma Engine renders geometry using the geometry index,
from the first index to the last. To make this consistent
with the heirachy arranged by the user, we walk the tree
model. We add the geometries by depth, adding everything
at depth 0, then depth 1 and so on.
*/
func compressGeometry(temp, newTemp *templates.Template, tree *gtk.TreeModel) {
	// build geo id map
	geoRename := make(map[int]int, len(template.Geometry))
	geoIters := make([]*gtk.TreeIter, 0, len(template.Geometry))

	i := 1
	first, ok := tree.GetIterFirst()
	if !ok {
		log.Printf("Error compressing geometry, no tree geometries")
		return
	}

	geoIters = append(geoIters, first)

	var childIter gtk.TreeIter
	for len(geoIters) > 0 {
		iter := geoIters[0]
		geoIters = geoIters[1:]

		geoNum, err := gtk_utils.ModelGetValue[int](tree, iter, GEO_NUM)
		if err != nil {
			log.Printf("Error getting geo num (%s)", err)
			continue
		}

		geoRename[geoNum] = i
		i++

		// add children of iter to the geoIters array
		ok := tree.IterChildren(iter, &childIter)
		if !ok {
			continue
		}

		newIter, err := childIter.Copy()
		if err != nil {
			log.Printf("Error copying iter (%s)", err)
			continue
		}

		geoIters = append(geoIters, newIter)
		ok = tree.IterNext(&childIter)
		for ok {
			newIter, err := childIter.Copy()
			if err != nil {
				log.Printf("Error copying iter (%s)", err)
				continue
			}

			geoIters = append(geoIters, newIter)
			ok = tree.IterNext(&childIter)
		}
	}

	// copy geo's from temp to newTemp
	for id, geo := range temp.Geometry {
		if geo == nil {
			continue
		}

		newID := geoRename[id]
		newTemp.Geometry[newID] = geo

		parentAttr := geo.Attr["parent"]
		if parentAttr == nil {
			continue
		}

		attr := parentAttr.(*attribute.IntAttribute)
		attr.Value = geoRename[attr.Value]
	}

    // update keyframe geo id's
    for _, frame := range temp.Keyframe {
        frame.FrameGeo = geoRename[frame.FrameGeo]
        frame.BindGeo = geoRename[frame.BindGeo]
        newTemp.Keyframe = append(newTemp.Keyframe, frame)
    }
}

func decompressGeometry(temp, newTemp *templates.Template) {
	// build a map of new geo id's to alloc geo id's
	geoRename := make(map[int]int, len(template.Geometry))

	// alloc new geo id's
	for id, geo := range newTemp.Geometry {
		geom, ok := geoms[geo.PropType]
		if !ok {
			log.Printf("Missing Geom %s", props.PropType(geo.PropType))
			continue
		}

		newID, err := geom.allocGeom()
		if err != nil {
			log.Print(err)
			continue
		}

		geoRename[id] = newID
	}

	// copy newTemp geo to temp using geoRename
	for id, geo := range newTemp.Geometry {
		newID := geoRename[id]
		temp.Geometry[newID] = geo

		parentAttr := geo.Attr["parent"]
		if parentAttr == nil {
			log.Printf("Missing parent attr for geo %s", geo.Name)
			continue
		}

		parent := parentAttr.(*attribute.IntAttribute)
		parent.Value = geoRename[parent.Value]
	}
}

func geometryToTreeView(tempView *TempTree, iter *gtk.TreeIter, propID int) {
	for id, geo := range template.Geometry {
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
		geometryToTreeView(tempView, newRow, id)
	}
}
