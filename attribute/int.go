package attribute

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/gotk3/gotk3/gtk"
)

type IntAttribute struct {
    fileName    string
	chromaName  string
	Value int
}

func NewIntAttribute(file, chroma string) *IntAttribute {
    intAttr := &IntAttribute{
        fileName: file, 
        chromaName: chroma,
    }

    return intAttr
}

func (intAttr *IntAttribute) String() string {
    return fmt.Sprintf("%s=%d#", intAttr.chromaName, intAttr.Value)
}

func (intAttr *IntAttribute) Encode() string {
    return fmt.Sprintf("%s %d;", intAttr.fileName, intAttr.Value)
}

func (intAttr *IntAttribute) Decode(s string) (err error) {
    line := strings.Split(s, " ")
    if len(line) != 2 {
        return fmt.Errorf("Incorrect int attr string (%s)", line)
    }

    intAttr.Value, err = strconv.Atoi(line[1])
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
	box       *gtk.Box
    button    *gtk.SpinButton
}

func NewIntEditor(name string, lower, upper float64, animate func()) (*IntEditor, error) {
    var err error
    intEdit := &IntEditor{}

    intEdit.box, err = gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 0)
    if err != nil {
        return nil, err
    }

    intEdit.box.SetVisible(true)
    label, err := gtk.LabelNew(name)
    if err != nil { 
        return nil, err
    }

    label.SetVisible(true)
    label.SetWidthChars(12)
    intEdit.box.PackStart(label, false, false, padding)

    intEdit.button, err = gtk.SpinButtonNewWithRange(lower, upper, 1)
    if err != nil { 
        return nil, err
    }

    intEdit.button.SetVisible(true)
    intEdit.button.SetValue(0)
    intEdit.box.PackStart(intEdit.button, false, false, 0)
    intEdit.button.Connect("value-changed", animate)

    return intEdit, nil
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
