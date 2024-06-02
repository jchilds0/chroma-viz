package artist

import (
	"chroma-viz/library/gtk_utils"
	"chroma-viz/library/pages"
	"chroma-viz/library/props"
	"chroma-viz/library/templates"
	"log"
	"strconv"

	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
)

const (
	GEO_TYPE = iota
	GEO_NAME
	GEO_NUM
	GEO_NUM_COLS
)

const (
	FRAME_NUM = iota
	FRAME_GEOMETRY
	FRAME_GEOMETRY_ID
	FRAME_ATTR
	FRAME_VALUE
	FRAME_USER_VALUE
	FRAME_MASK
	FRAME_EXPAND
	FRAME_BIND_FRAME
	FRAME_BIND_GEO
	FRAME_BIND_ATTR
	FRAME_NUM_COLS
)

type TempTree struct {
	geoModel *gtk.TreeStore
	keyModel *gtk.TreeStore
	geoView  *gtk.TreeView
	geoList  *gtk.TreeStore
	keyView  *gtk.TreeView
}

func NewTempTree(propToEditor func(propID int)) (*TempTree, error) {
	temp := &TempTree{}

	err := temp.createGeometryTree(propToEditor)
	if err != nil {
		return nil, err
	}

	err = temp.createKeyTree()
	if err != nil {
		return nil, err
	}

	return temp, nil
}

func (temp *TempTree) createGeometryTree(propToEditor func(propID int)) (err error) {
	temp.geoView, err = gtk.TreeViewNew()
	if err != nil {
		return
	}

	temp.geoView.Set("reorderable", true)

	typeCell, err := gtk.CellRendererTextNew()
	if err != nil {
		return
	}

	column, err := gtk.TreeViewColumnNewWithAttribute("Geometry", typeCell, "text", GEO_TYPE)
	if err != nil {
		return
	}
	temp.geoView.AppendColumn(column)

	nameCell, err := gtk.CellRendererTextNew()
	if err != nil {
		return
	}

	nameCell.SetProperty("editable", true)
	nameCell.Connect("edited", func(cell *gtk.CellRendererText, path, text string) {
		iter, err := temp.geoModel.GetIterFromString(path)
		if err != nil {
			log.Printf("Error editing geometry (%s)", err)
			return
		}

		model := temp.geoModel.ToTreeModel()
		geoID, err := gtk_utils.ModelGetValue[int](model, iter, GEO_NUM)
		if err != nil {
			log.Printf("Error editing geometry (%s)", err)
			return
		}

		geo := page.PropMap[geoID]
		if geo == nil {
			log.Printf("Error getting geometry %d", geoID)
			return
		}

		geo.Name = text
		temp.geoModel.SetValue(iter, GEO_NAME, text)
		temp.updateKeys(geoID, text)
		temp.updateGeometry(geoID, text)
	})

	column, err = gtk.TreeViewColumnNewWithAttribute("Name", nameCell, "text", GEO_NAME)
	if err != nil {
		return
	}
	temp.geoView.AppendColumn(column)

	temp.geoModel, err = gtk.TreeStoreNew(
		glib.TYPE_STRING, // GEO TYPE
		glib.TYPE_STRING, // GEO NAME
		glib.TYPE_INT,    // GEO NUM
	)
	if err != nil {
		return
	}

	temp.geoList, err = gtk.TreeStoreNew(
		glib.TYPE_STRING,
		glib.TYPE_STRING,
		glib.TYPE_INT,
	)
	if err != nil {
		return
	}

	temp.geoView.SetModel(temp.geoModel)

	temp.geoView.Connect("row-activated",
		func(tree *gtk.TreeView, path *gtk.TreePath, column *gtk.TreeViewColumn) {
			iter, err := temp.geoModel.GetIter(path)
			if err != nil {
				log.Printf("Error sending prop to editor (%s)", err)
				return
			}

			model := &temp.geoModel.TreeModel
			propID, err := gtk_utils.ModelGetValue[int](model, iter, GEO_NUM)
			if err != nil {
				log.Printf("Error sending prop to editor (%s)", err)
				return
			}

			propToEditor(propID)
		})

	return nil
}

func (temp *TempTree) createKeyTree() (err error) {
	temp.keyView, err = gtk.TreeViewNew()
	if err != nil {
		return
	}

	temp.keyView.Set("reorderable", true)

	temp.keyModel, err = gtk.TreeStoreNew(
		glib.TYPE_INT,     // Frame Num
		glib.TYPE_STRING,  // Geometry Name
		glib.TYPE_INT,     // Geometry Num
		glib.TYPE_STRING,  // Geometry Attr
		glib.TYPE_INT,     // Value Entry
		glib.TYPE_BOOLEAN, // User Value Selector
		glib.TYPE_BOOLEAN, // Mask Parent
		glib.TYPE_BOOLEAN, // Expand Children
		glib.TYPE_STRING,  // Derived Value Frame
		glib.TYPE_STRING,  // Derived Value Geo
		glib.TYPE_STRING,  // Derived Value Attr
	)
	if err != nil {
		return
	}

	temp.keyView.SetModel(temp.keyModel)

	// Frame Number
	{

		var frameNumCell *gtk.CellRendererText
		var column *gtk.TreeViewColumn

		frameNumCell, err = gtk.CellRendererTextNew()
		if err != nil {
			return
		}

		frameNumCell.SetProperty("editable", true)
		frameNumCell.Connect("edited", func(cell *gtk.CellRendererText, path, text string) {
			iter, err := temp.keyModel.GetIterFromString(path)
			if err != nil {
				log.Printf("Error editing keyframe (%s)", err)
				return
			}

			num, err := strconv.Atoi(text)
			if err != nil {
				log.Printf("Error editing keyframe (%s)", err)
				return
			}

			temp.keyModel.SetValue(iter, FRAME_NUM, num)
		})

		column, err = gtk.TreeViewColumnNewWithAttribute("Frame Number", frameNumCell, "text", FRAME_NUM)
		if err != nil {
			return
		}
		temp.keyView.AppendColumn(column)

	}

	// Geometry Name
	{

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

		temp.keyView.AppendColumn(column)

	}

	// Attribute Name
	{

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

		temp.keyView.AppendColumn(column)

	}

	// Set Value
	{

		var valueCell *gtk.CellRendererText
		var column *gtk.TreeViewColumn

		valueCell, err = gtk.CellRendererTextNew()
		if err != nil {
			return
		}

		valueCell.SetProperty("editable", true)
		valueCell.Connect("edited", func(cell *gtk.CellRendererText, path, text string) {
			iter, err := temp.keyModel.GetIterFromString(path)
			if err != nil {
				log.Printf("Error editing keyframe (%s)", err)
				return
			}

			num, err := strconv.Atoi(text)
			if err != nil {
				log.Printf("Error editing keyframe (%s)", err)
				return
			}

			temp.keyModel.SetValue(iter, FRAME_VALUE, num)
		})

		column, err = gtk.TreeViewColumnNewWithAttribute("Set Value", valueCell, "text", FRAME_VALUE)
		if err != nil {
			return
		}

		temp.keyView.AppendColumn(column)

	}

	// Bool Value
	{

		var toggleCell *gtk.CellRendererToggle
		var column *gtk.TreeViewColumn

		names := []string{"Mask", "Expand", "User Value"}
		cols := []int{FRAME_MASK, FRAME_EXPAND, FRAME_USER_VALUE}

		for i := range names {
			toggleCell, err = gtk.CellRendererToggleNew()
			if err != nil {
				return
			}

			toggleCell.SetProperty("activatable", true)
			toggleCell.Connect("toggled",
				func(cell *gtk.CellRendererToggle, path string) {
					iter, err := temp.keyModel.GetIterFromString(path)
					if err != nil {
						log.Printf("Error toggling toggle (%s)", err)
						return
					}

					model := temp.keyModel.ToTreeModel()

					state, err := gtk_utils.ModelGetValue[bool](model, iter, cols[i])
					if err != nil {
						log.Printf("Error toggling toggle (%s)", err)
						return
					}

					temp.keyModel.SetValue(iter, cols[i], !state)
				})

			column, err = gtk.TreeViewColumnNewWithAttribute(names[i], toggleCell, "active", cols[i])
			if err != nil {
				return
			}

			temp.keyView.AppendColumn(column)
		}

	}

	// Derived Value
	{

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
				iter, err := temp.keyModel.GetIterFromString(path)
				if err != nil {
					log.Printf("Error editing geometry (%s)", err)
					return
				}

				temp.keyModel.SetValue(iter, cols[i], text)
			})
		}

		temp.keyView.AppendColumn(column)

	}

	return nil
}

func (tempView *TempTree) updateGeometry(geoID int, name string) {
	iter, ok := tempView.geoList.GetIterFirst()
	model := tempView.geoList.ToTreeModel()

	for ok {
		currentID, err := gtk_utils.ModelGetValue[int](model, iter, GEO_NUM)
		if err != nil {
			log.Printf("Error getting geometry (%s)", err)
			ok = model.IterNext(iter)
			continue
		}

		if currentID == geoID {
			tempView.geoList.SetValue(iter, GEO_NAME, name)
		}

		ok = model.IterNext(iter)
	}
}

func (tempView *TempTree) removeGeometry(geoID int) {
	iter, ok := tempView.geoList.GetIterFirst()
	model := tempView.geoList.ToTreeModel()

	for ok {
		currentID, err := gtk_utils.ModelGetValue[int](model, iter, GEO_NUM)
		if err != nil {
			log.Printf("Error getting geometry (%s)", err)
			ok = tempView.geoList.IterNext(iter)
			continue
		}

		if currentID == geoID {
			tempView.geoList.Remove(iter)
			iter, ok = tempView.geoList.GetIterFirst()
		} else {
			ok = tempView.geoList.IterNext(iter)
		}
	}
}

func (tempView *TempTree) updateKeys(geoID int, name string) {
	iter, ok := tempView.keyModel.GetIterFirst()
	model := tempView.keyModel.ToTreeModel()

	for ok {
		currentID, err := gtk_utils.ModelGetValue[int](model, iter, FRAME_GEOMETRY_ID)
		if err != nil {
			log.Printf("Error getting keyframe geo id (%s)", err)
			ok = model.IterNext(iter)
			continue
		}

		if currentID == geoID {
			tempView.keyModel.SetValue(iter, FRAME_GEOMETRY, name)
		}

		ok = model.IterNext(iter)
	}
}

func (tempView *TempTree) removeKeys(geoID int) {
	iter, ok := tempView.keyModel.GetIterFirst()
	model := tempView.keyModel.ToTreeModel()

	for ok {
		currentID, err := gtk_utils.ModelGetValue[int](model, iter, FRAME_GEOMETRY_ID)
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

func (tempView *TempTree) keyframes(temp *templates.Template) {
	iter, ok := tempView.keyModel.GetIterFirst()
	keyModel := tempView.keyModel.ToTreeModel()

	var bindFrame, bindGeo, bindAttr string
	var user bool

	for ok {
		frame, err := tempView.getKeyframe(iter)
		if err != nil {
			goto ERROR
		}

		user, err = gtk_utils.ModelGetValue[bool](keyModel, iter, FRAME_USER_VALUE)
		if err != nil {
			goto ERROR
		}

		if user {
			keyframe := templates.NewUserFrame(frame)
			temp.UserFrame = append(temp.UserFrame, *keyframe)

			ok = tempView.keyModel.IterNext(iter)
			continue
		}

		bindFrame, err = gtk_utils.ModelGetValue[string](keyModel, iter, FRAME_BIND_FRAME)
		if err != nil {
			goto ERROR
		}

		bindGeo, err = gtk_utils.ModelGetValue[string](keyModel, iter, FRAME_BIND_GEO)
		if err != nil {
			goto ERROR
		}

		bindAttr, err = gtk_utils.ModelGetValue[string](keyModel, iter, FRAME_BIND_ATTR)
		if err != nil {
			goto ERROR
		}

		if bindFrame != "" && bindGeo != "" && bindAttr != "" {
			frameNum, _ := strconv.Atoi(bindFrame)
			geoNum, _ := strconv.Atoi(bindGeo)

			bind := templates.NewKeyFrame(frameNum, geoNum, bindAttr, templates.BIND_FRAME, false, false)

			keyframe := templates.NewBindFrame(frame, *bind)
			temp.BindFrame = append(temp.BindFrame, *keyframe)
			ok = tempView.keyModel.IterNext(iter)
			continue
		}

		{
			var value int
			value, err = gtk_utils.ModelGetValue[int](keyModel, iter, FRAME_VALUE)
			keyframe := templates.NewSetFrame(frame, value)

			temp.SetFrame = append(temp.SetFrame, *keyframe)
			ok = tempView.keyModel.IterNext(iter)
			continue
		}

	ERROR:
		log.Printf("Error getting keyframe (%s)", err)
		ok = tempView.keyModel.IterNext(iter)
		continue
	}
}

func (tempView *TempTree) getKeyframe(iter *gtk.TreeIter) (frame templates.Keyframe, err error) {
	keyModel := tempView.keyModel.ToTreeModel()

	frame.FrameNum, err = gtk_utils.ModelGetValue[int](keyModel, iter, FRAME_NUM)
	if err != nil {
		return
	}

	frame.GeoID, err = gtk_utils.ModelGetValue[int](keyModel, iter, FRAME_GEOMETRY_ID)
	if err != nil {
		return
	}

	frame.GeoAttr, err = gtk_utils.ModelGetValue[string](keyModel, iter, FRAME_ATTR)
	if err != nil {
		return
	}

	frame.Mask, err = gtk_utils.ModelGetValue[bool](keyModel, iter, FRAME_MASK)
	if err != nil {
		return
	}

	frame.Expand, err = gtk_utils.ModelGetValue[bool](keyModel, iter, FRAME_EXPAND)
	if err != nil {
		return
	}

	return
}

func (tempView *TempTree) addKeyframes(page *pages.Page, temp *templates.Template) {
	for _, frame := range temp.UserFrame {
		geo := page.PropMap[frame.GeoID]
		if geo == nil {
			log.Printf("Missing geometry %d for userframe", frame.FrameNum)
			continue
		}

		iter := tempView.keyModel.Append(nil)
		tempView.updateBaseKeyframe(iter, frame.Key(), geo)
		tempView.keyModel.SetValue(iter, FRAME_USER_VALUE, true)
	}

	for _, frame := range temp.BindFrame {
		geo := page.PropMap[frame.GeoID]
		if geo == nil {
			log.Printf("Missing geometry %d for bindframe", frame.FrameNum)
			continue
		}

		iter := tempView.keyModel.Append(nil)
		tempView.updateBaseKeyframe(iter, frame.Key(), geo)

		tempView.keyModel.SetValue(iter, FRAME_BIND_FRAME, frame.Bind.FrameNum)
		tempView.keyModel.SetValue(iter, FRAME_BIND_GEO, frame.Bind.GeoID)
		tempView.keyModel.SetValue(iter, FRAME_BIND_ATTR, frame.Bind.GeoAttr)
	}

	for _, frame := range temp.SetFrame {
		geo := page.PropMap[frame.GeoID]
		if geo == nil {
			log.Printf("Missing geometry %d for setframe", frame.FrameNum)
			continue
		}

		iter := tempView.keyModel.Append(nil)
		tempView.updateBaseKeyframe(iter, frame.Key(), geo)

		tempView.keyModel.SetValue(iter, FRAME_VALUE, frame.Value)
	}
}

func (tempView *TempTree) updateBaseKeyframe(iter *gtk.TreeIter, frame *templates.Keyframe, geo *props.Property) {
	tempView.keyModel.SetValue(iter, FRAME_NUM, frame.FrameNum)
	tempView.keyModel.SetValue(iter, FRAME_GEOMETRY, geo.Name)
	tempView.keyModel.SetValue(iter, FRAME_GEOMETRY_ID, frame.GeoID)
	tempView.keyModel.SetValue(iter, FRAME_ATTR, frame.GeoAttr)
	tempView.keyModel.SetValue(iter, FRAME_MASK, frame.Mask)
	tempView.keyModel.SetValue(iter, FRAME_EXPAND, frame.Expand)

}

func (tempView *TempTree) AddGeoRow(iter *gtk.TreeIter, name, propName string, propNum int) {
	tempView.geoModel.SetValue(iter, GEO_TYPE, propName)
	tempView.geoModel.SetValue(iter, GEO_NAME, name)
	tempView.geoModel.SetValue(iter, GEO_NUM, propNum)

	newIter := tempView.geoList.Append(nil)
	tempView.geoList.SetValue(newIter, GEO_TYPE, propName)
	tempView.geoList.SetValue(newIter, GEO_NAME, name)
	tempView.geoList.SetValue(newIter, GEO_NUM, propNum)
}

func (tempView *TempTree) AddKeyRow(iter *gtk.TreeIter, geoName string, geoID int, attrName string) {
	tempView.keyModel.SetValue(iter, FRAME_NUM, 0)
	tempView.keyModel.SetValue(iter, FRAME_GEOMETRY, geoName)
	tempView.keyModel.SetValue(iter, FRAME_GEOMETRY_ID, geoID)
	tempView.keyModel.SetValue(iter, FRAME_ATTR, attrName)
}

func (tempView *TempTree) Clean() {
	var err error
	tempView.geoModel, err = gtk.TreeStoreNew(glib.TYPE_STRING, glib.TYPE_STRING, glib.TYPE_INT)
	if err != nil {
		log.Print(err)
		return
	}

	tempView.geoView.SetModel(tempView.geoModel)

	tempView.keyModel, err = gtk.TreeStoreNew(
		glib.TYPE_INT,
		glib.TYPE_STRING,
		glib.TYPE_INT,
		glib.TYPE_STRING,
		glib.TYPE_INT,
		glib.TYPE_INT,
		glib.TYPE_BOOLEAN,
		glib.TYPE_BOOLEAN,
		glib.TYPE_BOOLEAN,
		glib.TYPE_STRING,
		glib.TYPE_STRING,
		glib.TYPE_STRING,
	)
	if err != nil {
		return
	}

	tempView.keyView.SetModel(tempView.keyModel)
}
