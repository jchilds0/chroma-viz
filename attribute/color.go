package attribute

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/gotk3/gotk3/gtk"
)

type ColorAttribute struct {
    r     int
    g     int
    b     int
    a     int
}

func NewColorAttribute() *ColorAttribute {
    colorAttr := &ColorAttribute{}
    return colorAttr
}

func (colorAttr *ColorAttribute) String() string {
    return fmt.Sprintf("color=%d %d %d %d#", colorAttr.r, colorAttr.g, colorAttr.b, colorAttr.a)
}

func (colorAttr *ColorAttribute) Encode() string {
    return fmt.Sprintf("color %d %d %d %d;", colorAttr.r, colorAttr.g, colorAttr.b, colorAttr.a)
}

func (colorAttr *ColorAttribute) Decode(s string) (err error) {
    line := strings.Split(s, " ")
    if len(line) != 5 {
        return fmt.Errorf("Incorrect color attr string (%s)", line)
    }

    colorAttr.r, err = strconv.Atoi(line[1])
    if err != nil {
        return err
    }

    colorAttr.g, err = strconv.Atoi(line[2])
    if err != nil {
        return err
    }

    colorAttr.b, err = strconv.Atoi(line[3])
    if err != nil {
        return err
    }

    colorAttr.a, err = strconv.Atoi(line[4])
    if err != nil {
        return err
    }

    return 
}

func (colorAttr *ColorAttribute) Update(edit Editor) error {
    colorEdit, ok := edit.(*ColorEditor)
    if !ok {
        return fmt.Errorf("ColorAttribute.Update requires ColorEditor") 
    }

    return nil
}

type ColorEditor struct {
    box *gtk.Box
}

func NewColorEditor(name string, animate func()) (*ColorEditor, error) {
    var err error
    colorEdit := &ColorEditor{}

    colorEdit.box, err = gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 0)
    if err != nil {
        return nil, err
    }

    colorEdit.box.SetVisible(true)
    label, err := gtk.LabelNew(name)
    if err != nil { 
        return nil, err
    }

    label.SetVisible(true)
    label.SetWidthChars(12)

    return colorEdit, nil
}

func (colorEdit *ColorEditor) Update(attr Attribute) error {
    colorAttr, ok := attr.(*ColorAttribute)
    if !ok {
        return fmt.Errorf("ColorEditor.Update requires ColorAttribute") 
    }

    return nil
}

func (colorEdit *ColorEditor) Box() *gtk.Box {
    return colorEdit.box
}
