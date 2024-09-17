package main

import (
	"chroma-viz/library/geometry"
	"chroma-viz/library/templates"
	"chroma-viz/library/util"
	"fmt"
	"log"
	"strconv"

	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
)

const iconSize = 24
const addFrameIcon = "cmd/chroma-artist/add-file-icon.svg"
const removeFrameIcon = "cmd/chroma-artist/remove-file-icon.svg"

const (
	ATTR_SELECT_TYPE = iota
	ATTR_SELECT_NAME
)

type Frames struct {
	keyTypeSelect *gtk.ComboBoxText
	keyGeoList    *gtk.ListStore
	keyGeoSelect  *gtk.ComboBox
	keyAttrList   *gtk.ListStore
	keyAttrSelect *gtk.ComboBox

	keyFrames     map[int]*KeyTree
	keyFrameList  *gtk.ListStore
	keyFrameStack *gtk.StackSidebar

	editor    *templates.Editor
	actions   *gtk.Box
	keyframes *gtk.Paned
}

func NewFrames(editor *templates.Editor, geoModel *gtk.ListStore,
	frameModel *gtk.ListStore) (frames *Frames, err error) {
	frames = &Frames{
		editor:       editor,
		keyFrameList: frameModel,
		keyGeoList:   geoModel,
	}

	frames.keyFrames = make(map[int]*KeyTree, 10)

	frames.keyFrameStack, err = gtk.StackSidebarNew()
	if err != nil {
		return
	}

	err = frames.InitActions()
	if err != nil {
		return
	}

	err = frames.InitFrames()
	return
}

func (frames *Frames) InitFrames() (err error) {
	frames.keyframes, err = gtk.PanedNew(gtk.ORIENTATION_HORIZONTAL)
	if err != nil {
		return
	}

	frames.keyFrameStack, err = gtk.StackSidebarNew()
	if err != nil {
		return
	}

	stack, err := gtk.StackNew()
	if err != nil {
		return
	}

	stack.SetVExpand(true)
	frames.keyFrameStack.SetStack(stack)

	frames.keyframes.Add1(frames.keyFrameStack)
	frames.keyframes.Add2(stack)

	return
}

func (frames *Frames) SelectedFrame() (frameNum int, err error) {
	stack := frames.keyFrameStack.GetStack()
	frameString := stack.GetVisibleChildName()

	frameNum, err = strconv.Atoi(frameString)
	if err != nil {
		err = fmt.Errorf("Error getting frame num: %s", err)
		return
	}

	return
}

func (frames *Frames) SelectedGeometry() (geoID int, geoName, geoType string, err error) {
	iter, err := frames.keyGeoSelect.GetActiveIter()
	if err != nil {
		return
	}

	geoID, err = util.ModelGetValue[int](frames.keyGeoList.ToTreeModel(), iter, GEO_NUM)
	if err != nil {
		return
	}

	geoName, err = util.ModelGetValue[string](frames.keyGeoList.ToTreeModel(), iter, GEO_NAME)
	if err != nil {
		return
	}

	geoType, err = util.ModelGetValue[string](frames.keyGeoList.ToTreeModel(), iter, GEO_TYPE)
	if err != nil {
		return
	}

	return
}

func (frames *Frames) SelectedAttribute() (attrType string, err error) {
	iter, err := frames.keyAttrSelect.GetActiveIter()
	if err != nil {
		return
	}

	attrType, err = util.ModelGetValue[string](frames.keyAttrList.ToTreeModel(), iter, ATTR_SELECT_TYPE)
	if err != nil {
		return
	}

	return
}

func (frames *Frames) RemoveGeo(geoID int) {
	for _, keyframe := range frames.keyFrames {
		keyframe.RemoveGeo(geoID)
	}
}

func (frames *Frames) sendEditor(frameID int, frameType string) {
	frameNum, err := frames.SelectedFrame()
	if err != nil {
		log.Print(err)
		return
	}

	frame := frames.keyFrames[frameNum]
	if frame == nil {
		log.Printf("Missing frame %d", frameNum)
		return
	}

	frames.editor.CurrentFrameID = frameNum
	frames.editor.CurrentKeyID = frameID
	frames.editor.ClearFrame()
	switch frameType {
	case templates.SET_FRAME:
		keyframe, ok := frame.setFrame[frames.editor.CurrentKeyID]
		if !ok {
			log.Printf("Missing keyframe %d", frames.editor.CurrentKeyID)
			break
		}

		frames.editor.SetFrame(keyframe)

	case templates.BIND_FRAME:
		keyframe, ok := frame.bindFrame[frames.editor.CurrentKeyID]
		if !ok {
			log.Printf("Missing keyframe %d", frames.editor.CurrentKeyID)
			break
		}

		frames.editor.BindFrame(keyframe)
	case templates.USER_FRAME:
	}
}

func (frames *Frames) AddFrame() (err error) {
	frameNum := len(frames.keyFrames) + 1
	for i := 1; i <= len(frames.keyFrames); i++ {
		if _, ok := frames.keyFrames[i]; ok {
			continue
		}

		frameNum = i
		break
	}

	frames.keyFrames[frameNum], err = NewKeyframeTree(frames.sendEditor)
	if err != nil {
		return
	}

	stack := frames.keyFrameStack.GetStack()
	name := fmt.Sprintf("   Frame %d   ", frameNum)
	stack.AddTitled(frames.keyFrames[frameNum].window, strconv.Itoa(frameNum), name)

	iter := frames.keyFrameList.Append()
	frames.keyFrameList.SetValue(iter, 0, frameNum)
	frames.keyFrameList.SetValue(iter, 1, fmt.Sprintf("Frame %d", frameNum))

	return
}

func (frames *Frames) RemoveFrame() (err error) {
	frameNum, err := frames.SelectedFrame()
	if err != nil {
		return
	}

	stack := frames.keyFrameStack.GetStack()
	stack.Remove(frames.keyFrames[frameNum].window)
	delete(frames.keyFrames, frameNum)

	return
}

func (frames *Frames) AddKeyframe() (err error) {
	geoID, geoName, _, err := frames.SelectedGeometry()
	if err != nil {
		err = fmt.Errorf("Error getting geo id: %s", err)
		return
	}

	attrType, err := frames.SelectedAttribute()
	if err != nil {
		err = fmt.Errorf("Error getting attribute: %s", err)
		return
	}

	frameType := frames.keyTypeSelect.GetActiveText()

	frameNum, err := frames.SelectedFrame()
	if err != nil {
		return
	}

	frames.keyFrames[frameNum].AddKeyframe(frameNum, geoID, geoName, attrType, frameType)
	return
}

func (frames *Frames) RemoveKeyframe() (err error) {
	frameNum, err := frames.SelectedFrame()
	if err != nil {
		return
	}

	frame := frames.keyFrames[frameNum]
	if frame == nil {
		err = fmt.Errorf("Error getting selected frame")
		return
	}

	err = frame.RemoveKeyframe()
	return
}

func (frames *Frames) ImportKeyframes(temp *templates.Template) (err error) {
	// add frames
	for range temp.MaxKeyframe() {
		err = frames.AddFrame()
		if err != nil {
			return
		}
	}

	for _, frame := range temp.UserFrame {
		geo := temp.Geos[frame.GeoID]
		if geo == nil {
			continue
		}

		frames.keyFrames[frame.FrameNum].ImportUserFrame(*geo, frame)
	}

	for _, frame := range temp.BindFrame {
		geo := temp.Geos[frame.GeoID]
		if geo == nil {
			continue
		}

		frames.keyFrames[frame.FrameNum].ImportBindFrame(*geo, frame)
	}

	for _, frame := range temp.SetFrame {
		geo := temp.Geos[frame.GeoID]
		if geo == nil {
			continue
		}

		frames.keyFrames[frame.FrameNum].ImportSetFrame(*geo, frame)
	}

	return
}

func (frames *Frames) ExportKeyframes(temp *templates.Template) (err error) {
	temp.SetFrame = make([]templates.SetFrame, 0, 128)
	temp.UserFrame = make([]templates.UserFrame, 0, 128)
	temp.BindFrame = make([]templates.BindFrame, 0, 128)

	for _, frame := range frames.keyFrames {
		frame.ExportKeyframes(temp)
	}

	return
}

func (frames *Frames) UpdateGeometryName(geoID int, name string) {
	for _, frame := range frames.keyFrames {
		frame.updateKeys(geoID, name)
	}
}

func (frames *Frames) Clear() {
	frames.keyAttrList.Clear()
	frames.keyFrameList.Clear()

	stack := frames.keyFrameStack.GetStack()
	for k, frame := range frames.keyFrames {
		stack.Remove(frame.window)
		delete(frames.keyFrames, k)
	}
}

func (frames *Frames) InitActions() (err error) {
	geoCell, err := gtk.CellRendererTextNew()
	if err != nil {
		return
	}

	// action box
	frames.actions, err = gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 0)
	if err != nil {
		return
	}

	var (
		button *gtk.Button
		buf    *gdk.Pixbuf
		img    *gtk.Image
	)

	{
		// Add Frame button
		button, err = gtk.ButtonNew()
		if err != nil {
			return
		}

		button.SetTooltipText("Add Frame")

		buf, err = gdk.PixbufNewFromFileAtSize(addFrameIcon, iconSize, iconSize)
		if err != nil {
			return
		}

		img, err = gtk.ImageNewFromPixbuf(buf)
		if err != nil {
			return
		}
		button.SetImage(img)

		button.Connect("clicked", func() {
			err := frames.AddFrame()
			if err != nil {
				log.Printf("Error adding frame: %s", err)
			}
		})

		frames.actions.PackStart(button, false, false, 10)
	}

	{
		// Remove Frame Button
		button, err = gtk.ButtonNew()
		if err != nil {
			return
		}

		button.SetTooltipText("Remove Frame")

		buf, err = gdk.PixbufNewFromFileAtSize(removeFrameIcon, iconSize, iconSize)
		if err != nil {
			return
		}

		img, err = gtk.ImageNewFromPixbuf(buf)
		if err != nil {
			return
		}

		button.SetImage(img)

		button.Connect("clicked", func() {
			err := frames.RemoveFrame()
			if err != nil {
				log.Printf("Error removing frame: %s", err)
			}
		})

		frames.actions.PackStart(button, false, false, 10)
	}

	{
		// Geometry Label
		var label *gtk.Label
		label, err = gtk.LabelNew("Geometry")
		label.SetWidthChars(12)

		frames.actions.PackStart(label, false, false, 0)
	}

	{
		// Geometry selector
		frames.keyGeoSelect, err = gtk.ComboBoxNewWithModel(frames.keyGeoList)
		if err != nil {
			return
		}

		frames.keyGeoSelect.Connect("changed", func() {
			_, _, geoType, err := frames.SelectedGeometry()
			if err != nil {
				return
			}

			geometry.UpdateAttrList(frames.keyAttrList, geoType)
		})

		frames.keyGeoSelect.PackStart(geoCell, true)
		frames.keyGeoSelect.CellLayout.AddAttribute(geoCell, "text", GEO_NAME)
		frames.keyGeoSelect.SetActive(GEO_NAME)
		frames.keyGeoSelect.SetModel(frames.keyGeoList)

		frames.actions.PackStart(frames.keyGeoSelect, false, false, 0)
	}

	{
		// Attribute Label
		var label *gtk.Label
		label, err = gtk.LabelNew("Attribute")
		label.SetWidthChars(12)

		frames.actions.PackStart(label, false, false, 0)
	}

	{
		// Attribute selector
		frames.keyAttrList, err = gtk.ListStoreNew(glib.TYPE_STRING, glib.TYPE_STRING)
		if err != nil {
			return
		}

		frames.keyAttrSelect, err = gtk.ComboBoxNewWithModel(frames.keyAttrList)
		if err != nil {
			return
		}

		frames.keyAttrSelect.PackStart(geoCell, true)
		frames.keyAttrSelect.CellLayout.AddAttribute(geoCell, "text", ATTR_SELECT_NAME)
		frames.keyAttrSelect.SetActive(ATTR_SELECT_NAME)
		frames.keyAttrSelect.SetModel(frames.keyAttrList)

		frames.actions.PackStart(frames.keyAttrSelect, false, false, 10)
	}

	{
		// Keyframe Type
		frames.keyTypeSelect, err = gtk.ComboBoxTextNew()
		if err != nil {
			return
		}

		frames.keyTypeSelect.InsertText(1, "UserFrame")
		frames.keyTypeSelect.InsertText(2, "SetFrame")
		frames.keyTypeSelect.InsertText(3, "BindFrame")

		frames.actions.PackStart(frames.keyTypeSelect, false, false, 10)
	}

	{
		// Add Keyframe button
		button, err = gtk.ButtonNew()
		if err != nil {
			return
		}

		img, err = gtk.ImageNewFromIconName("list-add", 3)
		if err != nil {
			return
		}

		button.SetImage(img)
		button.SetTooltipText("Add KeyFrame")

		button.Connect("clicked", func() {
			err := frames.AddKeyframe()
			if err != nil {
				log.Printf("Error adding keyframe: %s", err)
			}
		})

		frames.actions.PackStart(button, false, false, 10)
	}

	{
		// Remove keyframe button
		button, err = gtk.ButtonNew()
		if err != nil {
			return
		}

		img, err = gtk.ImageNewFromIconName("list-remove", 3)
		if err != nil {
			return
		}

		button.SetImage(img)
		button.SetTooltipText("Remove KeyFrame")

		button.Connect("clicked", func() {
			err := frames.RemoveKeyframe()
			if err != nil {
				log.Printf("Error removing keyframe: %s", err)
			}
		})

		frames.actions.PackStart(button, false, false, 10)
	}

	return
}
