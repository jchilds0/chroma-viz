package main

import (
	"chroma-viz/library/geometry"
	"chroma-viz/library/templates"
	"chroma-viz/library/util"
	"fmt"
	"log"

	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
)

const (
	GEO_TYPE = iota
	GEO_NAME
	GEO_NUM
	GEO_VISIBLE
	GEO_NUM_COLS
)

const (
	SELECTOR_GEO_NAME = iota
)

type GeoTree struct {
	geoModel *gtk.TreeStore
	geoList  *gtk.ListStore

	geoView     *gtk.TreeView
	geoTreeIter map[int]*gtk.TreeIter
	geoListIter map[int]*gtk.TreeIter
	geoSelector *gtk.ComboBox
}

func NewGeoTree(geoSelector *gtk.ComboBox, geoList *gtk.ListStore, geoToEditor func(geoID int), editName func(geoID int, name string)) (geoTree *GeoTree, err error) {
	geoTree = &GeoTree{
		geoSelector: geoSelector,
		geoList:     geoList,
	}

	geoTree.geoTreeIter = make(map[int]*gtk.TreeIter)
	geoTree.geoListIter = make(map[int]*gtk.TreeIter)

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

	column, err = gtk.TreeViewColumnNewWithAttribute("Geometry ID", typeCell, "text", GEO_NUM)
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

		geoTree.geoModel.SetValue(iter, GEO_NAME, text)

		iter, ok := geoTree.geoList.GetIterFirst()
		model = geoTree.geoList.ToTreeModel()

		for ok {
			currentID, err := util.ModelGetValue[int](model, iter, GEO_NUM)
			if err != nil {
				log.Printf("Error getting geometry (%s)", err)
				ok = model.IterNext(iter)
				continue
			}

			if currentID == geoID {
				geoTree.geoList.SetValue(iter, GEO_NAME, text)
			}

			ok = model.IterNext(iter)
		}

		editName(geoID, text)
	})

	column, err = gtk.TreeViewColumnNewWithAttribute("Name", nameCell, "text", GEO_NAME)
	if err != nil {
		return
	}
	geoTree.geoView.AppendColumn(column)

	toggleCell, err := gtk.CellRendererToggleNew()
	if err != nil {
		return
	}

	toggleCell.SetProperty("activatable", true)
	toggleCell.Connect("toggled",
		func(cell *gtk.CellRendererToggle, path string) {
			iter, err := geoTree.geoModel.GetIterFromString(path)
			if err != nil {
				log.Printf("Error toggling toggle (%s)", err)
				return
			}

			state, err := util.ModelGetValue[bool](geoTree.geoModel.ToTreeModel(), iter, GEO_VISIBLE)
			if err != nil {
				log.Printf("Error toggling toggle: %s", err)
				return
			}

			geoTree.geoModel.SetValue(iter, GEO_VISIBLE, !state)
		})

	column, err = gtk.TreeViewColumnNewWithAttribute("Visible", toggleCell, "active", GEO_VISIBLE)
	if err != nil {
		return
	}
	geoTree.geoView.AppendColumn(column)

	geoTree.geoModel, err = gtk.TreeStoreNew(
		glib.TYPE_STRING,  // GEO TYPE
		glib.TYPE_STRING,  // GEO NAME
		glib.TYPE_INT,     // GEO NUM
		glib.TYPE_BOOLEAN, // GEO VISIBILE
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
			geoID, err := util.ModelGetValue[int](model, iter, GEO_NUM)
			if err != nil {
				log.Printf("Error sending prop to editor (%s)", err)
				return
			}

			geoToEditor(geoID)
		})

	model, err := gtk.ListStoreNew(glib.TYPE_STRING)
	if err != nil {
		return
	}

	model.SetValue(model.Append(), SELECTOR_GEO_NAME, geometry.GEO_RECT)
	model.SetValue(model.Append(), SELECTOR_GEO_NAME, geometry.GEO_CIRCLE)
	model.SetValue(model.Append(), SELECTOR_GEO_NAME, geometry.GEO_TEXT)
	model.SetValue(model.Append(), SELECTOR_GEO_NAME, geometry.GEO_IMAGE)
	model.SetValue(model.Append(), SELECTOR_GEO_NAME, geometry.GEO_POLY)
	model.SetValue(model.Append(), SELECTOR_GEO_NAME, geometry.GEO_LIST)
	model.SetValue(model.Append(), SELECTOR_GEO_NAME, geometry.GEO_CLOCK)

	cell, err := gtk.CellRendererTextNew()
	if err != nil {
		return
	}

	geoSelector.SetModel(model)
	geoSelector.CellLayout.PackStart(cell, true)
	geoSelector.CellLayout.AddAttribute(cell, "text", SELECTOR_GEO_NAME)
	geoSelector.SetActive(1)

	return
}

func (geoTree *GeoTree) GetSelectedGeoName() (geoName string, err error) {
	iter, err := geoTree.geoSelector.GetActiveIter()
	if err != nil {
		return
	}

	model, err := geoTree.geoSelector.GetModel()
	if err != nil {
		return
	}

	geoName, err = util.ModelGetValue[string](model.ToTreeModel(), iter, SELECTOR_GEO_NAME)
	return
}

func (geoTree *GeoTree) GetSelectedGeometry() (iter *gtk.TreeIter, err error) {
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
	return
}

func (geoTree *GeoTree) RemoveGeo(geoID int) {
	iter, ok := geoTree.geoTreeIter[geoID]
	if !ok {
		return
	}

	delete(geoTree.geoTreeIter, geoID)
	geoTree.geoModel.Remove(iter)

	iter, ok = geoTree.geoListIter[geoID]
	if !ok {
		return
	}

	geoTree.geoList.Remove(iter)
	delete(geoTree.geoListIter, geoID)
}

func geometryToTreeView(temp *templates.Template, geoTree *GeoTree, parentID int) {
	for geoID, geo := range temp.Geos {
		if geo.Parent.Value != parentID {
			continue
		}

		geoTree.AddGeoRow(geoID, parentID, geo.Name, geo.GeoType, geo.Visible)
		if geoID == parentID {
			continue
		}

		geometryToTreeView(temp, geoTree, geoID)
	}
}

func (geoTree *GeoTree) ImportGeometry(temp *templates.Template) (err error) {
	geometryToTreeView(temp, geoTree, 0)
	return
}

func updateGeometry(temp *templates.Template, model *gtk.TreeModel, iter *gtk.TreeIter, parentID int) {
	nextIterExists := true
	for nextIterExists {
		geoID, err := util.ModelGetValue[int](model, iter, GEO_NUM)
		if err != nil {
			log.Printf("Error getting geometry id: %s", err.Error())
			nextIterExists = model.IterNext(iter)
			continue
		}

		name, err := util.ModelGetValue[string](model, iter, GEO_NAME)
		if err != nil {
			log.Printf("Error getting geometry name: %s", err.Error())
			nextIterExists = model.IterNext(iter)
			continue
		}

		visible, err := util.ModelGetValue[bool](model, iter, GEO_VISIBLE)
		if err != nil {
			log.Printf("Error getting geometry name: %s", err.Error())
			nextIterExists = model.IterNext(iter)
			continue
		}

		geo := temp.Geos[geoID]
		if geo == nil {
			log.Printf("Error: geometry %d is nil", geoID)
			nextIterExists = model.IterNext(iter)
			continue
		}

		geo.Name = name
		geo.Parent.Value = parentID
		geo.Visible = visible

		var childIter gtk.TreeIter
		if ok := model.IterChildren(iter, &childIter); ok {
			updateGeometry(temp, model, &childIter, geoID)
		}

		nextIterExists = model.IterNext(iter)
	}
}

func (geoTree *GeoTree) ExportGeometry(temp *templates.Template) {
	model := geoTree.geoModel.ToTreeModel()

	if iter, ok := model.GetIterFirst(); ok {
		updateGeometry(temp, model, iter, 0)
	}

	return
}

func (geoTree *GeoTree) AddGeoRow(geoID, parentID int, geoName, geoType string, geoVisible bool) {
	parentIter := geoTree.geoTreeIter[parentID]
	iter := geoTree.geoModel.Append(parentIter)

	geoTree.geoTreeIter[geoID] = iter
	geoTree.geoModel.SetValue(iter, GEO_TYPE, geoType)
	geoTree.geoModel.SetValue(iter, GEO_NAME, geoName)
	geoTree.geoModel.SetValue(iter, GEO_NUM, geoID)
	geoTree.geoModel.SetValue(iter, GEO_VISIBLE, geoVisible)

	iter = geoTree.geoList.Append()

	geoTree.geoListIter[geoID] = iter
	geoTree.geoList.SetValue(iter, GEO_TYPE, geoType)
	geoTree.geoList.SetValue(iter, GEO_NAME, geoName)
	geoTree.geoList.SetValue(iter, GEO_NUM, geoID)
	geoTree.geoList.SetValue(iter, GEO_VISIBLE, geoVisible)
}

func (geoTree *GeoTree) Clear() {
	geoTree.geoModel.Clear()
	geoTree.geoList.Clear()

	geoTree.geoTreeIter = make(map[int]*gtk.TreeIter)
}
