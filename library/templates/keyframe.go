package templates

import (
	"chroma-viz/library/geometry"
	"chroma-viz/library/util"
	"log"
	"math"

	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
)

const (
	BIND_FRAME = "BindFrame"
	SET_FRAME  = "SetFrame"
	USER_FRAME = "UserFrame"
)

type Keyframe struct {
	FrameNum int
	GeoID    int
	GeoAttr  string
	Type     string
	Expand   bool
}

func NewKeyFrame(num, geo int, attr string, expand bool) *Keyframe {
	frame := &Keyframe{
		FrameNum: num,
		GeoID:    geo,
		GeoAttr:  attr,
		Expand:   expand,
	}

	return frame
}

func (key *Keyframe) Key() *Keyframe {
	return key
}

type BindFrame struct {
	Keyframe
	Bind Keyframe
}

func NewBindFrame(frame, bind Keyframe) *BindFrame {
	frame.Type = BIND_FRAME

	return &BindFrame{
		Keyframe: frame,
		Bind:     bind,
	}
}

type SetFrame struct {
	Keyframe
	Value int
}

func NewSetFrame(frame Keyframe, value int) *SetFrame {
	frame.Type = SET_FRAME

	return &SetFrame{
		Keyframe: frame,
		Value:    value,
	}
}

type UserFrame struct {
	Keyframe
}

func NewUserFrame(frame Keyframe) *UserFrame {
	frame.Type = USER_FRAME

	return &UserFrame{Keyframe: frame}
}

type KeyframeEditor struct {
	Scroll *gtk.ScrolledWindow
	Box    *gtk.Box
	ID     int
}

func NewKeyFrameEditor(id int) (edit *KeyframeEditor, err error) {
	edit = &KeyframeEditor{
		ID: id,
	}

	edit.Scroll, err = gtk.ScrolledWindowNew(nil, nil)
	if err != nil {
		return
	}

	edit.Scroll.SetVisible(true)
	edit.Scroll.SetVExpand(true)

	edit.Box, err = gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0)
	if err != nil {
		return
	}

	edit.Box.SetVisible(true)
	edit.Scroll.Add(edit.Box)

	return
}

type SetFrameEditor struct {
	KeyframeEditor

	Value *gtk.SpinButton
}

func NewSetFrameEditor(frame SetFrame, id int) (edit *SetFrameEditor, err error) {
	keyEdit, err := NewKeyFrameEditor(id)
	if err != nil {
		return
	}

	edit = &SetFrameEditor{
		KeyframeEditor: *keyEdit,
	}

	box, err := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 0)
	if err != nil {
		return
	}

	box.SetVisible(true)
	edit.Box.PackStart(box, false, false, 10)

	label, err := gtk.LabelNew("Value")
	if err != nil {
		return
	}

	label.SetVisible(true)
	label.SetWidthChars(12)
	box.PackStart(label, false, false, 10)

	edit.Value, err = gtk.SpinButtonNewWithRange(float64(math.MinInt), float64(math.MaxInt), 1)
	if err != nil {
		return
	}

	edit.Value.SetVisible(true)
	edit.Value.SetValue(float64(frame.Value))
	box.PackStart(edit.Value, false, false, 0)

	return
}

func (edit *SetFrameEditor) UpdateKeyframe(frame SetFrame) SetFrame {
	frame.Value = int(edit.Value.GetValue())
	return frame
}

type BindFrameEditor struct {
	KeyframeEditor

	Frame     *gtk.ComboBox
	Geometry  *gtk.ComboBox
	Attribute *gtk.ComboBox
}

func NewBindFrameEditor(frame BindFrame, id int, frameModel *gtk.ListStore, geoModel *gtk.ListStore) (edit *BindFrameEditor, err error) {
	keyEdit, err := NewKeyFrameEditor(id)
	if err != nil {
		return
	}

	edit = &BindFrameEditor{
		KeyframeEditor: *keyEdit,
	}

	geoCell, err := gtk.CellRendererTextNew()
	if err != nil {
		return
	}

	var box *gtk.Box
	var label *gtk.Label
	{
		edit.Frame, err = gtk.ComboBoxNew()
		if err != nil {
			return
		}

		edit.Frame.PackStart(geoCell, true)
		edit.Frame.CellLayout.AddAttribute(geoCell, "text", 1)
		edit.Frame.SetActive(1)
		edit.Frame.SetModel(frameModel)

		label, err = gtk.LabelNew("Frame")
		if err != nil {
			return
		}

		label.SetVisible(true)
		label.SetWidthChars(12)

		box, err = gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 10)
		if err != nil {
			return
		}

		box.PackStart(label, false, false, 10)
		box.PackStart(edit.Frame, false, false, 0)
	}

	{
		edit.Geometry, err = gtk.ComboBoxNew()
		if err != nil {
			return
		}

		edit.Geometry.PackStart(geoCell, true)
		edit.Geometry.CellLayout.AddAttribute(geoCell, "text", 1)
		edit.Geometry.SetActive(1)
		edit.Geometry.SetModel(geoModel)

		label, err = gtk.LabelNew("Geometry")
		if err != nil {
			return
		}

		label.SetVisible(true)
		label.SetWidthChars(12)

		box, err = gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 10)
		if err != nil {
			return
		}

		box.PackStart(label, false, false, 10)
		box.PackStart(edit.Geometry, false, false, 0)
	}

	{
		var attrModel *gtk.ListStore

		edit.Attribute, err = gtk.ComboBoxNew()
		if err != nil {
			return
		}

		attrModel, err = gtk.ListStoreNew(glib.TYPE_STRING, glib.TYPE_STRING)
		if err != nil {
			return
		}

		edit.Attribute.PackStart(geoCell, true)
		edit.Attribute.CellLayout.AddAttribute(geoCell, "text", 1)
		edit.Attribute.SetActive(1)
		edit.Attribute.SetModel(attrModel)

		edit.Geometry.Connect("changed", func() {
			iter, err := edit.Geometry.GetActiveIter()
			if err != nil {
				log.Printf("Error selecting geometry: %s", err)
				return
			}

			geoType, err := util.ModelGetValue[string](geoModel.ToTreeModel(), iter, 0)
			if err != nil {
				log.Printf("Error selecting geometry: %s", err)
				return
			}

			geometry.UpdateAttrList(attrModel, geoType)
		})

		label, err = gtk.LabelNew("Attribute")
		if err != nil {
			return
		}

		label.SetVisible(true)
		label.SetWidthChars(12)

		box, err = gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 10)
		if err != nil {
			return
		}

		box.PackStart(label, false, false, 10)
		box.PackStart(edit.Attribute, false, false, 0)
	}

	return
}

func (edit *BindFrameEditor) UpdateKeyframe(frame BindFrame) (BindFrame, error) {
	frameNum, err := comboBoxSelection[int](edit.Frame, 0)
	if err != nil {
		return frame, err
	}

	frame.Bind.FrameNum = frameNum

	geoID, err := comboBoxSelection[int](edit.Geometry, 0)
	if err != nil {
		return frame, err
	}

	frame.Bind.GeoID = geoID

	attr, err := comboBoxSelection[string](edit.Attribute, 0)
	if err != nil {
		return frame, err
	}

	frame.Bind.GeoAttr = attr
	return frame, nil
}

func comboBoxSelection[T any](combo *gtk.ComboBox, col int) (t T, err error) {
	iter, err := combo.GetActiveIter()
	if err != nil {
		return
	}

	model, err := combo.GetModel()
	if err != nil {
		return
	}

	return util.ModelGetValue[T](model.ToTreeModel(), iter, col)
}
