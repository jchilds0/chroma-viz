package main

import (
	"chroma-viz/library/geometry"
	"chroma-viz/library/templates"
	"chroma-viz/library/util"
	"fmt"
	"log"
	"strconv"

	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
)

const (
	FRAME_GEO_ATTR = iota
	FRAME_ATTR_NAME
)

const (
	FRAME_GEOMETRY = iota
	FRAME_GEOMETRY_ID
	FRAME_ATTR_TYPE
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
	keyframeModel map[int]*gtk.ListStore
	keyframeView  map[int]*gtk.TreeView

	keyGeoList   *gtk.ListStore
	keyGeoSelect *gtk.ComboBox

	keyAttrList   *gtk.ListStore
	keyAttrSelect *gtk.ComboBox

	keyFrameStack *gtk.StackSidebar
}

func NewKeyframeTree(keyGeo, keyAttr *gtk.ComboBox, sideBar *gtk.StackSidebar) (keyTree *KeyTree) {
	keyTree = &KeyTree{
		keyGeoSelect:  keyGeo,
		keyAttrSelect: keyAttr,
		keyFrameStack: sideBar,
	}

	keyTree.keyframeModel = make(map[int]*gtk.ListStore, 20)
	keyTree.keyframeView = make(map[int]*gtk.TreeView, 20)

	geoCell, err := gtk.CellRendererTextNew()
	if err != nil {
		log.Fatal(err)
	}

	{

		var err error
		keyTree.keyGeoList, err = gtk.ListStoreNew(glib.TYPE_STRING, glib.TYPE_STRING, glib.TYPE_INT)
		if err != nil {
			log.Fatal(err)
		}

		keyTree.keyGeoSelect.PackStart(geoCell, true)
		keyTree.keyGeoSelect.CellLayout.AddAttribute(geoCell, "text", GEO_NAME)
		keyTree.keyGeoSelect.SetActive(GEO_NAME)
		keyTree.keyGeoSelect.SetModel(keyTree.keyGeoList)

	}

	{

		var err error
		keyTree.keyAttrList, err = gtk.ListStoreNew(glib.TYPE_STRING, glib.TYPE_STRING)
		if err != nil {
			log.Fatal(err)
		}

		keyTree.keyAttrSelect.PackStart(geoCell, true)
		keyTree.keyAttrSelect.CellLayout.AddAttribute(geoCell, "text", FRAME_ATTR_NAME)
		keyTree.keyAttrSelect.SetActive(GEO_NAME)
		keyTree.keyAttrSelect.SetModel(keyTree.keyAttrList)

	}

	return
}

func (keyTree *KeyTree) SelectedGeometry() (geoID int, geoName string, err error) {
	iter, err := keyTree.keyGeoSelect.GetActiveIter()
	if err != nil {
		log.Printf("No geometry selected")
		return
	}

	geoID, err = util.ModelGetValue[int](keyTree.keyGeoList.ToTreeModel(), iter, GEO_NUM)
	if err != nil {
		return
	}

	geoName, err = util.ModelGetValue[string](keyTree.keyGeoList.ToTreeModel(), iter, GEO_NAME)
	if err != nil {
		return
	}

	return
}

func (keyTree *KeyTree) SelectedAttribute() (attrType, attr string, err error) {
	iter, err := keyTree.keyAttrSelect.GetActiveIter()
	if err != nil {
		log.Printf("No attribute selected")
		return
	}

	attrType, err = util.ModelGetValue[string](keyTree.keyAttrList.ToTreeModel(), iter, FRAME_GEO_ATTR)
	if err != nil {
		return
	}

	attr, err = util.ModelGetValue[string](keyTree.keyAttrList.ToTreeModel(), iter, FRAME_ATTR_NAME)
	if err != nil {
		return
	}

	return
}

var keyframeAttrs = map[string]bool{
	"rel_x":        true,
	"rel_y":        true,
	"width":        true,
	"height":       true,
	"rounding":     true,
	"start_angle":  true,
	"end_angle":    true,
	"inner_radius": true,
	"outer_radius": true,
}

func (keyTree *KeyTree) UpdateAttrList(geoType string) {
	keyTree.keyAttrList.Clear()

	attrs := []string{geometry.ATTR_REL_X, geometry.ATTR_REL_Y}

	switch geoType {
	case geometry.GEO_RECT:
		attrs = append(attrs, geometry.ATTR_WIDTH)
		attrs = append(attrs, geometry.ATTR_HEIGHT)

	case geometry.GEO_CIRCLE:
		attrs = append(attrs, geometry.ATTR_INNER_RADIUS)
		attrs = append(attrs, geometry.ATTR_OUTER_RADIUS)
		attrs = append(attrs, geometry.ATTR_START_ANGLE)
		attrs = append(attrs, geometry.ATTR_END_ANGLE)

	case geometry.GEO_TEXT:
	case geometry.GEO_IMAGE:
	case geometry.GEO_POLY:
	case geometry.GEO_TICKER:
	case geometry.GEO_CLOCK:

	default:
		log.Fatalf("Unknown geometry type %s", geoType)
	}

	for _, name := range attrs {
		iter := keyTree.keyAttrList.Append()

		keyTree.keyAttrList.SetValue(iter, FRAME_GEO_ATTR, name)
		keyTree.keyAttrList.SetValue(iter, FRAME_ATTR_NAME, geometry.Attrs[name])
	}
}

func (keyTree *KeyTree) AddGeometry(name string, propNum int) {
	newIter := keyTree.keyGeoList.Append()

	keyTree.keyGeoList.SetValue(newIter, GEO_TYPE, name)
	keyTree.keyGeoList.SetValue(newIter, GEO_NAME, name)
	keyTree.keyGeoList.SetValue(newIter, GEO_NUM, propNum)
}

func (keyTree *KeyTree) RemoveGeo(geoID int) {
	iter, ok := keyTree.keyGeoList.GetIterFirst()
	model := keyTree.keyGeoList.ToTreeModel()

	for ok {
		currentID, err := util.ModelGetValue[int](model, iter, GEO_NUM)
		if err != nil {
			log.Printf("Error getting geometry (%s)", err)
			ok = keyTree.keyGeoList.IterNext(iter)
			continue
		}

		if currentID == geoID {
			keyTree.keyGeoList.Remove(iter)
			iter, ok = keyTree.keyGeoList.GetIterFirst()
		} else {
			ok = keyTree.keyGeoList.IterNext(iter)
		}
	}
}

func (keyTree *KeyTree) AddFrame() (err error) {
	model, err := gtk.ListStoreNew(
		glib.TYPE_STRING,  // Geometry Name
		glib.TYPE_INT,     // Geometry Num
		glib.TYPE_STRING,  // Geometry Attr Type
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

	frameNum := len(keyTree.keyframeView) + 1
	for i := 1; i <= len(keyTree.keyframeView); i++ {
		if _, ok := keyTree.keyframeView[i]; ok {
			continue
		}

		frameNum = i
		break
	}

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
						log.Printf("Error toggling toggle: %s", err)
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
					log.Printf("Error editing geometry: %s", err)
					return
				}

				model.SetValue(iter, cols[i], text)
			})
		}

		column.SetExpand(true)
		view.AppendColumn(column)

	}

	stack := keyTree.keyFrameStack.GetStack()
	name := fmt.Sprintf("   Frame %d   ", frameNum)
	stack.AddTitled(keyTree.keyframeView[frameNum], strconv.Itoa(frameNum), name)

	return
}

func (keyTree *KeyTree) RemoveFrame() (err error) {
	frameNum, err := keyTree.getSelectedFrameNum()
	if err != nil {
		return
	}

	stack := keyTree.keyFrameStack.GetStack()
	stack.Remove(keyTree.keyframeView[frameNum])
	delete(keyTree.keyframeView, frameNum)

	return
}

func (keyTree *KeyTree) AddKeyframe() (err error) {
	geoID, geoName, err := keyTree.SelectedGeometry()
	if err != nil {
		err = fmt.Errorf("Error getting geo id: %s", err)
		return
	}

	attrType, attrName, err := keyTree.SelectedAttribute()
	if err != nil {
		err = fmt.Errorf("Error getting attribute: %s", err)
		return
	}

	frameNum, err := keyTree.getSelectedFrameNum()
	if err != nil {
		return
	}

	model := keyTree.keyframeModel[frameNum]
	if model == nil {
		err = fmt.Errorf("Keyframe %d model does not exist", frameNum)
		return
	}

	iter := model.Append()
	model.SetValue(iter, FRAME_GEOMETRY, geoName)
	model.SetValue(iter, FRAME_GEOMETRY_ID, geoID)
	model.SetValue(iter, FRAME_ATTR_TYPE, attrType)
	model.SetValue(iter, FRAME_ATTR, attrName)
	return
}

func (keyTree *KeyTree) RemoveKeyframe() (err error) {
	frameNum, err := keyTree.getSelectedFrameNum()
	if err != nil {
		return
	}

	model := keyTree.keyframeModel[frameNum]
	if model == nil {
		err = fmt.Errorf("Error getting selected keyframe model")
		return
	}

	view := keyTree.keyframeView[frameNum]
	if view == nil {
		err = fmt.Errorf("Error getting selected keyframe view")
		return
	}

	selection, err := view.GetSelection()
	if err != nil {
		return
	}

	_, iter, ok := selection.GetSelected()
	if !ok {
		return
	}

	model.Remove(iter)
	return
}

func (keyTree *KeyTree) getSelectedFrameNum() (frameNum int, err error) {
	stack := keyTree.keyFrameStack.GetStack()
	frameString := stack.GetVisibleChildName()

	frameNum, err = strconv.Atoi(frameString)
	if err != nil {
		err = fmt.Errorf("Error getting frame num: %s", err)
		return
	}

	return
}

func (keyTree *KeyTree) ImportKeyframes(temp *templates.Template) (err error) {
	for range temp.MaxKeyframe() {
		err = keyTree.AddFrame()
		if err != nil {
			return
		}
	}

	var iter *gtk.TreeIter
	for _, frame := range temp.UserFrame {
		model := keyTree.keyframeModel[frame.FrameNum]
		geo := temp.Geos[frame.GeoID]

		iter, err = keyTree.addGeometryRow(&frame.Keyframe, geo)
		if err != nil {
			return
		}

		model.SetValue(iter, FRAME_USER_VALUE, true)
	}

	for _, frame := range temp.BindFrame {
		model := keyTree.keyframeModel[frame.FrameNum]
		geo := temp.Geos[frame.GeoID]

		iter, err = keyTree.addGeometryRow(&frame.Keyframe, geo)
		if err != nil {
			return
		}

		model.SetValue(iter, FRAME_BIND_FRAME, frame.Bind.FrameNum)
		model.SetValue(iter, FRAME_BIND_GEO, frame.Bind.GeoID)
		model.SetValue(iter, FRAME_BIND_ATTR, frame.Bind.GeoAttr)
	}

	for _, frame := range temp.SetFrame {
		model := keyTree.keyframeModel[frame.FrameNum]
		geo := temp.Geos[frame.GeoID]

		iter, err = keyTree.addGeometryRow(&frame.Keyframe, geo)
		if err != nil {
			return
		}

		model.SetValue(iter, FRAME_VALUE, frame.Value)
	}

	return
}

func (keyTree *KeyTree) addGeometryRow(frame *templates.Keyframe, geo *geometry.Geometry) (iter *gtk.TreeIter, err error) {
	model := keyTree.keyframeModel[frame.FrameNum]
	if model == nil {
		err = fmt.Errorf("Keyframe %d model does not exist", frame.FrameNum)
		return
	}

	if geo == nil {
		err = fmt.Errorf("Error adding keyframe %d: Geometry is nil", frame.FrameNum)
		return
	}

	iter = model.Append()
	model.SetValue(iter, FRAME_GEOMETRY_ID, frame.GeoID)
	model.SetValue(iter, FRAME_GEOMETRY, geo.Name)
	model.SetValue(iter, FRAME_ATTR_TYPE, frame.GeoAttr)
	model.SetValue(iter, FRAME_ATTR, geometry.Attrs[frame.GeoAttr])
	model.SetValue(iter, FRAME_EXPAND, frame.Expand)

	return
}

func (keyTree *KeyTree) ExportKeyframes(temp *templates.Template) (err error) {
	temp.SetFrame = make([]templates.SetFrame, 0, len(temp.SetFrame))
	temp.UserFrame = make([]templates.UserFrame, 0, len(temp.UserFrame))
	temp.BindFrame = make([]templates.BindFrame, 0, len(temp.BindFrame))

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

	frame.GeoAttr, err = util.ModelGetValue[string](model.ToTreeModel(), iter, FRAME_ATTR_TYPE)
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

func (keyTree *KeyTree) UpdateGeometryName(geoID int, name string) {
	iter, ok := keyTree.keyGeoList.GetIterFirst()
	model := keyTree.keyGeoList.ToTreeModel()

	for ok {
		currentID, err := util.ModelGetValue[int](model, iter, GEO_NUM)
		if err != nil {
			log.Printf("Error getting geometry (%s)", err)
			ok = model.IterNext(iter)
			continue
		}

		if currentID == geoID {
			keyTree.keyGeoList.SetValue(iter, GEO_NAME, name)
		}

		ok = model.IterNext(iter)
	}

	for frameNum := range keyTree.keyframeModel {
		if keyTree.keyframeModel[frameNum] == nil {
			log.Printf("Missing keyframe %d model", frameNum)
			continue
		}

		updateKeys(keyTree.keyframeModel[frameNum], geoID, name)
	}
}

func (keyTree *KeyTree) Clear() {
	keyTree.keyGeoList.Clear()
	keyTree.keyAttrList.Clear()

	for k, model := range keyTree.keyframeModel {
		model.Clear()
		delete(keyTree.keyframeModel, k)
	}

	for k, model := range keyTree.keyframeView {
		if model == nil {
			continue
		}

		keyTree.keyFrameStack.Remove(model)
		delete(keyTree.keyframeView, k)
	}
}
