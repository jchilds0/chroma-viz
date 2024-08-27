package attribute

import (
	"fmt"

	"github.com/gotk3/gotk3/gtk"
)

type FloatAttribute struct {
	Name  string
	Value float64
}

func (floatAttr *FloatAttribute) Encode() string {
	return fmt.Sprintf("%s=%f#", floatAttr.Name, floatAttr.Value)
}

func (floatAttr *FloatAttribute) UpdateAttribute(floatEdit *FloatEditor) (err error) {
	floatAttr.Value = floatEdit.button.GetValue()
	return
}

type FloatEditor struct {
	Name   string
	Box    *gtk.Box
	button *gtk.SpinButton
}

func NewFloatEditor(name string, lower, upper, scale float64) (floatEdit *FloatEditor, err error) {
	floatEdit = &FloatEditor{Name: name}

	floatEdit.Box, err = gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 0)
	if err != nil {
		return
	}

	floatEdit.Box.SetVisible(true)
	label, err := gtk.LabelNew(name)
	if err != nil {
		return
	}

	label.SetVisible(true)
	label.SetWidthChars(12)
	floatEdit.Box.PackStart(label, false, false, uint(padding))

	floatEdit.button, err = gtk.SpinButtonNewWithRange(lower, upper, scale)
	if err != nil {
		return
	}

	floatEdit.button.SetVisible(true)
	floatEdit.button.SetValue(0)
	floatEdit.Box.PackStart(floatEdit.button, false, false, 0)

	return
}

func (floatEdit *FloatEditor) UpdateEditor(floatAttr *FloatAttribute) (err error) {
	floatEdit.button.SetValue(floatAttr.Value)
	return
}
