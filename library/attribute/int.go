package attribute

import (
	"fmt"
	"strconv"

	"github.com/gotk3/gotk3/gtk"
)

type IntAttribute struct {
	Name  string
	Type  int
	Value int
}

func NewIntAttribute(name string) *IntAttribute {
	intAttr := &IntAttribute{
		Name: name,
		Type: INT,
	}

	return intAttr
}

func NewIntAttributeWithValue(name string, value int) *IntAttribute {
	intAttr := &IntAttribute{
		Name:  name,
		Value: value,
		Type:  INT,
	}

	return intAttr
}

func (intAttr *IntAttribute) String() string {
	return fmt.Sprintf("%s=%d#", intAttr.Name, intAttr.Value)
}

func (intAttr *IntAttribute) Encode() string {
	return fmt.Sprintf("%d", intAttr.Value)
}

func (intAttr *IntAttribute) Decode(value string) (err error) {
	intAttr.Value, err = strconv.Atoi(value)
	if err != nil {
		err = fmt.Errorf("Error decoding int attr (%s)", err)
		return
	}

	return
}

func (intAttr *IntAttribute) Copy(attr Attribute) (err error) {
	intAttrCopy, ok := attr.(*IntAttribute)
	if !ok {
		err = fmt.Errorf("Attribute not IntAttribute")
		return
	}

	intAttr.Value = intAttrCopy.Value
	return
}

func (intAttr *IntAttribute) Update(edit Editor) error {
	intEdit, ok := edit.(*IntEditor)
	if !ok {
		return fmt.Errorf("IntAttribute.Update requires IntEditor")
	}

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

func (intEdit *IntEditor) Update(attr Attribute) error {
	intAttr, ok := attr.(*IntAttribute)
	if !ok {
		return fmt.Errorf("IntEditor.Update requires IntAttribute")
	}

	intEdit.button.SetValue(float64(intAttr.Value))
	return nil
}

func (intEdit *IntEditor) Box() *gtk.Box {
	return intEdit.box
}

func (intEdit *IntEditor) Expand() bool {
	return false
}
