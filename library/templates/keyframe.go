package templates

import (
	"chroma-viz/library/geometry"
	"chroma-viz/library/util"
	"fmt"
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
	Value float64
}

func NewSetFrame(frame Keyframe, value float64) *SetFrame {
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
}

func NewKeyFrameEditor() (edit *KeyframeEditor, err error) {
	edit = &KeyframeEditor{}

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

func NewSetFrameEditor() (edit *SetFrameEditor, err error) {
	keyEdit, err := NewKeyFrameEditor()
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
	edit.Value.SetValue(0)
	box.PackStart(edit.Value, false, false, 0)

	return
}

func (edit *SetFrameEditor) UpdateKeyframe(frame SetFrame) SetFrame {
	frame.Value = edit.Value.GetValue()
	return frame
}

func (edit *SetFrameEditor) UpdateEditor(frame SetFrame) {
	edit.Value.SetValue(float64(frame.Value))
}

type BindFrameEditor struct {
	KeyframeEditor

	Frame     *gtk.ComboBox
	Geometry  *gtk.ComboBox
	Attribute *gtk.ComboBox
	attrModel *gtk.ListStore
}

func NewBindFrameEditor(frameModel *gtk.ListStore, geoModel *gtk.ListStore) (edit *BindFrameEditor, err error) {
	keyEdit, err := NewKeyFrameEditor()
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

		edit.Frame.SetVisible(true)
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

		box.SetVisible(true)
		box.PackStart(label, false, false, 10)
		box.PackStart(edit.Frame, false, false, 0)

		edit.Box.PackStart(box, false, false, 10)
	}

	{
		edit.Geometry, err = gtk.ComboBoxNew()
		if err != nil {
			return
		}

		edit.Geometry.SetVisible(true)
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

		box.SetVisible(true)
		box.PackStart(label, false, false, 10)
		box.PackStart(edit.Geometry, false, false, 0)

		edit.Box.PackStart(box, false, false, 10)
	}

	{
		edit.Attribute, err = gtk.ComboBoxNew()
		if err != nil {
			return
		}

		edit.attrModel, err = gtk.ListStoreNew(glib.TYPE_STRING, glib.TYPE_STRING)
		if err != nil {
			return
		}

		edit.Attribute.SetVisible(true)
		edit.Attribute.PackStart(geoCell, true)
		edit.Attribute.CellLayout.AddAttribute(geoCell, "text", 1)
		edit.Attribute.SetActive(1)
		edit.Attribute.SetModel(edit.attrModel)

		edit.Geometry.Connect("changed", func() {
			iter, err := edit.Geometry.GetActiveIter()
			if err != nil {
				return
			}

			geoType, err := util.ModelGetValue[string](geoModel.ToTreeModel(), iter, 0)
			if err != nil {
				log.Printf("Error selecting geometry: %s", err)
				return
			}

			geometry.UpdateAttrList(edit.attrModel, geoType)
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

		box.SetVisible(true)
		box.PackStart(label, false, false, 10)
		box.PackStart(edit.Attribute, false, false, 0)

		edit.Box.PackStart(box, false, false, 10)
	}

	return
}

func (edit *BindFrameEditor) UpdateKeyframe(frame BindFrame) (BindFrame, error) {
	var err error
	frame.Bind.FrameNum, err = comboBoxSelection[int](edit.Frame, 0)
	if err != nil {
		return frame, err
	}

	frame.Bind.GeoID, err = comboBoxSelection[int](edit.Geometry, 2)
	if err != nil {
		return frame, err
	}

	frame.Bind.GeoAttr, err = comboBoxSelection[string](edit.Attribute, 0)
	return frame, err
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

func (edit *BindFrameEditor) UpdateEditor(frame BindFrame) (err error) {
	_, err = setComboSelection(edit.Frame, frame.Bind.FrameNum, 0)
	if err != nil {
		return
	}

	iter, err := setComboSelection(edit.Geometry, frame.Bind.GeoID, 2)
	if err != nil {
		return
	}

	geoModel, err := edit.Geometry.GetModel()
	if err != nil {
		return
	}

	geoType, err := util.ModelGetValue[string](geoModel.ToTreeModel(), iter, 0)
	if err != nil {
		return
	}

	geometry.UpdateAttrList(edit.attrModel, geoType)
	_, err = setComboSelection(edit.Attribute, frame.Bind.GeoAttr, 0)
	return
}

func setComboSelection[T comparable](combo *gtk.ComboBox, value T, col int) (iter *gtk.TreeIter, err error) {
	imodel, err := combo.GetModel()
	if err != nil {
		return
	}

	model := imodel.ToTreeModel()
	iter, ok := model.GetIterFirst()

	var row T
	for {
		if !ok {
			err = fmt.Errorf("Selection %v not found", value)
			return
		}

		row, err = util.ModelGetValue[T](model, iter, col)
		if err != nil {
			return
		}

		if row == value {
			break
		}

		ok = model.IterNext(iter)
	}

	combo.SetActiveIter(iter)
	return
}
