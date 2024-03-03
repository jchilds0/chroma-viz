package attribute

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/gotk3/gotk3/gtk"
)

type FloatAttribute struct {
	name  string
	value float64
}

func NewFloatAttribute(name string) *FloatAttribute {
    floatAttr := &FloatAttribute{name: name}
    return floatAttr
}

func (floatAttr *FloatAttribute) String() string {
    return fmt.Sprintf("%s=%d#", floatAttr.name, floatAttr.value)
}

func (floatAttr *FloatAttribute) Encode() string {
    return fmt.Sprintf("%s %d;", floatAttr.name, floatAttr.value)
}

func (floatAttr *FloatAttribute) Decode(s string) (err error) {
    line := strings.Split(s, " ")
    floatAttr.name = line[0]
    floatAttr.value, err = strconv.ParseFloat(line[1], 64)

    return 
}

func (floatAttr *FloatAttribute) Update(edit Editor) error {
    floatEdit, ok := edit.(*FloatEditor)
    if !ok {
        return fmt.Errorf("FloatAttribute.Update requires FloatEditor") 
    }

    floatAttr.value = floatEdit.button.GetValue()
    return nil
}

type FloatEditor struct {
	box       *gtk.Box
    button    *gtk.SpinButton
}

func NewFloatEditor(name string, lower, upper float64, animate func()) (*FloatEditor, error) {
    var err error
    floatEdit := &FloatEditor{}

    floatEdit.box, err = gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 0)
    if err != nil {
        return nil, err
    }

    floatEdit.box.SetVisible(true)
    label, err := gtk.LabelNew(name)
    if err != nil { 
        return nil, err
    }

    label.SetVisible(true)
    label.SetWidthChars(12)
    floatEdit.box.PackStart(label, false, false, uint(padding))

    floatEdit.button, err = gtk.SpinButtonNewWithRange(lower, upper, 1)
    if err != nil { 
        return nil, err
    }

    floatEdit.button.SetVisible(true)
    floatEdit.button.SetValue(0)
    floatEdit.box.PackStart(floatEdit.button, false, false, 0)
    floatEdit.button.Connect("value-changed", animate)

    return floatEdit, nil
}

func (floatEdit *FloatEditor) Update(attr Attribute) error {
    floatAttr, ok := attr.(*FloatAttribute)
    if !ok {
        return fmt.Errorf("FloatEditor.Update requires FloatAttribute") 
    }

    floatEdit.button.SetValue(floatAttr.value)
    return nil
}

func (floatEdit *FloatEditor) Box() *gtk.Box {
    return floatEdit.box
}
