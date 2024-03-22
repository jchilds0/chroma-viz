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
        propID, err := gtk_utils.ModelGetValue[int](model, iter, PROP_NUM)
        if err != nil {
            log.Print("Error getting prop id")
            continue
        }

        prop := template.Geometry[propID]
        if prop == nil {
            log.Print("Error getting prop")
            continue
        }

        parentAttr := prop.Attr["parent"]
        if parentAttr == nil {
            log.Print("Error getting parent attr")
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

func compressGeometry(temp, newTemp *templates.Template) {
    // build geo id map 
    geoRename := make(map[int]int, len(template.Geometry))

    i := 1
    for id := range temp.Geometry {
        geoRename[id] = i 
        i++
    }

    // copy geo's from temp to newTemp
    for id, geo := range temp.Geometry {
        newID := geoRename[id]
        newTemp.Geometry[newID] = geo

        parentAttr := geo.Attr["parent"]
        if parentAttr == nil {
            continue
        }

        attr := parentAttr.(*attribute.IntAttribute)
        attr.Value = geoRename[attr.Value]
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

        geo.Visible = visible
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

        newRow := tempView.model.Append(iter)
        tempView.AddRow(newRow, geo.Name, geo_name[geo.PropType], id)
        geometryToTreeView(tempView, newRow, id)
    }
}
