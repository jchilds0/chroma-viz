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
	ATTR_SELECT_TYPE = iota
	ATTR_SELECT_NAME
)

const (
	FRAME_KEY_TYPE = iota
	FRAME_KEY_ID
	FRAME_GEOMETRY
	FRAME_GEOMETRY_ID
	FRAME_ATTR_TYPE
	FRAME_ATTR_NAME
)

type KeyTree struct {
	frameID   int
	userFrame map[int]templates.UserFrame
	bindFrame map[int]templates.BindFrame
	setFrame  map[int]templates.SetFrame

	keyframeModel map[int]*gtk.ListStore
	keyframeView  map[int]*gtk.TreeView

	keyTypeSelect *gtk.ComboBoxText

	keyFrameList *gtk.ListStore

	keyGeoList   *gtk.ListStore
	keyGeoSelect *gtk.ComboBox

	keyAttrList   *gtk.ListStore
	keyAttrSelect *gtk.ComboBox

	keyFrameStack *gtk.StackSidebar
	editor        *templates.Editor
}

func NewKeyframeTree(editor *templates.Editor, keyType *gtk.ComboBoxText,
	geoModel *gtk.ListStore, frameModel *gtk.ListStore,
	keyGeo, keyAttr *gtk.ComboBox, sideBar *gtk.StackSidebar) (keyTree *KeyTree) {
	keyTree = &KeyTree{
		frameID:       1,
		keyTypeSelect: keyType,
		keyGeoSelect:  keyGeo,
		keyAttrSelect: keyAttr,
		keyFrameStack: sideBar,
		editor:        editor,
	}

	keyTree.userFrame = make(map[int]templates.UserFrame, 20)
	keyTree.setFrame = make(map[int]templates.SetFrame, 20)
	keyTree.bindFrame = make(map[int]templates.BindFrame, 20)

	keyTree.keyframeModel = make(map[int]*gtk.ListStore, 20)
	keyTree.keyframeView = make(map[int]*gtk.TreeView, 20)

	geoCell, err := gtk.CellRendererTextNew()
	if err != nil {
		log.Fatal(err)
	}

	{

		keyTree.keyGeoList = geoModel
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
		keyTree.keyAttrSelect.CellLayout.AddAttribute(geoCell, "text", ATTR_SELECT_NAME)
		keyTree.keyAttrSelect.SetActive(ATTR_SELECT_NAME)
		keyTree.keyAttrSelect.SetModel(keyTree.keyAttrList)

	}

	return
}

func (keyTree *KeyTree) SelectedFrame() (frameNum int, err error) {
	stack := keyTree.keyFrameStack.GetStack()
	frameString := stack.GetVisibleChildName()

	frameNum, err = strconv.Atoi(frameString)
	if err != nil {
		err = fmt.Errorf("Error getting frame num: %s", err)
		return
	}

	return
}

func (keyTree *KeyTree) SelectedGeometry() (geoID int, geoName string, err error) {
	iter, err := keyTree.keyGeoSelect.GetActiveIter()
	if err != nil {
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

func (keyTree *KeyTree) SelectedAttribute() (attrType string, err error) {
	iter, err := keyTree.keyAttrSelect.GetActiveIter()
	if err != nil {
		return
	}

	attrType, err = util.ModelGetValue[string](keyTree.keyAttrList.ToTreeModel(), iter, ATTR_SELECT_TYPE)
	if err != nil {
		return
	}

	return
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
		glib.TYPE_STRING, // Keyframe Type
		glib.TYPE_INT,    // Keyframe ID
		glib.TYPE_STRING, // Geometry Name
		glib.TYPE_INT,    // Geometry Num
		glib.TYPE_STRING, // Geometry Attr Type
		glib.TYPE_STRING, // Geometry Attr Name
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

	view.Connect("row-activated",
		func(tree *gtk.TreeView, path *gtk.TreePath, column *gtk.TreeViewColumn) {
			iter, err := model.GetIter(path)
			if err != nil {
				log.Printf("Error sending keyframe to editor (%s)", err)
				return
			}

			frameID, err := util.ModelGetValue[int](model.ToTreeModel(), iter, FRAME_KEY_ID)
			if err != nil {
				log.Printf("Error sending keyframe to editor (%s)", err)
				return
			}

			frameType, err := util.ModelGetValue[string](model.ToTreeModel(), iter, FRAME_KEY_TYPE)
			if err != nil {
				log.Printf("Error sending keyframe to editor (%s)", err)
				return
			}

			keyTree.editor.CurrentKeyID = frameID
			switch frameType {
			case templates.SET_FRAME:
				frame, ok := keyTree.setFrame[keyTree.editor.CurrentKeyID]
				if !ok {
					log.Printf("Missing keyframe %d", keyTree.editor.CurrentKeyID)
					break
				}

				keyTree.editor.SetFrame(frame)

			case templates.BIND_FRAME:
				frame, ok := keyTree.bindFrame[keyTree.editor.CurrentKeyID]
				if !ok {
					log.Printf("Missing keyframe %d", keyTree.editor.CurrentKeyID)
					break
				}

				keyTree.editor.BindFrame(frame)
			case templates.USER_FRAME:
			}
		})

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

		column, err = gtk.TreeViewColumnNewWithAttribute("Attribute", attrCell, "text", FRAME_ATTR_NAME)
		if err != nil {
			return
		}

		column.SetResizable(true)
		view.AppendColumn(column)

	}

	{
		// Keyframe Type
		var frameCell *gtk.CellRendererText
		var column *gtk.TreeViewColumn

		frameCell, err = gtk.CellRendererTextNew()
		if err != nil {
			return
		}

		column, err = gtk.TreeViewColumnNewWithAttribute("Keyframe Type", frameCell, "text", FRAME_KEY_TYPE)
		if err != nil {
			return
		}

		column.SetResizable(true)
		view.AppendColumn(column)

	}

	stack := keyTree.keyFrameStack.GetStack()
	name := fmt.Sprintf("   Frame %d   ", frameNum)
	stack.AddTitled(keyTree.keyframeView[frameNum], strconv.Itoa(frameNum), name)

	return
}

func (keyTree *KeyTree) RemoveFrame() (err error) {
	frameNum, err := keyTree.SelectedFrame()
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

	attrType, err := keyTree.SelectedAttribute()
	if err != nil {
		err = fmt.Errorf("Error getting attribute: %s", err)
		return
	}

	frameType := keyTree.keyTypeSelect.GetActiveText()

	frameNum, err := keyTree.SelectedFrame()
	if err != nil {
		return
	}

	frame := templates.NewKeyFrame(frameNum, geoID, attrType, false)
	frame.Type = frameType

	frameID, err := keyTree.addGeometryRow(geoName, *frame)
	if err != nil {
		return
	}

	switch frameType {
	case templates.USER_FRAME:
		keyTree.userFrame[frameID] = *templates.NewUserFrame(*frame)
	case templates.SET_FRAME:
		keyTree.setFrame[frameID] = *templates.NewSetFrame(*frame, 0)
	case templates.BIND_FRAME:
		keyTree.bindFrame[frameID] = *templates.NewBindFrame(*frame, *templates.NewKeyFrame(0, 0, "", false))
	}

	return
}

func (keyTree *KeyTree) RemoveKeyframe() (err error) {
	frameNum, err := keyTree.SelectedFrame()
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

func (keyTree *KeyTree) ImportKeyframes(temp *templates.Template) (err error) {
	var frameID int

	// add frames
	for range temp.MaxKeyframe() {
		err = keyTree.AddFrame()
		if err != nil {
			return
		}
	}

	for _, frame := range temp.UserFrame {
		geo := temp.Geos[frame.GeoID]

		frameID, err = keyTree.addGeometryRow(geo.Name, frame.Keyframe)
		if err != nil {
			return
		}

		keyTree.userFrame[frameID] = frame
	}

	for _, frame := range temp.BindFrame {
		geo := temp.Geos[frame.GeoID]

		frameID, err = keyTree.addGeometryRow(geo.Name, frame.Keyframe)
		if err != nil {
			return
		}

		keyTree.bindFrame[frameID] = frame
	}

	for _, frame := range temp.SetFrame {
		geo := temp.Geos[frame.GeoID]

		frameID, err = keyTree.addGeometryRow(geo.Name, frame.Keyframe)
		if err != nil {
			return
		}

		keyTree.setFrame[frameID] = frame
	}

	return
}

func (keyTree *KeyTree) addGeometryRow(geoName string, frame templates.Keyframe) (id int, err error) {
	model := keyTree.keyframeModel[frame.FrameNum]
	if model == nil {
		err = fmt.Errorf("Keyframe %d model does not exist", frame.FrameNum)
		return
	}

	iter := model.Append()
	id = keyTree.frameID
	model.SetValue(iter, FRAME_GEOMETRY_ID, frame.GeoID)
	model.SetValue(iter, FRAME_GEOMETRY, geoName)
	model.SetValue(iter, FRAME_ATTR_TYPE, frame.GeoAttr)
	model.SetValue(iter, FRAME_ATTR_NAME, geometry.Attrs[frame.GeoAttr])
	model.SetValue(iter, FRAME_KEY_TYPE, frame.Type)
	model.SetValue(iter, FRAME_KEY_ID, keyTree.frameID)

	keyTree.frameID++
	return
}

func (keyTree *KeyTree) ExportKeyframes(temp *templates.Template) (err error) {
	temp.SetFrame = make([]templates.SetFrame, 0, len(keyTree.setFrame))
	temp.UserFrame = make([]templates.UserFrame, 0, len(keyTree.userFrame))
	temp.BindFrame = make([]templates.BindFrame, 0, len(keyTree.bindFrame))

	for _, frame := range keyTree.setFrame {
		temp.SetFrame = append(temp.SetFrame, frame)
	}

	for _, frame := range keyTree.bindFrame {
		temp.BindFrame = append(temp.BindFrame, frame)
	}

	for _, frame := range keyTree.userFrame {
		temp.UserFrame = append(temp.UserFrame, frame)
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

	stack := keyTree.keyFrameStack.GetStack()

	for k, model := range keyTree.keyframeView {
		if model == nil {
			continue
		}

		stack.Remove(model)
		delete(keyTree.keyframeView, k)
	}
}
