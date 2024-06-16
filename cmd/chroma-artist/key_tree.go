package main

import (
	"chroma-viz/library/templates"
	"chroma-viz/library/util"
	"fmt"
	"log"
	"strconv"

	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
)

const (
	FRAME_GEOMETRY = iota
	FRAME_GEOMETRY_ID
	FRAME_ATTR
	FRAME_VALUE
	FRAME_USER_VALUE
	FRAME_EXPAND
	FRAME_BIND_FRAME
	FRAME_BIND_GEO
	FRAME_BIND_ATTR
	FRAME_NUM_COLS
)

type KeyTree struct {
	nextFrame     int
	keyframeModel map[int]*gtk.ListStore
	keyframeView  map[int]*gtk.TreeView
}

func NewKeyframeTree() (keyTree *KeyTree) {
	keyTree = &KeyTree{
		nextFrame: 1,
	}

	keyTree.keyframeModel = make(map[int]*gtk.ListStore, 20)
	keyTree.keyframeView = make(map[int]*gtk.TreeView, 20)

	return
}

func (keyTree *KeyTree) AddFrame() (frameNum int, err error) {
	model, err := gtk.ListStoreNew(
		glib.TYPE_STRING,  // Geometry Name
		glib.TYPE_INT,     // Geometry Num
		glib.TYPE_STRING,  // Geometry Attr
		glib.TYPE_INT,     // Value Entry
		glib.TYPE_BOOLEAN, // User Value Selector
		glib.TYPE_BOOLEAN, // Expand Children
		glib.TYPE_STRING,  // Derived Value Frame
		glib.TYPE_STRING,  // Derived Value Geo
		glib.TYPE_STRING,  // Derived Value Attr
	)
	if err != nil {
		return
	}

	view, err := gtk.TreeViewNew()
	if err != nil {
		return
	}
	view.SetReorderable(true)
	view.SetModel(model)
	view.SetVisible(true)

	frameNum = keyTree.nextFrame
	keyTree.nextFrame++

	keyTree.keyframeModel[frameNum] = model
	keyTree.keyframeView[frameNum] = view

	{
		// Geometry Name
		var geoCell *gtk.CellRendererText
		var column *gtk.TreeViewColumn

		geoCell, err = gtk.CellRendererTextNew()
		if err != nil {
			return
		}

		column, err = gtk.TreeViewColumnNewWithAttribute("Geometry", geoCell, "text", FRAME_GEOMETRY)
		if err != nil {
			return
		}

		column.SetResizable(true)
		view.AppendColumn(column)

	}

	{
		// Attribute Name
		var attrCell *gtk.CellRendererText
		var column *gtk.TreeViewColumn

		attrCell, err = gtk.CellRendererTextNew()
		if err != nil {
			return
		}

		column, err = gtk.TreeViewColumnNewWithAttribute("Attribute", attrCell, "text", FRAME_ATTR)
		if err != nil {
			return
		}

		column.SetResizable(true)
		view.AppendColumn(column)

	}

	{
		// Set Value
		var valueCell *gtk.CellRendererText
		var column *gtk.TreeViewColumn

		valueCell, err = gtk.CellRendererTextNew()
		if err != nil {
			return
		}

		valueCell.SetProperty("editable", true)
		valueCell.SetProperty("xalign", 1.0)
		valueCell.SetProperty("xpad", 15)
		valueCell.Connect("edited", func(cell *gtk.CellRendererText, path, text string) {
			iter, err := model.GetIterFromString(path)
			if err != nil {
				log.Printf("Error editing keyframe (%s)", err)
				return
			}

			num, err := strconv.Atoi(text)
			if err != nil {
				log.Printf("Error editing keyframe (%s)", err)
				return
			}

			model.SetValue(iter, FRAME_VALUE, num)
		})

		column, err = gtk.TreeViewColumnNewWithAttribute("Set Value", valueCell, "text", FRAME_VALUE)
		if err != nil {
			return
		}

		column.SetResizable(true)
		view.AppendColumn(column)

	}

	{
		// Bool Value
		var toggleCell *gtk.CellRendererToggle
		var column *gtk.TreeViewColumn

		names := []string{"Expand", "User Value"}
		cols := []int{FRAME_EXPAND, FRAME_USER_VALUE}

		for i := range names {
			toggleCell, err = gtk.CellRendererToggleNew()
			if err != nil {
				return
			}

			toggleCell.SetProperty("activatable", true)
			toggleCell.Connect("toggled",
				func(cell *gtk.CellRendererToggle, path string) {
					iter, err := model.GetIterFromString(path)
					if err != nil {
						log.Printf("Error toggling toggle (%s)", err)
						return
					}

					state, err := util.ModelGetValue[bool](model.ToTreeModel(), iter, cols[i])
					if err != nil {
						log.Printf("Error toggling toggle (%s)", err)
						return
					}

					model.SetValue(iter, cols[i], !state)
				})

			column, err = gtk.TreeViewColumnNewWithAttribute(names[i], toggleCell, "active", cols[i])
			if err != nil {
				return
			}

			column.SetResizable(true)
			view.AppendColumn(column)
		}

	}

	{
		// Derived Value
		var valueText, valueCell *gtk.CellRendererText
		var column *gtk.TreeViewColumn

		column, err = gtk.TreeViewColumnNew()
		if err != nil {
			return
		}

		column.SetTitle("Value From Keyframe")

		names := []string{"Frame", "Geometry", "Attr"}
		cols := []int{FRAME_BIND_FRAME, FRAME_BIND_GEO, FRAME_BIND_ATTR}

		for i, name := range names {
			valueText, err = gtk.CellRendererTextNew()
			if err != nil {
				return
			}

			valueText.SetProperty("text", name+": ")

			valueCell, err = gtk.CellRendererTextNew()
			if err != nil {
				return
			}

			valueCell.SetProperty("editable", true)

			column.PackStart(valueText, false)
			column.PackStart(valueCell, true)

			column.AddAttribute(valueCell, "text", cols[i])

			valueCell.Connect("edited", func(cell *gtk.CellRendererText, path, text string) {
				iter, err := model.GetIterFromString(path)
				if err != nil {
					log.Printf("Error editing geometry (%s)", err)
					return
				}

				model.SetValue(iter, cols[i], text)
			})
		}

		column.SetExpand(true)
		view.AppendColumn(column)

	}

	return
}

func (keyTree *KeyTree) AddKeyframe(frameNum, geoID int, geoName, attrName string) (err error) {
	model := keyTree.keyframeModel[frameNum]
	if model == nil {
		err = fmt.Errorf("Keyframe %d model does not exist", frameNum)
		return
	}

	iter := model.Append()
	model.SetValue(iter, FRAME_GEOMETRY, geoName)
	model.SetValue(iter, FRAME_GEOMETRY_ID, geoID)
	model.SetValue(iter, FRAME_ATTR, attrName)
	return
}

func (keyTree *KeyTree) ImportKeyframes(temp *templates.Template) (err error) {
	for range temp.MaxKeyframe() {
		_, err = keyTree.AddFrame()
		if err != nil {
			return
		}
	}

	for _, frame := range temp.UserFrame {
		model := keyTree.keyframeModel[frame.FrameNum]
		if model == nil {
			err = fmt.Errorf("Keyframe %d model does not exist", frame.FrameNum)
			return
		}

		iter := model.Append()
		model.SetValue(iter, FRAME_GEOMETRY_ID, frame.GeoID)
		model.SetValue(iter, FRAME_ATTR, frame.GeoAttr)
		model.SetValue(iter, FRAME_EXPAND, frame.Expand)

		model.SetValue(iter, FRAME_USER_VALUE, true)
	}

	for _, frame := range temp.BindFrame {
		model := keyTree.keyframeModel[frame.FrameNum]
		if model == nil {
			err = fmt.Errorf("Keyframe %d model does not exist", frame.FrameNum)
			return
		}

		iter := model.Append()
		model.SetValue(iter, FRAME_GEOMETRY_ID, frame.GeoID)
		model.SetValue(iter, FRAME_ATTR, frame.GeoAttr)
		model.SetValue(iter, FRAME_EXPAND, frame.Expand)

		model.SetValue(iter, FRAME_BIND_FRAME, frame.Bind.FrameNum)
		model.SetValue(iter, FRAME_BIND_GEO, frame.Bind.GeoID)
		model.SetValue(iter, FRAME_BIND_ATTR, frame.Bind.GeoAttr)
	}

	for _, frame := range temp.SetFrame {
		model := keyTree.keyframeModel[frame.FrameNum]
		if model == nil {
			err = fmt.Errorf("Keyframe %d model does not exist", frame.FrameNum)
			return
		}

		iter := model.Append()
		model.SetValue(iter, FRAME_GEOMETRY_ID, frame.GeoID)
		model.SetValue(iter, FRAME_ATTR, frame.GeoAttr)
		model.SetValue(iter, FRAME_EXPAND, frame.Expand)

		model.SetValue(iter, FRAME_VALUE, frame.Value)
	}

	return
}

func (keyTree *KeyTree) ExportKeyframes(temp *templates.Template) (err error) {
	for frameNum := range keyTree.keyframeModel {
		err := keyTree.exportFrame(temp, frameNum)
		if err != nil {
			log.Print(err)
			continue
		}
	}

	return
}

func getKeyframeFromIter(model *gtk.ListStore, iter *gtk.TreeIter, frameNum int) (frame templates.Keyframe, err error) {
	frame.FrameNum = frameNum
	frame.GeoID, err = util.ModelGetValue[int](model.ToTreeModel(), iter, FRAME_GEOMETRY_ID)
	if err != nil {
		return
	}

	frame.GeoAttr, err = util.ModelGetValue[string](model.ToTreeModel(), iter, FRAME_ATTR)
	if err != nil {
		return
	}

	frame.Expand, err = util.ModelGetValue[bool](model.ToTreeModel(), iter, FRAME_EXPAND)
	if err != nil {
		return
	}

	return
}

func (keyTree *KeyTree) exportFrame(temp *templates.Template, frameNum int) (err error) {
	keyModel := keyTree.keyframeModel[frameNum]
	if keyModel == nil {
		err = fmt.Errorf("Missing keyframe %d model", frameNum)
		return
	}

	iter, ok := keyModel.GetIterFirst()

	var bindFrame, bindGeo, bindAttr string
	var user bool

	for ok {
		frame, err := getKeyframeFromIter(keyModel, iter, frameNum)
		if err != nil {
			log.Printf("Error getting keyframe: %s", err)
			ok = keyModel.IterNext(iter)
			continue
		}

		user, err = util.ModelGetValue[bool](keyModel.ToTreeModel(), iter, FRAME_USER_VALUE)
		if err != nil {
			log.Printf("Error getting keyframe: %s", err)
			ok = keyModel.IterNext(iter)
			continue
		}

		if user {
			keyframe := templates.NewUserFrame(frame)
			temp.UserFrame = append(temp.UserFrame, *keyframe)

			ok = keyModel.IterNext(iter)
			continue
		}

		bindFrame, err = util.ModelGetValue[string](keyModel.ToTreeModel(), iter, FRAME_BIND_FRAME)
		if err != nil {
			log.Printf("Error getting keyframe: %s", err)
			ok = keyModel.IterNext(iter)
			continue
		}

		bindGeo, err = util.ModelGetValue[string](keyModel.ToTreeModel(), iter, FRAME_BIND_GEO)
		if err != nil {
			log.Printf("Error getting keyframe: %s", err)
			ok = keyModel.IterNext(iter)
			continue
		}

		bindAttr, err = util.ModelGetValue[string](keyModel.ToTreeModel(), iter, FRAME_BIND_ATTR)
		if err != nil {
			log.Printf("Error getting keyframe: %s", err)
			ok = keyModel.IterNext(iter)
			continue
		}

		if bindFrame != "" && bindGeo != "" && bindAttr != "" {
			frameNum, _ := strconv.Atoi(bindFrame)
			geoNum, _ := strconv.Atoi(bindGeo)

			bind := templates.NewKeyFrame(frameNum, geoNum, bindAttr, false)

			keyframe := templates.NewBindFrame(frame, *bind)
			temp.BindFrame = append(temp.BindFrame, *keyframe)
			ok = keyModel.IterNext(iter)
			continue
		}

		var value int
		value, err = util.ModelGetValue[int](keyModel.ToTreeModel(), iter, FRAME_VALUE)
		keyframe := templates.NewSetFrame(frame, value)

		temp.SetFrame = append(temp.SetFrame, *keyframe)
		ok = keyModel.IterNext(iter)
		continue
	}

	return
}

func updateKeys(model *gtk.ListStore, geoID int, name string) {
	iter, ok := model.GetIterFirst()

	for ok {
		currentID, err := util.ModelGetValue[int](model.ToTreeModel(), iter, FRAME_GEOMETRY_ID)
		if err != nil {
			log.Printf("Error getting keyframe geo id (%s)", err)
			ok = model.IterNext(iter)
			continue
		}

		if currentID == geoID {
			model.SetValue(iter, FRAME_GEOMETRY, name)
		}

		ok = model.IterNext(iter)
	}
}

func (keyTree *KeyTree) UpdateKeys(geoID int, name string) {
	for frameNum := range keyTree.keyframeModel {
		if keyTree.keyframeModel[frameNum] == nil {
			log.Printf("Missing keyframe %d model", frameNum)
			continue
		}

		updateKeys(keyTree.keyframeModel[frameNum], geoID, name)
	}
}

/*
func (tempView *GeoTree) removeKeys(geoID int) {
	iter, ok := tempView.keyModel.GetIterFirst()
	model := tempView.keyModel.ToTreeModel()

	for ok {
		currentID, err := util.ModelGetValue[int](model, iter, FRAME_GEOMETRY_ID)
		if err != nil {
			log.Printf("Error getting keyframe geo id (%s)", err)
			ok = tempView.keyModel.IterNext(iter)
			continue
		}

		if currentID == geoID {
			tempView.keyModel.Remove(iter)
			iter, ok = tempView.keyModel.GetIterFirst()
		} else {
			ok = tempView.keyModel.IterNext(iter)
		}
	}
}

*/

func (keyTree *KeyTree) Clear() {
	keyTree.nextFrame = 1

	for k := range keyTree.keyframeModel {
		delete(keyTree.keyframeModel, k)
	}

	for k := range keyTree.keyframeView {
		delete(keyTree.keyframeView, k)
	}
}
