package geometry

import (
	"chroma-viz/library/attribute"
	"strings"
)

type Rectangle struct {
	Geometry

	Width    attribute.IntAttribute
	Height   attribute.IntAttribute
	Rounding attribute.IntAttribute
	Color    attribute.ColorAttribute
}

func (r *Rectangle) UpdateGeometry(rEdit *RectangleEditor) (err error) {
	err = r.Geometry.UpdateGeometry(&rEdit.GeometryEditor)
	if err != nil {
		return
	}

	r.Width.UpdateAttribute(&rEdit.Width)
	if err != nil {
		return
	}

	r.Height.UpdateAttribute(&rEdit.Height)
	if err != nil {
		return
	}

	r.Rounding.UpdateAttribute(&rEdit.Rounding)
	if err != nil {
		return
	}

	r.Color.UpdateAttribute(&rEdit.Color)
	if err != nil {
		return
	}

	return
}

func (r *Rectangle) EncodeEngine(b strings.Builder) {

}

type RectangleEditor struct {
	GeometryEditor

	Width    attribute.IntEditor
	Height   attribute.IntEditor
	Rounding attribute.IntEditor
	Color    attribute.ColorAttribute
}

func NewRectangleEditor() *RectangleEditor {
	return nil
}

func (rEdit *RectangleEditor) UpdateEditor(r *Rectangle) (err error) {
	err = rEdit.GeometryEditor.UpdateEditor(&r.Geometry)
	if err != nil {
		return
	}

	err = rEdit.Width.UpdateEditor(&r.Width)
	if err != nil {
		return
	}

	err = rEdit.Height.UpdateEditor(&r.Height)
	if err != nil {
		return
	}

	err = rEdit.Rounding.UpdateEditor(&r.Rounding)
	if err != nil {
		return
	}

	err = rEdit.Color.UpdateEditor(&r.Color)
	return

}
