package main

import (
	"chroma-viz/library/attribute"
	"chroma-viz/library/pages"
	"chroma-viz/library/util"
	"log"

	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
)

const (
	GEO_TYPE = iota
	GEO_NAME
	GEO_NUM
	GEO_NUM_COLS
)

type GeoTree struct {
	currentPage *pages.Page
	geoModel    *gtk.TreeStore
	geoView     *gtk.TreeView
	geoList     *gtk.TreeStore
}

func NewGeoTree(propToEditor func(propID int)) (geoTree *GeoTree, err error) {
	geoTree = &GeoTree{}

	geoTree.geoView, err = gtk.TreeViewNew()
	if err != nil {
		return
	}

	geoTree.geoView.Set("reorderable", true)

	typeCell, err := gtk.CellRendererTextNew()
	if err != nil {
		return
	}

	column, err := gtk.TreeViewColumnNewWithAttribute("Geometry", typeCell, "text", GEO_TYPE)
	if err != nil {
		return
	}
	geoTree.geoView.AppendColumn(column)

	column, err = gtk.TreeViewColumnNewWithAttribute("Geo ID", typeCell, "text", GEO_NUM)
	if err != nil {
		return
	}
	column.SetResizable(true)
	geoTree.geoView.AppendColumn(column)

	nameCell, err := gtk.CellRendererTextNew()
	if err != nil {
		return
	}

	nameCell.SetProperty("editable", true)
	nameCell.Connect("edited", func(cell *gtk.CellRendererText, path, text string) {
		iter, err := geoTree.geoModel.GetIterFromString(path)
		if err != nil {
			log.Printf("Error editing geometry (%s)", err)
			return
		}

		model := geoTree.geoModel.ToTreeModel()
		geoID, err := util.ModelGetValue[int](model, iter, GEO_NUM)
		if err != nil {
			log.Printf("Error editing geometry (%s)", err)
			return
		}

		geo := geoTree.currentPage.PropMap[geoID]
		if geo == nil {
			log.Printf("Error getting geometry %d", geoID)
			return
		}

		geo.Name = text
		geoTree.geoModel.SetValue(iter, GEO_NAME, text)
		geoTree.updateGeometry(geoID, text)
	})

	column, err = gtk.TreeViewColumnNewWithAttribute("Name", nameCell, "text", GEO_NAME)
	if err != nil {
		return
	}
	geoTree.geoView.AppendColumn(column)

	geoTree.geoModel, err = gtk.TreeStoreNew(
		glib.TYPE_STRING, // GEO TYPE
		glib.TYPE_STRING, // GEO NAME
		glib.TYPE_INT,    // GEO NUM
	)
	if err != nil {
		return
	}

	geoTree.geoList, err = gtk.TreeStoreNew(
		glib.TYPE_STRING,
		glib.TYPE_STRING,
		glib.TYPE_INT,
	)
	if err != nil {
		return
	}

	geoTree.geoView.SetModel(geoTree.geoModel)

	geoTree.geoView.Connect("row-activated",
		func(tree *gtk.TreeView, path *gtk.TreePath, column *gtk.TreeViewColumn) {
			iter, err := geoTree.geoModel.GetIter(path)
			if err != nil {
				log.Printf("Error sending prop to editor (%s)", err)
				return
			}

			model := &geoTree.geoModel.TreeModel
			propID, err := util.ModelGetValue[int](model, iter, GEO_NUM)
			if err != nil {
				log.Printf("Error sending prop to editor (%s)", err)
				return
			}

			propToEditor(propID)
		})

	return
}

func (geoTree *GeoTree) updateGeometry(geoID int, name string) {
	iter, ok := geoTree.geoList.GetIterFirst()
	model := geoTree.geoList.ToTreeModel()

	for ok {
		currentID, err := util.ModelGetValue[int](model, iter, GEO_NUM)
		if err != nil {
			log.Printf("Error getting geometry (%s)", err)
			ok = model.IterNext(iter)
			continue
		}

		if currentID == geoID {
			geoTree.geoList.SetValue(iter, GEO_NAME, name)
		}

		ok = model.IterNext(iter)
	}
}

func (geoTree *GeoTree) removeGeometry(geoID int) {
	iter, ok := geoTree.geoList.GetIterFirst()
	model := geoTree.geoList.ToTreeModel()

	for ok {
		currentID, err := util.ModelGetValue[int](model, iter, GEO_NUM)
		if err != nil {
			log.Printf("Error getting geometry (%s)", err)
			ok = geoTree.geoList.IterNext(iter)
			continue
		}

		if currentID == geoID {
			geoTree.geoList.Remove(iter)
			iter, ok = geoTree.geoList.GetIterFirst()
		} else {
			ok = geoTree.geoList.IterNext(iter)
		}
	}
}

func geometryToTreeView(page *pages.Page, tempView *GeoTree, iter *gtk.TreeIter, propID int) {
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
func (geoTree *GeoTree) ImportGeometry(page *pages.Page) (err error) {
	geometryToTreeView(page, geoTree, nil, 0)

	return
}

func (geoTree *GeoTree) ExportGeometry(page *pages.Page) (err error) {
	return
}

func (geoTree *GeoTree) AddGeoRow(iter *gtk.TreeIter, name, propName string, propNum int) {
	geoTree.geoModel.SetValue(iter, GEO_TYPE, propName)
	geoTree.geoModel.SetValue(iter, GEO_NAME, name)
	geoTree.geoModel.SetValue(iter, GEO_NUM, propNum)

	newIter := geoTree.geoList.Append(nil)
	geoTree.geoList.SetValue(newIter, GEO_TYPE, propName)
	geoTree.geoList.SetValue(newIter, GEO_NAME, name)
	geoTree.geoList.SetValue(newIter, GEO_NUM, propNum)
}

func (geoTree *GeoTree) Clear() {
	geoTree.geoList.Clear()
	geoTree.geoModel.Clear()
}
