package attribute

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/gotk3/gotk3/gtk"
)

type FloatAttribute struct {
    FileName    string
	ChromaName  string
    Type        int
	Value       float64
}

func NewFloatAttribute(file, chroma string) *FloatAttribute {
    floatAttr := &FloatAttribute{
        ChromaName: chroma,
        FileName: file,
        Type: FLOAT,
    }

    return floatAttr
}

func (floatAttr *FloatAttribute) String() string {
    return fmt.Sprintf("%s=%f#", floatAttr.ChromaName, floatAttr.Value)
}

func (floatAttr *FloatAttribute) Encode() string {
    return fmt.Sprintf("%s %f;", floatAttr.FileName, floatAttr.Value)
}

func (floatAttr *FloatAttribute) Decode(s string) (err error) {
    line := strings.Split(s, " ")
    floatAttr.Value, err = strconv.ParseFloat(line[1], 64)

    return 
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
	box       *gtk.Box
    button    *gtk.SpinButton
}

func NewFloatEditor(name string, lower, upper, scale float64) *FloatEditor {
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

    return floatEdit
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
