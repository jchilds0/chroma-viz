package attribute

import (
	"fmt"

	"github.com/gotk3/gotk3/gtk"
)

type IntAttribute struct {
	Name  string
	Value int
}

func (intAttr *IntAttribute) Encode() string {
	return fmt.Sprintf("%s=%d#", intAttr.Name, intAttr.Value)
}

func (intAttr *IntAttribute) UpdateAttribute(intEdit *IntEditor) error {
	intAttr.Value = intEdit.button.GetValueAsInt()
	return nil
}

type IntEditor struct {
	box    *gtk.Box
	button *gtk.SpinButton
	name   string
}

func NewIntEditor(name string, lower, upper float64) (intEdit *IntEditor, err error) {
	intEdit = &IntEditor{name: name}

	intEdit.box, err = gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 0)
	if err != nil {
		return
	}

	intEdit.box.SetVisible(true)
	label, err := gtk.LabelNew(name)
	if err != nil {
		return
	}

	label.SetVisible(true)
	label.SetWidthChars(12)
	intEdit.box.PackStart(label, false, false, padding)

	intEdit.button, err = gtk.SpinButtonNewWithRange(lower, upper, 1)
	if err != nil {
		return
	}

	intEdit.button.SetVisible(true)
	intEdit.button.SetValue(0)
	intEdit.box.PackStart(intEdit.button, false, false, 0)

	return
}

func (intEdit *IntEditor) Name() string {
	return intEdit.name
}

func (intEdit *IntEditor) UpdateEditor(intAttr *IntAttribute) error {
	intEdit.button.SetValue(float64(intAttr.Value))
	return nil
}
