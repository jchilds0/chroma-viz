package main

import (
	"chroma-viz/library/geometry"
	"chroma-viz/library/templates"
	"chroma-viz/library/util"
	"log"

	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
)

const (
	FRAME_KEY_TYPE = iota
	FRAME_KEY_ID
	FRAME_GEOMETRY
	FRAME_GEOMETRY_ID
	FRAME_ATTR_TYPE
	FRAME_ATTR_NAME
	FRAME_EXPAND
)

type KeyTree struct {
	frameID int

	userFrame map[int]templates.UserFrame
	bindFrame map[int]templates.BindFrame
	setFrame  map[int]templates.SetFrame

	model *gtk.ListStore
	view  *gtk.TreeView
}

func NewKeyframeTree(sendEditor func(frameID int, keyType string)) (keyTree *KeyTree, err error) {
	keyTree = &KeyTree{}

	keyTree.userFrame = make(map[int]templates.UserFrame, 20)
	keyTree.setFrame = make(map[int]templates.SetFrame, 20)
	keyTree.bindFrame = make(map[int]templates.BindFrame, 20)

	keyTree.model, err = gtk.ListStoreNew(
		glib.TYPE_STRING,  // Keyframe Type
		glib.TYPE_INT,     // Keyframe ID
		glib.TYPE_STRING,  // Geometry Name
		glib.TYPE_INT,     // Geometry Num
		glib.TYPE_STRING,  // Geometry Attr Type
		glib.TYPE_STRING,  // Geometry Attr Name
		glib.TYPE_BOOLEAN, // Expand Attr
	)
	if err != nil {
		return
	}

	keyTree.view, err = gtk.TreeViewNewWithModel(keyTree.model)
	if err != nil {
		return
	}

	keyTree.view.SetReorderable(true)
	keyTree.view.SetVisible(true)

	keyTree.view.Connect("row-activated",
		func(tree *gtk.TreeView, path *gtk.TreePath, column *gtk.TreeViewColumn) {
			iter, err := keyTree.model.GetIter(path)
			if err != nil {
				log.Printf("Error sending keyframe to editor (%s)", err)
				return
			}

			frameID, err := util.ModelGetValue[int](keyTree.model.ToTreeModel(), iter, FRAME_KEY_ID)
			if err != nil {
				log.Printf("Error sending keyframe to editor (%s)", err)
				return
			}

			frameType, err := util.ModelGetValue[string](keyTree.model.ToTreeModel(), iter, FRAME_KEY_TYPE)
			if err != nil {
				log.Printf("Error sending keyframe to editor (%s)", err)
				return
			}

			sendEditor(frameID, frameType)
		})

	var column *gtk.TreeViewColumn

	geoCell, err := gtk.CellRendererTextNew()
	if err != nil {
		return
	}

	{
		// Geometry Name
		column, err = gtk.TreeViewColumnNewWithAttribute("Geometry", geoCell, "text", FRAME_GEOMETRY)
		if err != nil {
			return
		}

		column.SetResizable(true)
		keyTree.view.AppendColumn(column)

	}

	{
		// Attribute Name
		column, err = gtk.TreeViewColumnNewWithAttribute("Attribute", geoCell, "text", FRAME_ATTR_NAME)
		if err != nil {
			return
		}

		column.SetResizable(true)
		keyTree.view.AppendColumn(column)
	}

	{
		// Keyframe Type
		column, err = gtk.TreeViewColumnNewWithAttribute("Keyframe Type", geoCell, "text", FRAME_KEY_TYPE)
		if err != nil {
			return
		}

		column.SetResizable(true)
		keyTree.view.AppendColumn(column)
	}

	{
		// Bool Value
		var toggleCell *gtk.CellRendererToggle
		var column *gtk.TreeViewColumn

		toggleCell, err = gtk.CellRendererToggleNew()
		if err != nil {
			return
		}

		toggleCell.SetProperty("activatable", true)
		toggleCell.Connect("toggled",
			func(cell *gtk.CellRendererToggle, path string) {
				iter, err := keyTree.model.GetIterFromString(path)
				if err != nil {
					log.Printf("Error toggling toggle (%s)", err)
					return
				}

				state, err := util.ModelGetValue[bool](keyTree.model.ToTreeModel(), iter, FRAME_EXPAND)
				if err != nil {
					log.Printf("Error toggling toggle: %s", err)
					return
				}

				keyTree.model.SetValue(iter, FRAME_EXPAND, !state)
				frameID, err := util.ModelGetValue[int](keyTree.model.ToTreeModel(), iter, FRAME_KEY_ID)
				if err != nil {
					log.Printf("Error toggling toggle: %s", err)
					return
				}

				if frame, ok := keyTree.bindFrame[frameID]; ok {
					frame.Expand = !state
					keyTree.bindFrame[frameID] = frame
				} else if frame, ok := keyTree.setFrame[frameID]; ok {
					frame.Expand = !state
					keyTree.setFrame[frameID] = frame
				} else if frame, ok := keyTree.userFrame[frameID]; ok {
					frame.Expand = !state
					keyTree.userFrame[frameID] = frame
				}
			})

		column, err = gtk.TreeViewColumnNewWithAttribute("Expand", toggleCell, "active", FRAME_EXPAND)
		if err != nil {
			return
		}

		column.SetResizable(true)
		keyTree.view.AppendColumn(column)

	}

	return
}

func (keyTree *KeyTree) AddKeyframe(frameNum, geoID int, geoName, attrType, frameType string) {
	frame := templates.NewKeyFrame(frameNum, geoID, attrType, false)
	frame.Type = frameType

	frameID := keyTree.addGeometryRow(geoName, *frame)
	switch frameType {
	case templates.USER_FRAME:
		keyTree.userFrame[frameID] = *templates.NewUserFrame(*frame)
	case templates.SET_FRAME:
		keyTree.setFrame[frameID] = *templates.NewSetFrame(*frame, 0)
	case templates.BIND_FRAME:
		keyTree.bindFrame[frameID] = *templates.NewBindFrame(*frame, *templates.NewKeyFrame(0, 0, "", false))
	}
}

func (keyTree *KeyTree) addGeometryRow(geoName string, frame templates.Keyframe) (id int) {
	iter := keyTree.model.Append()

	id = keyTree.frameID
	keyTree.model.SetValue(iter, FRAME_GEOMETRY_ID, frame.GeoID)
	keyTree.model.SetValue(iter, FRAME_GEOMETRY, geoName)
	keyTree.model.SetValue(iter, FRAME_ATTR_TYPE, frame.GeoAttr)
	keyTree.model.SetValue(iter, FRAME_ATTR_NAME, geometry.Attrs[frame.GeoAttr])
	keyTree.model.SetValue(iter, FRAME_KEY_TYPE, frame.Type)
	keyTree.model.SetValue(iter, FRAME_KEY_ID, keyTree.frameID)

	keyTree.frameID++
	return
}

func (keyTree *KeyTree) RemoveKeyframe() (err error) {
	selection, err := keyTree.view.GetSelection()
	if err != nil {
		return
	}

	_, iter, ok := selection.GetSelected()
	if !ok {
		return
	}

	frameID, err := util.ModelGetValue[int](keyTree.model.ToTreeModel(), iter, FRAME_KEY_ID)
	if err != nil {
		return
	}

	delete(keyTree.bindFrame, frameID)
	delete(keyTree.setFrame, frameID)
	delete(keyTree.userFrame, frameID)

	keyTree.model.Remove(iter)
	return
}

func (keyTree *KeyTree) UpdateBindFrame(edit *templates.BindFrameEditor, frameID int) {
	frame, ok := keyTree.bindFrame[frameID]
	if !ok {
		return
	}

	newFrame, err := edit.UpdateKeyframe(frame)
	if err != nil {
		log.Printf("Error updating template keyframe: %s", err)
		return
	}

	keyTree.bindFrame[frameID] = newFrame
}

func (keyTree *KeyTree) UpdateSetFrame(edit *templates.SetFrameEditor, frameID int) {
	frame, ok := keyTree.setFrame[frameID]
	if !ok {
		return
	}

	newFrame := edit.UpdateKeyframe(frame)
	keyTree.setFrame[frameID] = newFrame
}

func (keyTree *KeyTree) ImportUserFrame(geo geometry.Geometry, frame templates.UserFrame) {
	frameID := keyTree.addGeometryRow(geo.Name, frame.Keyframe)
	keyTree.userFrame[frameID] = frame
	return
}

func (keyTree *KeyTree) ImportBindFrame(geo geometry.Geometry, frame templates.BindFrame) {
	frameID := keyTree.addGeometryRow(geo.Name, frame.Keyframe)
	keyTree.bindFrame[frameID] = frame
}

func (keyTree *KeyTree) ImportSetFrame(geo geometry.Geometry, frame templates.SetFrame) {
	frameID := keyTree.addGeometryRow(geo.Name, frame.Keyframe)
	keyTree.setFrame[frameID] = frame
}

func (keyTree *KeyTree) ExportKeyframes(temp *templates.Template) (err error) {
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

func (keyTree *KeyTree) updateKeys(geoID int, name string) {
	iter, ok := keyTree.model.GetIterFirst()

	for ; ok; ok = keyTree.model.IterNext(iter) {
		currentID, err := util.ModelGetValue[int](keyTree.model.ToTreeModel(), iter, FRAME_GEOMETRY_ID)
		if err != nil {
			log.Printf("Error getting keyframe geo id (%s)", err)
			continue
		}

		if currentID == geoID {
			keyTree.model.SetValue(iter, FRAME_GEOMETRY, name)
		}
	}
}
