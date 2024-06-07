package attribute

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"
)

type ColorAttribute struct {
	Name  string
	Type  int
	Red   float64
	Green float64
	Blue  float64
	Alpha float64
}

func NewColorAttribute(name string) *ColorAttribute {
	colorAttr := &ColorAttribute{
		Name: name, Type: COLOR,
		Red: 1.0, Green: 1.0, Blue: 1.0, Alpha: 1.0,
	}
	return colorAttr
}

func (colorAttr *ColorAttribute) String() string {
	return fmt.Sprintf("%s=%f %f %f %f#", colorAttr.Name,
		colorAttr.Red, colorAttr.Green, colorAttr.Blue, colorAttr.Alpha)
}

func (colorAttr *ColorAttribute) Encode() string {
	return fmt.Sprintf("%f %f %f %f", colorAttr.Red, colorAttr.Green, colorAttr.Blue, colorAttr.Alpha)
}

func (colorAttr *ColorAttribute) Decode(value string) (err error) {
	if value == "" {
		return
	}

	s := strings.Split(value, " ")
	if len(s) < 4 {
		err = fmt.Errorf("Error decoding color attr (%s)", value)
		return
	}

	colorAttr.Red, err = strconv.ParseFloat(s[0], 64)
	if err != nil {
		return
	}

	colorAttr.Green, err = strconv.ParseFloat(s[1], 64)
	if err != nil {
		return
	}

	colorAttr.Blue, err = strconv.ParseFloat(s[2], 64)
	if err != nil {
		return
	}

	colorAttr.Alpha, err = strconv.ParseFloat(s[3], 64)
	if err != nil {
		return
	}

	return
}

func (colorAttr *ColorAttribute) Copy(attr Attribute) (err error) {
	colorAttrCopy, ok := attr.(*ColorAttribute)
	if !ok {
		err = fmt.Errorf("Attribute not ColorAttribute")
		return
	}

	colorAttr.Red = colorAttrCopy.Red
	colorAttr.Green = colorAttrCopy.Green
	colorAttr.Blue = colorAttrCopy.Blue
	colorAttr.Alpha = colorAttrCopy.Alpha
	return
}

func (colorAttr *ColorAttribute) Update(edit Editor) (err error) {
	colorEdit, ok := edit.(*ColorEditor)
	if !ok {
		err = fmt.Errorf("ColorAttribute.Update requires ColorEditor")
		return
	}

	rgba := colorEdit.color.GetRGBA()
	colorAttr.Red = rgba.GetRed()
	colorAttr.Green = rgba.GetGreen()
	colorAttr.Blue = rgba.GetBlue()
	colorAttr.Alpha = colorEdit.opacity.GetValue()

	return
}

type ColorEditor struct {
	box     *gtk.Box
	color   *gtk.ColorButton
	opacity *gtk.Scale
	name    string
}

func NewColorEditor(name string) (colorEdit *ColorEditor, err error) {
	colorEdit = &ColorEditor{name: name}

	colorEdit.box, err = gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 0)
	if err != nil {
		return
	}

	colorEdit.box.SetVisible(true)
	label1, err := gtk.LabelNew(name)
	if err != nil {
		return
	}

	label1.SetVisible(true)
	label1.SetWidthChars(12)
	colorEdit.box.PackStart(label1, false, false, padding)

	colorEdit.color, err = gtk.ColorButtonNew()
	if err != nil {
		return
	}

	colorEdit.color.SetVisible(true)
	colorEdit.box.PackStart(colorEdit.color, false, false, padding)

	label2, err := gtk.LabelNew("Opacity")
	if err != nil {
		return
	}

	label2.SetVisible(true)
	label2.SetWidthChars(12)
	colorEdit.box.PackStart(label2, false, false, padding)

	colorEdit.opacity, err = gtk.ScaleNewWithRange(gtk.ORIENTATION_HORIZONTAL, 0, 1, 0.01)
	if err != nil {
		return
	}

	colorEdit.opacity.SetVisible(true)
	colorEdit.box.PackStart(colorEdit.opacity, true, true, padding)

	return
}

func (colorEdit *ColorEditor) Name() string {
	return colorEdit.name
}

func (colorEdit *ColorEditor) Update(attr Attribute) error {
	colorAttr, ok := attr.(*ColorAttribute)
	if !ok {
		return fmt.Errorf("ColorEditor.Update requires ColorAttribute")
	}

	rgb := gdk.NewRGBA(
		colorAttr.Red,
		colorAttr.Green,
		colorAttr.Blue,
		colorAttr.Alpha,
	)

	colorEdit.color.SetRGBA(rgb)
	colorEdit.opacity.SetValue(colorAttr.Alpha)

	return nil
}

func (colorEdit *ColorEditor) Box() *gtk.Box {
	return colorEdit.box
}

func (colorEdit *ColorEditor) Expand() bool {
	return false
}
