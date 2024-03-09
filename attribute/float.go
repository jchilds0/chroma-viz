package attribute

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/gotk3/gotk3/gtk"
)

type FloatAttribute struct {
    fileName    string
	chromaName  string
	value float64
}

func NewFloatAttribute(file, chroma string) *FloatAttribute {
    floatAttr := &FloatAttribute{
        chromaName: chroma,
        fileName: file,
    }

    return floatAttr
}

func (floatAttr *FloatAttribute) String() string {
    return fmt.Sprintf("%s=%f#", floatAttr.chromaName, floatAttr.value)
}

func (floatAttr *FloatAttribute) Encode() string {
    return fmt.Sprintf("%s %f;", floatAttr.fileName, floatAttr.value)
}

func (floatAttr *FloatAttribute) Decode(s string) (err error) {
    line := strings.Split(s, " ")
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

func NewFloatEditor(name string, lower, upper, scale float64, animate func()) *FloatEditor {
    var err error
    floatEdit := &FloatEditor{}

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
    floatEdit.button.Connect("value-changed", animate)

    return floatEdit
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
