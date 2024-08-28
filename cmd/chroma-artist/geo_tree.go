package main

import (
	"chroma-viz/library/geometry"
	"chroma-viz/library/pages"
	"chroma-viz/library/util"
	"fmt"
	"log"

	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
)

var propNames = map[string]string{
	geometry.GEO_RECT:   "Rectangle",
	geometry.GEO_CIRCLE: "Circle",
	geometry.GEO_TEXT:   "Text",
	geometry.GEO_TICKER: "Ticker",
	geometry.GEO_CLOCK:  "Clock",
	geometry.GEO_IMAGE:  "Image",
	geometry.GEO_POLY:   "Polygon",
}

const (
	GEO_TYPE = iota
	GEO_NAME
	GEO_NUM
	GEO_NUM_COLS
)

const (
	SELECTOR_PROP_NAME = iota
	SELECTOR_GEO_NAME
)

type GeoTree struct {
	currentPage *pages.Page
	geoModel    *gtk.TreeStore
	geoView     *gtk.TreeView
	geoIter     map[int]*gtk.TreeIter
	geoSelector *gtk.ComboBox
}

func NewGeoTree(geoSelector *gtk.ComboBox, propToEditor func(propID int), editName func(geoID int, name string)) (geoTree *GeoTree, err error) {
	geoTree = &GeoTree{
		geoSelector: geoSelector,
	}

	geoTree.geoIter = make(map[int]*gtk.TreeIter)
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
			err = fmt.Errorf("Error getting geometry %d", geoID)
			return
		}

		geo.Name = text
		geoTree.geoModel.SetValue(iter, GEO_NAME, text)

		editName(geoID, text)
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

	model, err := gtk.ListStoreNew(glib.TYPE_STRING, glib.TYPE_STRING)
	if err != nil {
		return
	}

	for propName, name := range propNames {
		iter := model.Append()

		model.SetValue(iter, SELECTOR_PROP_NAME, propName)
		model.SetValue(iter, SELECTOR_GEO_NAME, name)
	}

	cell, err := gtk.CellRendererTextNew()
	if err != nil {
		return
	}

	geoSelector.SetModel(model)
	geoSelector.CellLayout.PackStart(cell, true)
	geoSelector.CellLayout.AddAttribute(cell, "text", 1)
	geoSelector.SetActive(1)

	return
}

func (geoTree *GeoTree) GetSelectedPropName() (propType string, err error) {
	iter, err := geoTree.geoSelector.GetActiveIter()
	if err != nil {
		return
	}

	model, err := geoTree.geoSelector.GetModel()
	if err != nil {
		return
	}

	propType, err = util.ModelGetValue[string](model.ToTreeModel(), iter, SELECTOR_PROP_NAME)
	return
}

func (geoTree *GeoTree) GetSelectedGeoID() (geoID int, err error) {
	selection, err := geoTree.geoView.GetSelection()
	if err != nil {
		err = fmt.Errorf("Error getting selected: %s", err)
		return
	}

	_, iter, ok := selection.GetSelected()
	if !ok {
		err = fmt.Errorf("No geometry selected")
		return
	}

	model := geoTree.geoModel.ToTreeModel()
	geoID, err = util.ModelGetValue[int](model, iter, GEO_NUM)
	return
}

func (geoTree *GeoTree) RemoveGeo(geoID int) {
	iter := geoTree.geoIter[geoID]
	geoTree.geoModel.Remove(iter)

	delete(geoTree.geoIter, geoID)
}

func geometryToTreeView(page *pages.Page, geoTree *GeoTree, propID int) {
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

		geoTree.AddGeoRow(id, propID, geo.Name, propNames[geo.PropType])
		if id == propID {
			continue
		}

		geometryToTreeView(page, geoTree, id)
	}
}

func (geoTree *GeoTree) ImportGeometry(page *pages.Page) (err error) {
	geometryToTreeView(page, geoTree, 0)

	return
}

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

func (geoTree *GeoTree) ExportGeometry(page *pages.Page) {
	model := geoTree.geoModel.ToTreeModel()

	if iter, ok := model.GetIterFirst(); ok {
		updateParentGeometry(page, model, iter, 0)
	}

	return
}

func (geoTree *GeoTree) AddGeoRow(geoID, parentID int, name, propName string) {
	parentIter := geoTree.geoIter[parentID]
	iter := geoTree.geoModel.Append(parentIter)

	geoTree.geoIter[geoID] = iter
	geoTree.geoModel.SetValue(iter, GEO_TYPE, propName)
	geoTree.geoModel.SetValue(iter, GEO_NAME, name)
	geoTree.geoModel.SetValue(iter, GEO_NUM, geoID)
}

func (geoTree *GeoTree) Clear() {
	geoTree.geoModel.Clear()

	geoTree.geoIter = make(map[int]*gtk.TreeIter)
}
