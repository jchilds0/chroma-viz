package attribute

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"
)

type ColorAttribute struct {
    fileName    string 
    chromaName  string
    r     float64
    g     float64
    b     float64
    a     float64
}

func NewColorAttribute(file, chroma string) *ColorAttribute {
    colorAttr := &ColorAttribute{
        fileName: file, chromaName: chroma,
        r: 1.0, g: 1.0, b: 1.0, a: 1.0,
    }
    return colorAttr
}

func (colorAttr *ColorAttribute) String() string {
    return fmt.Sprintf("%s=%f %f %f %f#", colorAttr.chromaName, 
        colorAttr.r, colorAttr.g, colorAttr.b, colorAttr.a)
}

func (colorAttr *ColorAttribute) Encode() string {
    return fmt.Sprintf("%s %f %f %f %f;", colorAttr.fileName, 
        colorAttr.r, colorAttr.g, colorAttr.b, colorAttr.a)
}

func (colorAttr *ColorAttribute) Decode(s string) (err error) {
    line := strings.Split(s, " ")
    if len(line) != 5 {
        return fmt.Errorf("Incorrect color attr string (%s)", line)
    }

    colorAttr.r, err = strconv.ParseFloat(line[0], 64)
    if err != nil {
        return err
    }

    colorAttr.g, err = strconv.ParseFloat(line[1], 64)
    if err != nil {
        return err
    }

    colorAttr.b, err = strconv.ParseFloat(line[2], 64)
    if err != nil {
        return err
    }

    colorAttr.a, err = strconv.ParseFloat(line[3], 64)
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

    rgba := colorEdit.color.GetRGBA()
    colorAttr.r = rgba.GetRed()
    colorAttr.g = rgba.GetGreen()
    colorAttr.b = rgba.GetBlue()
    colorAttr.a = rgba.GetAlpha()

    return nil
}

type ColorEditor struct {
    box      *gtk.Box
    color    *gtk.ColorButton
}

func NewColorEditor(name string) *ColorEditor {
    var err error
    colorEdit := &ColorEditor{}

    colorEdit.box, err = gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 0)
    if err != nil {
        log.Print(err)
        return nil 
    }

    colorEdit.box.SetVisible(true)
    label, err := gtk.LabelNew(name)
    if err != nil { 
        log.Print(err)
        return nil 
    }

    label.SetVisible(true)
    label.SetWidthChars(12)
    colorEdit.box.PackStart(label, false, false, padding)

    colorEdit.color, err = gtk.ColorButtonNew()
    if err != nil {
        log.Print(err)
        return nil
    }

    colorEdit.color.SetVisible(true)
    colorEdit.box.PackStart(colorEdit.color, false, false, padding)

    return colorEdit
}

func (colorEdit *ColorEditor) Update(attr Attribute) error {
    colorAttr, ok := attr.(*ColorAttribute)
    if !ok {
        return fmt.Errorf("ColorEditor.Update requires ColorAttribute") 
    }

    rgb := gdk.NewRGBA(
        colorAttr.r, 
        colorAttr.g, 
        colorAttr.b, 
        colorAttr.a,
    )

    colorEdit.color.SetRGBA(rgb)

    return nil
}

func (colorEdit *ColorEditor) Box() *gtk.Box {
    return colorEdit.box
}
