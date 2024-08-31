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
	Red   float64
	Green float64
	Blue  float64
	Alpha float64
}

func (colorAttr *ColorAttribute) ToString() string {
	return fmt.Sprintf("%f %f %f %f", colorAttr.Red, colorAttr.Blue, colorAttr.Green, colorAttr.Alpha)
}

func (colorAttr *ColorAttribute) FromString(s string) (err error) {
	nums := strings.Split(s, " ")
	if len(nums) != 4 {
		return fmt.Errorf("Incorrect color format: %s", s)
	}

	colorAttr.Red, err = strconv.ParseFloat(nums[0], 64)
	if err != nil {
		return
	}

	colorAttr.Blue, err = strconv.ParseFloat(nums[1], 64)
	if err != nil {
		return
	}

	colorAttr.Green, err = strconv.ParseFloat(nums[2], 64)
	if err != nil {
		return
	}

	colorAttr.Alpha, err = strconv.ParseFloat(nums[3], 64)
	return
}

func (colorAttr *ColorAttribute) UpdateAttribute(colorEdit *ColorEditor) (err error) {
	rgba := colorEdit.color.GetRGBA()
	colorAttr.Red = rgba.GetRed()
	colorAttr.Green = rgba.GetGreen()
	colorAttr.Blue = rgba.GetBlue()
	colorAttr.Alpha = colorEdit.opacity.GetValue()

	return
}

type ColorEditor struct {
	Box     *gtk.Box
	color   *gtk.ColorButton
	opacity *gtk.Scale
	Name    string
}

func NewColorEditor(name string) (colorEdit *ColorEditor, err error) {
	colorEdit = &ColorEditor{Name: name}

	colorEdit.Box, err = gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 0)
	if err != nil {
		return
	}

	colorEdit.Box.SetVisible(true)
	label1, err := gtk.LabelNew(name)
	if err != nil {
		return
	}

	label1.SetVisible(true)
	label1.SetWidthChars(12)
	colorEdit.Box.PackStart(label1, false, false, padding)

	colorEdit.color, err = gtk.ColorButtonNew()
	if err != nil {
		return
	}

	colorEdit.color.SetVisible(true)
	colorEdit.Box.PackStart(colorEdit.color, false, false, padding)

	label2, err := gtk.LabelNew("Opacity")
	if err != nil {
		return
	}

	label2.SetVisible(true)
	label2.SetWidthChars(12)
	colorEdit.Box.PackStart(label2, false, false, padding)

	colorEdit.opacity, err = gtk.ScaleNewWithRange(gtk.ORIENTATION_HORIZONTAL, 0, 1, 0.01)
	if err != nil {
		return
	}

	colorEdit.opacity.SetVisible(true)
	colorEdit.Box.PackStart(colorEdit.opacity, true, true, padding)

	return
}

func (colorEdit *ColorEditor) UpdateEditor(colorAttr *ColorAttribute) error {
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
