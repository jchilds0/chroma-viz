package attribute

import (
	"fmt"
	"log"
	"strconv"

	"github.com/gotk3/gotk3/gtk"
)

type FloatAttribute struct {
	Name  string
	Type  int
	Value float64
}

func NewFloatAttribute(name string) *FloatAttribute {
	floatAttr := &FloatAttribute{
		Name: name,
		Type: FLOAT,
	}

	return floatAttr
}

func (floatAttr *FloatAttribute) String() string {
	return fmt.Sprintf("%s=%f#", floatAttr.Name, floatAttr.Value)
}

func (floatAttr *FloatAttribute) Encode() string {
	return fmt.Sprintf("{'name': '%s', 'value': '%f'}",
		floatAttr.Name, floatAttr.Value)
}

func (floatAttr *FloatAttribute) Decode(value string) {
	var err error

	floatAttr.Value, err = strconv.ParseFloat(value, 64)
	if err != nil {
		log.Printf("Error decoding float attr (%s)", err)
	}
}

func (floatAttr *FloatAttribute) Copy(attr Attribute) {
	floatAttrCopy, ok := attr.(*FloatAttribute)
	if !ok {
		log.Print("Attribute not FloatAttribute")
		return
	}

	floatAttr.Value = floatAttrCopy.Value
}

func (floatAttr *FloatAttribute) Update(edit Editor) error {
	floatEdit, ok := edit.(*FloatEditor)
	if !ok {
		return fmt.Errorf("FloatAttribute.Update requires FloatEditor")
	}

	floatAttr.Value = floatEdit.button.GetValue()
	return nil
}

type FloatEditor struct {
	box    *gtk.Box
	button *gtk.SpinButton
	name   string
}

func NewFloatEditor(name string, lower, upper, scale float64) *FloatEditor {
	var err error
	floatEdit := &FloatEditor{name: name}

	floatEdit.box, err = gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 0)
	if err != nil {
		log.Print(err)
		return nil
	}

	floatEdit.box.SetVisible(true)
	label, err := gtk.LabelNew(name)
	if err != nil {
		log.Print(err)
		return nil
	}

	label.SetVisible(true)
	label.SetWidthChars(12)
	floatEdit.box.PackStart(label, false, false, uint(padding))

	floatEdit.button, err = gtk.SpinButtonNewWithRange(lower, upper, scale)
	if err != nil {
		log.Print(err)
		return nil
	}

	floatEdit.button.SetVisible(true)
	floatEdit.button.SetValue(0)
	floatEdit.box.PackStart(floatEdit.button, false, false, 0)

	return floatEdit
}

func (floatEdit *FloatEditor) Name() string {
	return floatEdit.name
}

func (floatEdit *FloatEditor) Update(attr Attribute) error {
	floatAttr, ok := attr.(*FloatAttribute)
	if !ok {
		return fmt.Errorf("FloatEditor.Update requires FloatAttribute")
	}

	floatEdit.button.SetValue(floatAttr.Value)
	return nil
}

func (floatEdit *FloatEditor) Box() *gtk.Box {
	return floatEdit.box
}

func (floatEdit *FloatEditor) Expand() bool {
	return false
}
